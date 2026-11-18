package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"testing"
)

func TestOperationFailureLogsRemoteOutput(t *testing.T) {
	var logs bytes.Buffer
	previous := slog.Default()
	slog.SetDefault(slog.New(slog.NewJSONHandler(&logs, nil)))
	t.Cleanup(func() { slog.SetDefault(previous) })

	failure := errors.New("Process exited with status 1")
	got := operationFailure(context.Background(), 101, "details",
		"\x1b[31msu: user gameserver does not exist\x1b[0m\n", failure)
	if got != failure {
		t.Fatal("the original error must be returned unchanged")
	}

	var entry map[string]any
	if err := json.Unmarshal(logs.Bytes(), &entry); err != nil {
		t.Fatal(err)
	}
	if entry["level"] != "ERROR" || entry["vmid"] != float64(101) ||
		entry["action"] != "details" || entry["error"] != failure.Error() ||
		entry["output"] != "su: user gameserver does not exist" {
		t.Fatalf("missing failure diagnostics: %#v", entry)
	}
}
