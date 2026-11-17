package api

import (
	"context"
	"crypto/rand"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

type requestLogKey struct{}

type requestLog struct {
	id     string
	fields []any
}

// One completion record per request; handlers attach failure diagnostics here.
func requestLogging() gin.HandlerFunc {
	return func(c *gin.Context) {
		started := time.Now()
		record := &requestLog{id: rand.Text()}
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), requestLogKey{}, record))
		c.Header("X-Request-ID", record.id)
		c.Next()

		level := slog.LevelInfo
		if c.Writer.Status() >= 500 {
			level = slog.LevelError
		} else if c.Writer.Status() >= 400 {
			level = slog.LevelWarn
		}
		fields := []any{"request_id", record.id, "method", c.Request.Method, "path", c.Request.URL.Path,
			"status", c.Writer.Status(), "duration_ms", time.Since(started).Milliseconds()}
		fields = append(fields, record.fields...)
		slog.Log(c.Request.Context(), level, "request completed", fields...)
	}
}

func requestID(ctx context.Context) string {
	if record, ok := ctx.Value(requestLogKey{}).(*requestLog); ok {
		return record.id
	}
	return "unavailable"
}

func addRequestFields(ctx context.Context, fields ...any) {
	if record, ok := ctx.Value(requestLogKey{}).(*requestLog); ok {
		record.fields = append(record.fields, fields...)
	} else {
		slog.ErrorContext(ctx, "operation failed", fields...)
	}
}
