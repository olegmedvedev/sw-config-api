package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/ogen-go/ogen/ogenerrors"
)

// CustomErrorHandler provides better error messages for API consumers
func CustomErrorHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
	// Set content type
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Handle different types of errors using errors.As for wrapped errors
	var decodeParamsErr *ogenerrors.DecodeParamsError
	if errors.As(err, &decodeParamsErr) {
		// Handle parameter decoding errors
		handleDecodeParamsError(w, decodeParamsErr)
		return
	}

	var decodeParamErr *ogenerrors.DecodeParamError
	if errors.As(err, &decodeParamErr) {
		// Handle single parameter decoding errors
		handleDecodeParamError(w, decodeParamErr)
		return
	}

	// Handle other errors
	handleGenericError(w, err)
}

func handleDecodeParamsError(w http.ResponseWriter, err *ogenerrors.DecodeParamsError) {
	// Extract the underlying parameter error
	var decodeErr *ogenerrors.DecodeParamError
	if errors.As(err.Err, &decodeErr) {
		handleDecodeParamError(w, decodeErr)
		return
	}

	// Fallback for other decode errors
	w.WriteHeader(http.StatusBadRequest)
	response := map[string]interface{}{
		"error": map[string]interface{}{
			"code":    400,
			"message": "Invalid request parameters",
		},
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Log error but can't do much more since headers are already written
		fmt.Printf("Failed to encode error response: %v\n", err)
	}
}

func handleDecodeParamError(w http.ResponseWriter, err *ogenerrors.DecodeParamError) {
	w.WriteHeader(http.StatusBadRequest)

	// Create user-friendly error message
	message := fmt.Sprintf("Missing required parameter: %s", err.Name)

	response := map[string]interface{}{
		"error": map[string]interface{}{
			"code":    400,
			"message": message,
		},
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Log error but can't do much more since headers are already written
		fmt.Printf("Failed to encode error response: %v\n", err)
	}
}

func handleGenericError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)

	// Check if it's a not found error
	if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "Configuration not found") {
		w.WriteHeader(http.StatusNotFound)
		response := map[string]interface{}{
			"error": map[string]interface{}{
				"code":    404,
				"message": "Configuration not found",
			},
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			// Log error but can't do much more since headers are already written
			fmt.Printf("Failed to encode error response: %v\n", err)
		}
		return
	}

	// Default internal server error
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
