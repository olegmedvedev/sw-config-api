package service

import (
	"sw-config-api/internal/errors"
)

// NotFoundError is an alias for errors.NotFoundError
type NotFoundError = errors.NotFoundError

// IsNotFoundError is an alias for errors.IsNotFoundError
func IsNotFoundError(err error) bool {
	return errors.IsNotFoundError(err)
}
