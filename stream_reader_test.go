package coze

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock event processor for testing
func mockEventProcessor(line []byte, reader *bufio.Reader) (*WorkflowEvent, bool, error) {
	if len(line) == 0 {
		return nil, false, nil
	}

	// Parse event data
	event := &WorkflowEvent{
		ID:    0,
		Event: WorkflowEventTypeMessage,
		Message: &WorkflowEventMessage{
			Content: string(line),
		},
	}

	// Check if this is the last event
	isDone := string(line) == "done"
	if isDone {
		event.Event = WorkflowEventTypeDone
	}
	return event, isDone, nil
}

func TestStreamReader(t *testing.T) {
	t.Run("successful event processing", func(t *testing.T) {
		// Create mock response with multiple events
		events := []string{
			"first",
			"second",
			"done",
		}
		resp := createMockResponse(events)

		// Create stream reader
		reader := &streamReader[WorkflowEvent]{
			reader:    bufio.NewReader(resp.Body),
			response:  resp,
			logID:     "test_log_id",
			processor: mockEventProcessor,
		}
		defer reader.Close()

		// Read first event
		event, err := reader.Recv()
		require.NoError(t, err)
		assert.Equal(t, WorkflowEventTypeMessage, event.Event)
		assert.Equal(t, "first", event.Message.Content)
		assert.False(t, reader.isFinished)

		// Read second event
		event, err = reader.Recv()
		require.NoError(t, err)
		assert.Equal(t, WorkflowEventTypeMessage, event.Event)
		assert.Equal(t, "second", event.Message.Content)
		assert.False(t, reader.isFinished)

		// Read final event
		event, err = reader.Recv()
		require.NoError(t, err)
		assert.Equal(t, WorkflowEventTypeDone, event.Event)
		assert.True(t, reader.isFinished)

		// Try reading after done
		event, err = reader.Recv()
		assert.Equal(t, io.EOF, err)
		assert.Nil(t, event)
	})

	t.Run("empty lines are skipped", func(t *testing.T) {
		events := []string{
			"",
			"test",
			"",
			"done",
		}
		resp := createMockResponse(events)

		reader := &streamReader[WorkflowEvent]{
			reader:    bufio.NewReader(resp.Body),
			response:  resp,
			logID:     "test_log_id",
			processor: mockEventProcessor,
		}
		defer reader.Close()

		// First non-empty event
		event, err := reader.Recv()
		require.NoError(t, err)
		assert.Equal(t, WorkflowEventTypeMessage, event.Event)
		assert.Equal(t, "test", event.Message.Content)

		// Second non-empty event
		event, err = reader.Recv()
		require.NoError(t, err)
		assert.Equal(t, WorkflowEventTypeDone, event.Event)
	})

	t.Run("error response handling", func(t *testing.T) {
		// Create mock error response
		errorResp := &http.Response{
			StatusCode: http.StatusBadRequest,
			Header: http.Header{
				"Content-Type": []string{"application/json"},
			},
			Body: io.NopCloser(strings.NewReader(`{
				"log_id": "error_log_id",
				"error": {
					"code": 400,
					"message": "Bad Request"
				}
			}`)),
		}

		reader := &streamReader[WorkflowEvent]{
			reader:    bufio.NewReader(errorResp.Body),
			response:  errorResp,
			logID:     "error_log_id",
			processor: mockEventProcessor,
		}
		defer reader.Close()

		// Attempt to read should return error
		event, err := reader.Recv()
		assert.Error(t, err)
		assert.Nil(t, event)
	})

	t.Run("LogID method", func(t *testing.T) {
		reader := &streamReader[WorkflowEvent]{
			logID: "test_log_id",
		}
		assert.Equal(t, "test_log_id", reader.LogID())
	})
}

// Helper function to create mock response with events
func createMockResponse(events []string) *http.Response {
	// Join events with newlines
	body := strings.Join(events, "\n")

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(body)),
	}
}