package main

import (
	"context"
	"example/Orcha/api"
	"example/Orcha/logger"
	"example/Orcha/proxmoxClient"
	"example/Orcha/sshClient"
	"example/Orcha/util"
	"log/slog"
)

func main() {
	// Defines default logger with slog
	slog.Info("[Starting Orcha] - Setting up default logger...")
	err := logger.LoadLogger()
	if err != nil {
		slog.Error("[Starting Orcha] - Could not set a default logger: " + err.Error())
		return
	}

	// Loads config
	slog.Info("[Starting Orcha] - Loading configuration file...")
	config, err := util.LoadConfig(".")
	if err != nil {
		slog.Error("[Starting Orcha] - Could not load from config: " + err.Error())
		return
	}

	// Creates context
	slog.Info("[Starting Orcha] - Creating default context...")
	ctx := context.Background()

	// Proxmox Client
	slog.Info("[Starting Orcha] - Establishing Proxmox connection...")
	pmClient, err := proxmoxClient.NewClient(config, ctx)
	if err != nil {
		slog.Error("[Starting Orcha] - Could not establish Proxmox connection: " + err.Error())
		return
	}

	// SSH Client
	slog.Info("[Starting Orcha] - Establishing SSH connection...")
	sshClient, err := sshClient.NewSshClient(config)
	if err != nil {
		slog.Error("[Starting Orcha] - Could not establish ssh connection: " + err.Error())
		return
	}

	// Close SSH Connection
	defer sshClient.CloseConnection()

	// API Server
	slog.Info("[Starting Orcha] - Creating API server...")
	server, err := api.NewServer(config, ctx, pmClient, sshClient)
	if err != nil {
		slog.Error("[Starting Orcha] - Could not create server: " + err.Error())
		return
	}

	// API server start
	slog.Info("[Starting Orcha] - Starting API server...")
	err = server.Start(":8090")
	if err != nil {
		slog.Error("[Starting Orcha] - Could not start server: " + err.Error())
		return
	}
}
