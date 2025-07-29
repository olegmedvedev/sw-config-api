package middleware

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"sw-config-api/internal/api"

	"github.com/ogen-go/ogen/middleware"
	"github.com/rs/xid"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const (
	requestIDKey contextKey = "request_id"
)

// LoggingMiddleware creates middleware for logging HTTP requests
func LoggingMiddleware(logger *slog.Logger) api.Middleware {
	return func(req middleware.Request, next middleware.Next) (middleware.Response, error) {
		start := time.Now()

		// Generate unique request ID
		requestID := generateRequestID()

		// Add request ID to context
		ctx := context.WithValue(req.Context, requestIDKey, requestID)
		req.Context = ctx

		// Create logger with request context
		requestLogger := logger.With(
			"request_id", requestID,
			"method", req.Raw.Method,
			"path", req.Raw.URL.Path,
			"query", req.Raw.URL.RawQuery,
			"user_agent", req.Raw.UserAgent(),
			"remote_addr", req.Raw.RemoteAddr,
			"operation", req.OperationName,
			"operation_summary", req.OperationSummary,
		)

		// Log request start
		requestLogger.Info("request started")

		// Log request parameters at DEBUG level
		if len(req.Params) > 0 {
			params := make(map[string]interface{})
			for key, value := range req.Params {
				params[key.Name] = value
			}
			requestLogger.Debug("request parameters", "params", params)
		}

		// Execute next middleware/handler
		response, err := next(req)

		// Calculate request duration
		duration := time.Since(start)

		// Log request completion
		if err != nil {
			requestLogger.Error("request failed",
				"duration_ms", duration.Milliseconds(),
				"error", err.Error(),
			)
		} else {
			requestLogger.Info("request completed",
				"duration_ms", duration.Milliseconds(),
			)
		}

		return response, err
	}
}

// generateRequestID generates a unique ID for the request using strings.Builder
func generateRequestID() string {
	var builder strings.Builder

	// Add timestamp prefix
	builder.WriteString(time.Now().Format("20060102150405"))
	builder.WriteString("-")

	// Add unique identifier using xid
	builder.WriteString(xid.New().String())

	return builder.String()
}
