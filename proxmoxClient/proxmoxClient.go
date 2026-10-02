package proxmoxClient

import (
	"context"
	"fmt"
	"time"

	"github.com/BrunoBiz/Orcha/util"

	"github.com/luthermonson/go-proxmox"
)

type ProxmoxClient struct {
	Client *proxmox.Client
	Node   *proxmox.Node
	config util.Config
	ctx    context.Context
}

func NewClient(config util.Config, ctx context.Context) (*ProxmoxClient, error) {
	// Proxmox API Client connection
	client := proxmox.NewClient(config.PVEUrl,
		proxmox.WithAPIToken(config.PVEUserRealm+"!"+config.PVETokenID, config.PVEToken),
		proxmox.WithTimeout(30*time.Second), // http.DefaultClient has no timeout
	)

	// Checks if the client is answering
	_, err := client.Version(ctx)
	if err != nil || client == nil {
		return nil, err
	}

	// Main node - This API will load only one node from the proxmox server
	node, err := client.Node(ctx, config.PVENodeName)
	if err != nil {
		return nil, err
	}

	proxmoxClient := &ProxmoxClient{
		Client: client,
		Node:   node,
		config: config,
		ctx:    ctx,
	}

	return proxmoxClient, nil
}

// Testing -- ignore
func (pmClient *ProxmoxClient) PingCtn() {
	ctnList, _ := pmClient.Node.Containers(pmClient.ctx)
	var totalMem uint64

	for i := 0; i < len(ctnList); i++ {
		ctn := ctnList[i]

		fmt.Println(ctn.Name)
		fmt.Println(ctn.MaxMem / 1048576) // Memória máxima configurada para o container atual - Convertido para MB (de bytes)

		if ctn.Status == "running" { // Caso esteja rodando o container - soma o valor total de memória
			totalMem += ctn.MaxMem / 1048576
		}

		err := pmClient.Client.Get(pmClient.ctx, fmt.Sprintf("/nodes/%s/lxc/%d/config", ctn.Node, ctn.VMID), &ctn.ContainerConfig)
		if err == nil {
			fmt.Println(ctn.ContainerConfig.Hostname)
			fmt.Println(ctn.ContainerConfig.Nets)
			fmt.Println(ctn.ContainerConfig.OnBoot)
		}
		fmt.Println("--------------------------")
	}

	fmt.Printf("\nTotal Memory Used: %dMb | %dGb", totalMem, totalMem/1024)
}
