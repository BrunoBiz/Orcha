package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/BrunoBiz/Orcha/logger"
)

func (server *Server) postStartServer(c context.Context, input *ServerRequest) (*ServerResponse, error) {
	slog.Log(c, logger.LevelFile, "[postStartServer] - API CALL")

	var cntID uint64
	var serverResponse = NewErrorResponse("An error occurred while processing the request.", "INTERNAL_SERVER_ERROR")

	// Parameter sent via URL
	cntID = input.CntID
	slog.Log(c, logger.LevelFile, "[postStartServer] - cntID: "+strconv.FormatUint(cntID, 10))
	slog.Log(c, logger.LevelFile, "[postStartServer] - User: "+input.Body.User)

	// Prepares the command to start the server
	commandStart := fmt.Sprintf(`pct exec %d -- bash -c "su -s /bin/bash %s -c 'cd ~ && ./Narwhal start'"`,
		cntID,
		input.Body.User)
	slog.Log(c, logger.LevelFile, "[postStartServer] - commandDetails: "+commandStart)

	// Sends the command via SSH, returns the combined output - stdout + stderr
	optStartReturn, err := server.sshClient.NewSession(commandStart)

	if err != nil {
		slog.Error("[postStartServer] - SSH New Session: " + err.Error())
		serverResponse = NewErrorResponse("An error occurred while processing the request.", "INTERNAL_SERVER_ERROR")
		return &serverResponse, err
	}

	// Started
	if strings.Contains(optStartReturn, "[  OK  ] Starting") || strings.Contains(optStartReturn, "MESSAGE: Server started") {
		slog.Log(c, logger.LevelFile, "[postStartServer] - Server started successfully")
		serverResponse = NewSuccessfulResponse(true, "Server started successfully.")
		return &serverResponse, nil
	}

	// Already running
	if strings.Contains(optStartReturn, "is already running") {
		slog.Log(c, logger.LevelFile, "[postStartServer] - Server is already running")
		serverResponse = NewSuccessfulResponse(false, "Server is already running.")
		return &serverResponse, nil
	}

	// Error - If can't start server, returns an error
	return &serverResponse, errors.New("Unable to start server")
}
