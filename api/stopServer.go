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

func (server *Server) postStopServer(c context.Context, input *ServerRequest) (*ServerResponse, error) {
	slog.Log(c, logger.LevelFile, "[postStopServer] - API CALL")

	var cntID uint64
	var serverResponse = NewErrorResponse("An error occurred while processing the request.", "INTERNAL_SERVER_ERROR")

	// Parameter sent via URL
	cntID = input.CntID
	slog.Log(c, logger.LevelFile, "[postStartServer] - cntID: "+strconv.FormatUint(cntID, 10))

	// Prepares the command to start the server
	commandStop := fmt.Sprintf(`pct exec %d -- bash -c "su -s /bin/bash gameserver -c 'cd ~ && ./Narwhal stop'"`,
		cntID)

	// Sends the command via SSH, returns the combined output - stdout + stderr
	optStopReturn, err := server.sshClient.NewSession(commandStop)
	slog.Log(c, logger.LevelFile, "[postStopServer] - commandDetails: "+optStopReturn)

	if err != nil {
		slog.Error("[postStopServer] - SSH New Session: " + err.Error())
		serverResponse = NewErrorResponse("An error occurred while processing the request.", "INTERNAL_SERVER_ERROR")
		return &serverResponse, err
	}

	// Server stopped
	if strings.Contains(optStopReturn, "[  OK  ] Stopping") || strings.Contains(optStopReturn, "MESSAGE: Server stopped") {
		slog.Log(c, logger.LevelFile, "[postStopServer] - Server stopped successfully")
		serverResponse = NewSuccessfulResponse(true, "Server stopped successfully")
		return &serverResponse, nil
	}

	// Server is already stopped
	if strings.Contains(optStopReturn, "is already stopped") {
		slog.Log(c, logger.LevelFile, "[postStopServer] - Server is already stopped")
		serverResponse = NewSuccessfulResponse(false, "Server is already stopped")
		return &serverResponse, nil
	}

	// Error - If can't stop server, returns an error
	slog.Log(c, logger.LevelFile, "[postStopServer] - OK")
	return &serverResponse, errors.New("Unable to stop server")
}
