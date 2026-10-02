package api

import (
	"context"

	"github.com/BrunoBiz/Orcha/proxmoxClient"
	"github.com/BrunoBiz/Orcha/sshClient"
	"github.com/BrunoBiz/Orcha/util"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gin-gonic/gin"
	"github.com/luthermonson/go-proxmox"
)

// API Request Body
type ServerRequest struct {
	CntID uint64 `path:"id" maxLength:"5" example:"101" doc:"Container VMID"`
	Body  struct {
		User string `json:"user" binding:"required" example:"gameserver" doc:"Linux user where the game is running from."`
	}
}

// API response body
type ServerResponse struct {
	Body struct {
		Successful bool   `json:"successful" binding:"required" doc:"If the request was successful or not - False when an error occurs"`                                                                                                        // If the request was successful or not - False when an error occurs
		Status     bool   `json:"status" binding:"required" doc:"If the request could achieve what was intended / If the request is based on a boolean response E.g. If the server is running - True / If the container was not found - False"` // If the request could achieve what was intended / If the request is based on a boolean response E.g. If the server is running - True / If the container was not found - False
		Message    string `json:"message" binding:"required"`
		ErrorCode  string `json:"errorcode" binding:"required" doc:"Unused in a successful response"`
	}
}

// Container info from the proxmox package, formatted to be used with HUMA
type ContainerInfo struct {
	CPUs    int                    `json:"cpus"`
	MaxDisk uint64                 `json:"maxdisk"`
	MaxMem  uint64                 `json:"maxmem"`
	MaxSwap uint64                 `json:"maxswap"`
	Name    string                 `json:"name"`
	Node    string                 `json:"node"`
	Status  string                 `json:"status"`
	Tags    string                 `json:"tags"`
	Uptime  uint64                 `json:"uptime"`
	VMID    proxmox.StringOrUint64 `json:"vmid"`
}

// An array of containers
type ContainerOutput struct {
	Body struct {
		Containers []ContainerInfo
	}
}

type Server struct {
	config    util.Config
	ctx       context.Context
	router    *gin.Engine
	pmClient  *proxmoxClient.ProxmoxClient
	sshClient *sshClient.SshClient
	humaAPI   huma.API
}

type ConteinerInput struct {
	CntID uint64 `path:"id" maxLength:"5" example:"101" doc:"Container VMID"`
}
