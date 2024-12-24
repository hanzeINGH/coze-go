package coze_error

import (
	"errors"
	"fmt"
)

type CozeError struct {
	ErrorCode    int
	ErrorMessage string
	LogID        string
}

func NewCozeError(code int, msg, logID string) *CozeError {
	return &CozeError{
		ErrorCode:    code,
		ErrorMessage: msg,
		LogID:        logID,
	}
}

// Error implements the error interface
func (e *CozeError) Error() string {
	return fmt.Sprintf("Code: %d, ErrorMessage: %s, LogID: %s",
		e.ErrorCode,
		e.ErrorMessage,
		e.LogID)
}

// AsCozeError checks if the error is of type CozeError
func AsCozeError(err error) (*CozeError, bool) {
	var cozeErr *CozeError
	if errors.As(err, &cozeErr) {
		return cozeErr, true
	}
	return nil, false
}
