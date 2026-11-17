package logger

import (
	"log/slog"
	"testing"
)

func TestLogLevel(t *testing.T) {
	for _, tc := range []struct {
		value string
		want  slog.Level
	}{{"", slog.LevelInfo}, {"INFO", slog.LevelInfo}, {"debug", slog.LevelDebug}, {"warn", slog.LevelWarn}, {"error", slog.LevelError}} {
		got, err := logLevel(tc.value)
		if err != nil || got != tc.want {
			t.Fatalf("logLevel(%q) = %v, %v", tc.value, got, err)
		}
	}
	if _, err := logLevel("verbose"); err == nil {
		t.Fatal("invalid log level should fail")
	}
}
