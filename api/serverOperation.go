package api

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/acarl005/stripansi"
	"github.com/danielgtaylor/huma/v2"
)

var startedStatus = regexp.MustCompile(`(?mi)(status:)\s+(started)\b`)
var stoppedStatus = regexp.MustCompile(`(?mi)(status:)\s+(stopped)\b`)

const maxDiagnosticOutput = 8192

func diagnosticOutput(output string) string {
	output = strings.TrimSpace(stripansi.Strip(output))
	if len(output) > maxDiagnosticOutput {
		return strings.ToValidUTF8(output[:maxDiagnosticOutput/2]+"\n[... output truncated ...]\n"+output[len(output)-maxDiagnosticOutput/2:], "�")
	}
	return strings.ToValidUTF8(output, "�")
}

// Capture both non-zero remote exits and successful commands whose output is
// not recognized. Raw command output stays in server logs, not HTTP responses.
func operationFailure(ctx context.Context, vmid uint64, action, command, output string, err error) error {
	code := "UNRECOGNIZED_MANAGER_OUTPUT"
	fields := []any{"vmid", vmid, "action", action, "command", command, "output", diagnosticOutput(output),
		"output_truncated", len(strings.TrimSpace(stripansi.Strip(output))) > maxDiagnosticOutput}
	if err != nil {
		code = "SSH_COMMAND_FAILED"
		fields = append(fields, "error", err.Error())
		var exitErr interface{ ExitStatus() int }
		if errors.As(err, &exitErr) {
			fields = append(fields, "exit_code", exitErr.ExitStatus())
		}
	}
	fields = append(fields, "error_code", code)
	addRequestFields(ctx, fields...)
	return huma.Error500InternalServerError(fmt.Sprintf("Server %s failed for VMID %d (%s). Request ID: %s", action, vmid, code, requestID(ctx)))
}

func (server *Server) runServerOperation(ctx context.Context, input *ServerRequest, action string) (*ServerResponse, error) {
	command := fmt.Sprintf(`pct exec %d -- bash -c "su -s /bin/bash gameserver -c 'cd ~ && ./Narwhal %s'"`, input.CntID, action)
	output, err := server.sshClient.NewSession(command)
	if err != nil {
		return nil, operationFailure(ctx, input.CntID, action, command, output, err)
	}

	clean := stripansi.Strip(output)
	var response ServerResponse
	switch action {
	case "details":
		if startedStatus.MatchString(clean) {
			response = NewSuccessfulResponse(true, "Server running")
		} else if stoppedStatus.MatchString(clean) {
			response = NewSuccessfulResponse(false, "Server stopped")
		} else {
			return nil, operationFailure(ctx, input.CntID, action, command, output, nil)
		}
	case "start", "restart":
		if strings.Contains(clean, "[  OK  ] Starting") || strings.Contains(clean, "MESSAGE: Server started") {
			response = NewSuccessfulResponse(true, "Server started successfully.")
		} else if strings.Contains(clean, "is already running") {
			response = NewSuccessfulResponse(false, "Server is already running.")
		} else {
			return nil, operationFailure(ctx, input.CntID, action, command, output, nil)
		}
	case "stop":
		if strings.Contains(clean, "[  OK  ] Stopping") || strings.Contains(clean, "MESSAGE: Server stopped") {
			response = NewSuccessfulResponse(true, "Server stopped successfully")
		} else if strings.Contains(clean, "is already stopped") {
			response = NewSuccessfulResponse(false, "Server is already stopped")
		} else {
			return nil, operationFailure(ctx, input.CntID, action, command, output, nil)
		}
	}
	return &response, nil
}
