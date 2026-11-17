package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakeCommander struct {
	output  string
	err     error
	command string
}

func (f *fakeCommander) NewSession(command string) (string, error) {
	f.command = command
	return f.output, f.err
}

type remoteExit struct{}

func (remoteExit) Error() string   { return "Process exited with status 1" }
func (remoteExit) ExitStatus() int { return 1 }

func operationRequest(t *testing.T, action, output string, err error) (*httptest.ResponseRecorder, map[string]any, string) {
	t.Helper()
	var logs bytes.Buffer
	previous := slog.Default()
	slog.SetDefault(slog.New(slog.NewJSONHandler(&logs, &slog.HandlerOptions{Level: slog.LevelInfo})))
	t.Cleanup(func() { slog.SetDefault(previous) })
	command := &fakeCommander{output: output, err: err}
	server := &Server{sshClient: command}
	server.setupRouter()
	request := httptest.NewRequest("POST", "/containers/server/111/"+action+"?token=hidden-query", strings.NewReader("hidden-body"))
	request.Header.Set("Authorization", "Bearer hidden-token")
	response := httptest.NewRecorder()
	server.router.ServeHTTP(response, request)
	lines := strings.Split(strings.TrimSpace(logs.String()), "\n")
	if len(lines) != 1 {
		t.Fatalf("expected one request record, got %d: %s", len(lines), logs.String())
	}
	var record map[string]any
	if err := json.Unmarshal([]byte(lines[0]), &record); err != nil {
		t.Fatal(err)
	}
	for _, value := range []string{"hidden-query", "hidden-body", "hidden-token"} {
		if strings.Contains(logs.String(), value) {
			t.Fatalf("request data leaked into log: %s", value)
		}
	}
	if record["request_id"] != response.Header().Get("X-Request-ID") || record["request_id"] == "" {
		t.Fatal("missing or mismatched request ID")
	}
	return response, record, command.command
}

func TestRemoteFailureCapturesOutputAndExitCode(t *testing.T) {
	response, record, command := operationRequest(t, "restart", "\x1b[31msu: user gameserver does not exist\x1b[0m\n", remoteExit{})
	if response.Code != 500 {
		t.Fatalf("HTTP %d: %s", response.Code, response.Body.String())
	}
	if record["output"] != "su: user gameserver does not exist" || record["exit_code"] != float64(1) || record["vmid"] != float64(111) || record["action"] != "restart" {
		t.Fatalf("missing diagnostics: %#v", record)
	}
	if record["command"] != command || !strings.Contains(command, "./Narwhal restart") {
		t.Fatal("wrong failed command")
	}
	if !strings.Contains(response.Body.String(), response.Header().Get("X-Request-ID")) {
		t.Fatal("error response missing correlation ID")
	}
	if strings.Contains(response.Body.String(), "user gameserver does not exist") {
		t.Fatal("raw manager output should stay in logs")
	}
}

func TestUnrecognizedOutputIsAlsoCaptured(t *testing.T) {
	response, record, _ := operationRequest(t, "restart", "manager: unknown startup failure", nil)
	if response.Code != 500 || record["error_code"] != "UNRECOGNIZED_MANAGER_OUTPUT" || record["output"] != "manager: unknown startup failure" {
		t.Fatalf("missing parse diagnostics: %#v", record)
	}
}

func TestDetailsReportsActualRunningState(t *testing.T) {
	for _, tc := range []struct {
		output string
		online bool
	}{{"Status: STARTED", true}, {"\x1b[31mStatus: STOPPED\x1b[0m", false}} {
		t.Run(fmt.Sprint(tc.online), func(t *testing.T) {
			response, record, _ := operationRequest(t, "details", tc.output, nil)
			var body map[string]any
			if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			if response.Code != 200 || body["successful"] != true || body["status"] != tc.online {
				t.Fatalf("wrong details result: %s", response.Body.String())
			}
			if record["level"] != "INFO" || record["output"] != nil {
				t.Fatalf("successful poll should be concise: %#v", record)
			}
		})
	}
}

func TestDiagnosticOutputIsBounded(t *testing.T) {
	response, record, _ := operationRequest(t, "start", "first error\n"+strings.Repeat("x", 20000)+"\nlast error", remoteExit{})
	output := record["output"].(string)
	if response.Code != 500 || record["output_truncated"] != true || len(output) > maxDiagnosticOutput+100 || !strings.HasPrefix(output, "first error") || !strings.HasSuffix(output, "last error") {
		t.Fatalf("unexpected truncation: %d bytes", len(output))
	}
}
