package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sw-config-api/internal/api"
	apperr "sw-config-api/internal/errors"

	"github.com/ogen-go/ogen/ogenerrors"
)

// CustomErrorHandler provides better error messages for API consumers
func CustomErrorHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, err error, logger *slog.Logger) {
	// Set content type
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Handle different types of errors using errors.As for wrapped errors
	var decodeParamsErr *ogenerrors.DecodeParamsError
	if errors.As(err, &decodeParamsErr) {
		// Handle parameter decoding errors
		handleDecodeParamsError(logger, w, decodeParamsErr)
		return
	}

	var decodeParamErr *ogenerrors.DecodeParamError
	if errors.As(err, &decodeParamErr) {
		// Handle single parameter decoding errors
		handleDecodeParamError(logger, w, decodeParamErr)
		return
	}

	// Handle other errors
	handleGenericError(logger, w, err)
}

func handleDecodeParamsError(logger *slog.Logger, w http.ResponseWriter, err *ogenerrors.DecodeParamsError) {
	// Extract the underlying parameter error
	var decodeErr *ogenerrors.DecodeParamError
	if errors.As(err.Err, &decodeErr) {
		handleDecodeParamError(logger, w, decodeErr)
		return
	}

	// Fallback for other decode errors
	w.WriteHeader(http.StatusBadRequest)

	// Log the validation error
	logger.Warn("request validation failed",
		"error", "Invalid request parameters",
	)

	// Use generated error types
	errorResponse := &api.ConfigGetBadRequest{
		Error: api.NewOptConfigGetBadRequestError(api.ConfigGetBadRequestError{
			Code:    api.NewOptInt(400),
			Message: api.NewOptString("Invalid request parameters"),
		}),
	}

	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		// Log error but can't do much more since headers are already written
		fmt.Printf("Failed to encode error response: %v\n", err)
	}
}

func handleDecodeParamError(logger *slog.Logger, w http.ResponseWriter, err *ogenerrors.DecodeParamError) {
	w.WriteHeader(http.StatusBadRequest)

	// Create user-friendly error message
	message := fmt.Sprintf("Missing required parameter: %s", err.Name)

	// Log the validation error
	logger.Warn("request validation failed",
		"error", message,
		"parameter", err.Name,
	)

	// Use generated error types
	errorResponse := &api.ConfigGetBadRequest{
		Error: api.NewOptConfigGetBadRequestError(api.ConfigGetBadRequestError{
			Code:    api.NewOptInt(400),
			Message: api.NewOptString(message),
		}),
	}

	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		// Log error but can't do much more since headers are already written
		fmt.Printf("Failed to encode error response: %v\n", err)
	}
}

func handleGenericError(logger *slog.Logger, w http.ResponseWriter, err error) {
	// Check if it's a not found error
	if apperr.IsNotFoundError(err) {
		w.WriteHeader(http.StatusNotFound)

		// Use generated error types
		errorResponse := &api.ConfigGetNotFound{
			Error: api.NewOptConfigGetNotFoundError(api.ConfigGetNotFoundError{
				Code:    api.NewOptInt(404),
				Message: api.NewOptString("Configuration not found"),
			}),
		}

		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			// Log error but can't do much more since headers are already written
			fmt.Printf("Failed to encode error response: %v\n", err)
		}
		return
	}

	// Default internal server error
	w.WriteHeader(http.StatusInternalServerError)

	// For 500 errors, we still need to use a generic response since there's no generated type
	response := map[string]interface{}{
		"error": map[string]interface{}{
			"code":    500,
			"message": "Internal server error",
		},
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Log error but can't do much more since headers are already written
		fmt.Printf("Failed to encode error response: %v\n", err)
	}
}
