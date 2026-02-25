package apperrors

import "errors"

var (
	ErrServiceUnavailable = errors.New("external service is unavailable")
	ErrBadRequestPayload  = errors.New("failed to prepare request payload")
)
