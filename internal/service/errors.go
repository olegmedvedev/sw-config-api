package service

import (
	"errors"
	"fmt"
)

// NotFoundError represents a not found error with details
type NotFoundError struct {
	Platform   string
	AppVersion string
}

func (e *NotFoundError) Error() string {
	if e.AppVersion != "" {
		return fmt.Sprintf("configuration not found for appVersion %s (%s)", e.AppVersion, e.Platform)
	}
	return fmt.Sprintf("configuration not found for %s", e.Platform)
}

func IsNotFoundError(err error) bool {
	return errors.Is(err, &NotFoundError{})
}
