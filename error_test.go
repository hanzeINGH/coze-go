package coze

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCozeError(t *testing.T) {
	// 测试创建新的 CozeError
	err := NewCozeError(1001, "test error", "test-log-id")
	assert.NotNil(t, err)
	assert.Equal(t, 1001, err.Code)
	assert.Equal(t, "test error", err.Message)
	assert.Equal(t, "test-log-id", err.LogID)
}

func TestCozeError_Error(t *testing.T) {
	// 测试 Error() 方法
	err := NewCozeError(1001, "test error", "test-log-id")
	expectedMsg := "Code: 1001, Message: test error, LogID: test-log-id"
	assert.Equal(t, expectedMsg, err.Error())
}

func TestAsCozeError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantErr  *CozeError
		wantBool bool
	}{
		{
			name:     "nil error",
			err:      nil,
			wantErr:  nil,
			wantBool: false,
		},
		{
			name:     "non-CozeError",
			err:      errors.New("standard error"),
			wantErr:  nil,
			wantBool: false,
		},
		{
			name:     "CozeError",
			err:      NewCozeError(1001, "test error", "test-log-id"),
			wantErr:  NewCozeError(1001, "test error", "test-log-id"),
			wantBool: true,
		},
		{
			name: "wrapped CozeError",
			err: fmt.Errorf("wrapped: %w",
				NewCozeError(1001, "test error", "test-log-id")),
			wantErr:  NewCozeError(1001, "test error", "test-log-id"),
			wantBool: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr, gotBool := AsCozeError(tt.err)
			assert.Equal(t, tt.wantBool, gotBool)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.Code, gotErr.Code)
				assert.Equal(t, tt.wantErr.Message, gotErr.Message)
				assert.Equal(t, tt.wantErr.LogID, gotErr.LogID)
			} else {
				assert.Nil(t, gotErr)
			}
		})
	}
}
