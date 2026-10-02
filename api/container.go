package api

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/BrunoBiz/Orcha/logger"
)

func (server *Server) getContainers(c context.Context, input *struct{}) (*ContainerOutput, error) {
	slog.Log(c, logger.LevelFile, "[getContainers] - API CALL")

	var containerInfo ContainerInfo
	var containerOutput ContainerOutput

	// Container list in main node
	ctnList, err := server.pmClient.Node.Containers(server.ctx)
	slog.Log(c, logger.LevelFile, "[getContainers] - ctnList: "+strconv.Itoa(len(ctnList)))

	if err != nil {
		slog.Log(c, logger.LevelFile, "[getContainers] - ERROR: "+err.Error())
		return nil, err
	}

	for _, ctnRange := range ctnList {
		containerInfo.CPUs = ctnRange.CPUs
		containerInfo.MaxDisk = ctnRange.MaxDisk
		containerInfo.MaxMem = ctnRange.MaxMem
		containerInfo.MaxSwap = ctnRange.MaxSwap
		containerInfo.Name = ctnRange.Name
		containerInfo.Node = ctnRange.Node
		containerInfo.Status = ctnRange.Status
		containerInfo.Tags = ctnRange.Tags
		containerInfo.Uptime = ctnRange.Uptime
		containerInfo.VMID = ctnRange.VMID

		containerOutput.Body.Containers = append(containerOutput.Body.Containers, containerInfo)
	}

	slog.Log(c, logger.LevelFile, "[getContainers] - OK")
	return &containerOutput, err
}

func (server *Server) getContainerById(c context.Context, input *ContainerInput) (*ContainerOutput, error) {
	slog.Log(c, logger.LevelFile, "[getContainerById] - API CALL")

	var cntID uint64
	var containerInfo ContainerInfo
	var containerOutput ContainerOutput

	cntID = input.CntID
	slog.Log(c, logger.LevelFile, "[getContainerById] - cntID:"+strconv.FormatUint(cntID, 10))
	slog.Log(c, logger.LevelFile, "[getContainerById] - input.cntID:"+strconv.FormatUint(input.CntID, 10))

	// Container list in main node
	ctnList, err := server.pmClient.Node.Containers(server.ctx)
	slog.Log(c, logger.LevelFile, "[getContainerById] - ctnList: "+strconv.Itoa(len(ctnList)))

	if err != nil {
		slog.Log(c, logger.LevelFile, "[getContainerById] - ERROR: "+err.Error())
		return nil, err
	}

	// Looks for the container
	for i := 0; i < len(ctnList); i++ {
		if uint64(ctnList[i].VMID) == cntID {
			slog.Log(c, logger.LevelFile, "[getContainerById] - Container found")
			containerInfo.CPUs = ctnList[i].CPUs
			containerInfo.MaxDisk = ctnList[i].MaxDisk
			containerInfo.MaxMem = ctnList[i].MaxMem
			containerInfo.MaxSwap = ctnList[i].MaxSwap
			containerInfo.Name = ctnList[i].Name
			containerInfo.Node = ctnList[i].Node
			containerInfo.Status = ctnList[i].Status
			containerInfo.Tags = ctnList[i].Tags
			containerInfo.Uptime = ctnList[i].Uptime
			containerInfo.VMID = ctnList[i].VMID

			containerOutput.Body.Containers = append(containerOutput.Body.Containers, containerInfo)
		}
	}

	slog.Log(c, logger.LevelFile, "[getContainerById] - OK")
	return &containerOutput, nil
}

func (server *Server) getContainerStatusById(c context.Context, input *ContainerInput) (*ServerResponse, error) {
	slog.Log(c, logger.LevelFile, "[getContainerStatusById] - API CALL")

	var cntID uint64
	var serverResponse = NewSuccessfulResponse(false, "No container found") // Default response - no container found

	cntID = input.CntID
	slog.Log(c, logger.LevelFile, "[getContainerById] - cntID:"+strconv.FormatUint(cntID, 10))
	slog.Log(c, logger.LevelFile, "[getContainerById] - input.cntID:"+strconv.FormatUint(input.CntID, 10))

	// Container list in main node
	ctnList, err := server.pmClient.Node.Containers(server.ctx)
	slog.Log(c, logger.LevelFile, "[getContainerStatusById] - ctnList: "+strconv.Itoa(len(ctnList)))

	if err != nil {
		slog.Log(c, logger.LevelFile, "[getContainerStatusById] - ERROR: "+err.Error())
		serverResponse = NewErrorResponse("An error occurred while processing the request.", "INTERNAL_SERVER_ERROR")
		return &serverResponse, err
	}

	// Looks for the container
	for i := 0; i < len(ctnList); i++ {
		if uint64(ctnList[i].VMID) == cntID {
			slog.Log(c, logger.LevelFile, "[getContainerStatusById] - Container found")
			serverResponse = NewSuccessfulResponse(true, ctnList[i].Status)
		}
	}

	slog.Log(c, logger.LevelFile, "[getContainerStatusById] - OK")
	return &serverResponse, nil
}
