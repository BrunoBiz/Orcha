package api

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/BrunoBiz/Orcha/logger"
	"github.com/BrunoBiz/Orcha/proxmoxClient"
	"github.com/BrunoBiz/Orcha/sshClient"
	"github.com/BrunoBiz/Orcha/util"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
)

func NewServer(config util.Config, ctx context.Context, proxmoxClient *proxmoxClient.ProxmoxClient, sshClient *sshClient.SshClient) (*Server, error) {
	server := &Server{
		config:    config,
		ctx:       ctx,
		pmClient:  proxmoxClient,
		sshClient: sshClient,
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	defaultLogger := slog.Default()

	configSlogGin := sloggin.Config{
		DefaultLevel: logger.LevelGIN,
		//ClientErrorLevel: slog.LevelWarn,
		//ServerErrorLevel: slog.LevelError,

		WithUserAgent:      false,
		WithRequestID:      true,
		WithRequestBody:    true,
		WithRequestHeader:  true,
		WithResponseBody:   true,
		WithResponseHeader: true,
		WithSpanID:         false,
		WithTraceID:        false,
		WithClientIP:       false,

		WithCustomMessage: func(c *gin.Context) string {
			return "GIN-Log"
		},
	}

	// Sets up the slog middleware for GIN
	router.Use(sloggin.NewWithConfig(defaultLogger, configSlogGin))

	// Sets up HUMA
	humaApi := humagin.New(router, huma.DefaultConfig("Orcha - Proxmox Orchestration API", "1.0.0"))
	openApiSpecs(humaApi)

	var tagContainers = []string{"Containers"}
	var tagServices = []string{"Services"}

	// Sets up routing
	// Uses the Proxmox API
	huma.Register(humaApi, huma.Operation{ // Returns info about all containers
		OperationID: "get-containers",
		Method:      http.MethodGet,
		Path:        "/containers",
		Summary:     "Retrieve all containers",
		Description: "Lists all available containers and information about each of them, such as id, name, uptime, status, etc.",
		Tags:        tagContainers,
	}, server.getContainers)

	huma.Register(humaApi, huma.Operation{ // Returns info about a specific container
		OperationID: "get-containersById",
		Method:      http.MethodGet,
		Path:        "/containers/{id}",
		Summary:     "Retrieve container via id",
		Description: "Returns information about a specific container, if found, via it's vmid.",
		Tags:        tagContainers,
	}, server.getContainerById)

	huma.Register(humaApi, huma.Operation{ // Returns if a specific container is running
		OperationID: "get-containerStatusById",
		Method:      http.MethodGet,
		Path:        "/containers/{id}/status",
		Summary:     "Check if container is running",
		Description: "Retrieves the current running state of the specified Proxmox LXC container.",
		Tags:        tagContainers,
	}, server.getContainerStatusById)

	// Uses SSH to connect to a container and run the commands
	huma.Register(humaApi, huma.Operation{ // Start server
		OperationID: "post-startServer",
		Method:      http.MethodPost,
		Path:        "/containers/server/{id}/start",
		Summary:     "Start server",
		Description: "Starts the container's associated game server",
		Tags:        tagServices,
	}, server.postStartServer)

	huma.Register(humaApi, huma.Operation{ // Server status -> Online/Offline
		OperationID: "post-detailsServer",
		Method:      http.MethodPost,
		Path:        "/containers/server/{id}/details",
		Summary:     "Check server status",
		Description: "Check if the container's associated game server is running",
		Tags:        tagServices,
	}, server.postDetailsServer)

	huma.Register(humaApi, huma.Operation{ // Stop server
		OperationID: "post-stopServer",
		Method:      http.MethodPost,
		Path:        "/containers/server/{id}/stop",
		Summary:     "Stop server",
		Description: "Stops the container's associated game server",
		Tags:        tagServices,
	}, server.postStopServer)

	huma.Register(humaApi, huma.Operation{ // Restart server
		OperationID: "post-restartServer",
		Method:      http.MethodPost,
		Path:        "/containers/server/{id}/restart",
		Summary:     "Restart server",
		Description: "Restarts the container's associated game server",
		Tags:        tagServices,
	}, server.postRestartServer)

	server.router = router
	server.humaAPI = humaApi
}

func (server *Server) Start(address string) error {
	return http.ListenAndServe(address, server.router)
}

func NewSuccessfulResponse(status bool, message string) ServerResponse {
	response := ServerResponse{}

	response.Body.Successful = true
	response.Body.Status = status
	response.Body.Message = message
	response.Body.ErrorCode = "" // Not used in a SUCCESSFUL response

	return response
}

func NewErrorResponse(message string, errorCode string) ServerResponse { // Will only be considered as an error status codes in the 5xx range
	response := ServerResponse{}

	response.Body.Successful = false
	response.Body.Status = false // Not used in an UNSUCCESSFUL response
	response.Body.Message = message
	response.Body.ErrorCode = errorCode

	return response
}

// Sets up required info for OpenAPI documentation
func openApiSpecs(api huma.API) {
	var contact = huma.Contact{
		Name:  "Bruno Biz Dias de Castro",
		Email: "brunobizdc@gmail.com",
		URL:   "https://github.com/BrunoBiz",
	}

	var license = huma.License{
		Name:       "MIT License",
		Identifier: "MIT",
	}

	var info = huma.Info{
		Title:       "Orcha - Proxmox Orchestration API",
		Description: "A RESTful API that provides access to Proxmox environment information and enables administrators to manage the lifecycle of game servers running inside Proxmox LXC containers.",
		Contact:     &contact,
		Version:     "1.0.0",
		License:     &license,
	}
	api.OpenAPI().Info = &info

	var servers = huma.Server{
		URL:         "http://192.168.18.162:8090",
		Description: "Local Proxmox container IP",
	}
	api.OpenAPI().Servers = append(api.OpenAPI().Servers, &servers)

	var tagContainer = huma.Tag{
		Name:        "Containers",
		Description: "Operations for retrieving and managing Proxmox LXC containers.",
	}

	var tagServices = huma.Tag{
		Name:        "Services",
		Description: "Operations for managing game server services running within Proxmox LXC containers.",
	}

	api.OpenAPI().Tags = append(api.OpenAPI().Tags, &tagContainer)
	api.OpenAPI().Tags = append(api.OpenAPI().Tags, &tagServices)

}
