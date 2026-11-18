package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"regexp"

	"github.com/BrunoBiz/Orcha/logger"

	"github.com/acarl005/stripansi"
)

func (server *Server) postDetailsServer(c context.Context, input *ServerRequest) (*ServerResponse, error) {

	var serverResponse = NewErrorResponse("An error occurred while processing the request.", "INTERNAL_SERVER_ERROR")
	var cntID uint64

	// Parameter sent via URL
	cntID = input.CntID

	// Prepares the command to check server details
	commandDetails := fmt.Sprintf(`pct exec %d -- bash -c "su -s /bin/bash gameserver -c 'cd ~ && ./Narwhal details'"`,
		cntID)

	// Sends the command via SSH, returns the combined output - stdout + stderr
	optDetailsReturn, err := server.sshClient.NewSession(commandDetails)

	if err != nil {
		serverResponse = NewErrorResponse("An error occurred while processing the request.", "INTERNAL_SERVER_ERROR")
		return &serverResponse, operationFailure(c, cntID, "details", optDetailsReturn, err)
	}

	// Server ONLINE
	if regexp.MustCompile(`(?mi)(status:)\s+(started)`).MatchString(stripansi.Strip(optDetailsReturn)) {
		slog.Log(c, logger.LevelFile, "[postDetailsServer] - Server running")
		serverResponse = NewSuccessfulResponse(true, "Server running")
		return &serverResponse, nil
	}

	// Server OFFLINE
	if regexp.MustCompile(`(?mi)(status:)\s+(stopped)`).MatchString(stripansi.Strip(optDetailsReturn)) {
		slog.Log(c, logger.LevelFile, "[postDetailsServer] - Server stopped")
		serverResponse = NewSuccessfulResponse(false, "Server stopped")
		return &serverResponse, nil
	}

	// Error - If can't check server details, returns an error
	return &serverResponse, operationFailure(c, cntID, "details", optDetailsReturn, errors.New("Unable to check server details"))
}
