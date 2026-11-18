package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/BrunoBiz/Orcha/logger"
)

func (server *Server) postRestartServer(c context.Context, input *ServerRequest) (*ServerResponse, error) {

	var cntID uint64
	var serverResponse = NewErrorResponse("An error occurred while processing the request.", "INTERNAL_SERVER_ERROR")

	// Parameter sent via URL
	cntID = input.CntID

	// Prepares the command to restart the server
	commandRestart := fmt.Sprintf(`pct exec %d -- bash -c "su -s /bin/bash gameserver -c 'cd ~ && ./Narwhal restart'"`,
		cntID)

	// Sends the command via SSH, returns the combined output - stdout + stderr
	optRestartReturn, err := server.sshClient.NewSession(commandRestart)

	if err != nil {
		serverResponse = NewErrorResponse("An error occurred while processing the request.", "INTERNAL_SERVER_ERROR")
		return &serverResponse, operationFailure(c, cntID, "restart", optRestartReturn, err)
	}

	// Started
	if strings.Contains(optRestartReturn, "[  OK  ] Starting") || strings.Contains(optRestartReturn, "MESSAGE: Server started") {
		slog.Log(c, logger.LevelFile, "[postRestartServer] - Server started successfully")
		serverResponse = NewSuccessfulResponse(true, "Server started successfully.")
		return &serverResponse, nil
	}

	// Already running
	if strings.Contains(optRestartReturn, "is already running") {
		slog.Log(c, logger.LevelFile, "[postRestartServer] - Server is already running")
		serverResponse = NewSuccessfulResponse(false, "Server is already running.")
		return &serverResponse, nil
	}

	// Error - If can't restart server, returns an error
	return &serverResponse, operationFailure(c, cntID, "restart", optRestartReturn, errors.New("Unable to restart server"))
}
