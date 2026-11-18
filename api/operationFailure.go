package api

import (
	"context"
	"log/slog"
	"strings"

	"github.com/acarl005/stripansi"
)

// Log the output already captured over SSH before returning the original error.
func operationFailure(ctx context.Context, vmid uint64, action, output string, err error) error {
	slog.ErrorContext(ctx, "Server operation failed",
		"vmid", vmid,
		"action", action,
		"error", err,
		"output", strings.TrimSpace(stripansi.Strip(output)),
	)
	return err
}
