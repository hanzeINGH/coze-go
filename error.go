package coze

import (
	"errors"
	"fmt"
)

type CozeError struct {
	Code    int
	Message string
	LogID   string
}

func NewCozeError(code int, msg, logID string) *CozeError {
	return &CozeError{
		Code:    code,
		Message: msg,
		LogID:   logID,
	}
}

// Error implements the error interface
func (e *CozeError) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s, LogID: %s",
		e.Code,
		e.Message,
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
