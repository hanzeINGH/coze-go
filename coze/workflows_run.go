package coze

import (
	"bufio"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/coze/coze/internal"
)

// WorkflowRunResult represents the result of a workflow run
type WorkflowRunResult struct {
	DebugUrl string `json:"debug_url"`

	// Workflow execution result, usually a JSON serialized string. In some scenarios, a string with a
	// non-JSON structure may be returned.
	Data string `json:"data"`

	// Execution ID of asynchronous execution. Only returned when the workflow is executed
	// asynchronously (is_async=true). You can use execute_id to call the Query Workflow Asynchronous
	// Execution Result API to obtain the final execution result of the workflow.
	ExecuteID string `json:"execute_id"`
}

// WorkflowEvent represents an event in a workflow
type WorkflowEvent struct {
	// The event ID of this message in the interface response. It starts from 0.
	ID int `json:"id"`

	// The current streaming data packet event.
	Event WorkflowEventType `json:"event"`

	Message   *WorkflowEventMessage   `json:"message,omitempty"`
	Interrupt *WorkflowEventInterrupt `json:"interrupt,omitempty"`
	Error     *WorkflowEventError     `json:"error,omitempty"`
	LogID     string                  `json:"log_id"`
}

func parseWorkflowEventMessage(id int, data string, logID string) (*WorkflowEvent, error) {
	var message WorkflowEventMessage
	if err := json.Unmarshal([]byte(data), &message); err != nil {
		return nil, err
	}

	return &WorkflowEvent{
		ID:      id,
		Event:   WorkflowEventTypeMessage,
		Message: &message,
		LogID:   logID,
	}, nil
}

func parseWorkflowEventInterrupt(id int, data string, logID string) (*WorkflowEvent, error) {
	var interrupt WorkflowEventInterrupt
	if err := json.Unmarshal([]byte(data), &interrupt); err != nil {
		return nil, err
	}

	return &WorkflowEvent{
		ID:        id,
		Event:     WorkflowEventTypeInterrupt,
		Interrupt: &interrupt,
		LogID:     logID,
	}, nil
}

func parseWorkflowEventError(id int, data string, logID string) (*WorkflowEvent, error) {
	var errorEvent WorkflowEventError
	if err := json.Unmarshal([]byte(data), &errorEvent); err != nil {
		return nil, err
	}

	return &WorkflowEvent{
		ID:    id,
		Event: WorkflowEventTypeError,
		Error: &errorEvent,
		LogID: logID,
	}, nil
}

func parseWorkflowEventDone(id int, logID string) *WorkflowEvent {
	return &WorkflowEvent{
		ID:    id,
		Event: WorkflowEventTypeDone,
		LogID: logID,
	}
}

func ParseWorkflowEvent(eventLine map[string]string, logID string) (*WorkflowEvent, error) {
	id, _ := strconv.Atoi(eventLine["id"])
	event := WorkflowEventType(eventLine["event"])
	data := eventLine["data"]

	switch event {
	case WorkflowEventTypeMessage:
		return parseWorkflowEventMessage(id, data, logID)
	case WorkflowEventTypeInterrupt:
		return parseWorkflowEventInterrupt(id, data, logID)
	case WorkflowEventTypeError:
		return parseWorkflowEventError(id, data, logID)
	case WorkflowEventTypeDone:
		return parseWorkflowEventDone(id, logID), nil
	default:
		return parseWorkflowEventMessage(id, data, logID)
	}
}

func (e *WorkflowEvent) IsDone() bool {
	return e.Event == WorkflowEventTypeDone
}

// WorkflowEventError represents an error event in a workflow
type WorkflowEventError struct {
	// Status code. 0 represents a successful API call. Other values indicate that the call has
	// failed. You can determine the detailed reason for the error through the error_message field.
	ErrorCode int `json:"error_code"`

	// Status message. You can get detailed error information when the API call fails.
	ErrorMessage string `json:"error_message"`
}

// WorkflowEventInterrupt represents an interruption event in a workflow
type WorkflowEventInterrupt struct {
	// The content of interruption event.
	InterruptData *WorkflowEventInterruptData `json:"interrupt_data"`

	// The name of the node that outputs the message, such as "Question".
	NodeTitle string `json:"node_title"`
}

// WorkflowEventInterruptData represents the data of an interruption event
type WorkflowEventInterruptData struct {
	// The workflow interruption event ID, which should be passed back when resuming the workflow.
	EventID string `json:"event_id"`

	// The type of workflow interruption, which should be passed back when resuming the workflow.
	Type int `json:"type"`
}

// ParseWorkflowEventError parses JSON string to WorkflowEventError
func ParseWorkflowEventError(data string) (*WorkflowEventError, error) {
	var err WorkflowEventError
	if err := json.Unmarshal([]byte(data), &err); err != nil {
		return nil, err
	}
	return &err, nil
}

// ParseWorkflowEventInterrupt parses JSON string to WorkflowEventInterrupt
func ParseWorkflowEventInterrupt(data string) (*WorkflowEventInterrupt, error) {
	var interrupt WorkflowEventInterrupt
	if err := json.Unmarshal([]byte(data), &interrupt); err != nil {
		return nil, err
	}
	return &interrupt, nil
}

// WorkflowEventMessage represents a message event in a workflow
type WorkflowEventMessage struct {
	// The content of the streamed output message.
	Content string `json:"content"`

	// The name of the node that outputs the message, such as the message node or end node.
	NodeTitle string `json:"node_title"`

	// The message ID of this message within the node, starting at 0.
	NodeSeqID string `json:"node_seq_id"`

	// Whether the current message is the last data packet for this node.
	NodeIsFinish bool `json:"node_is_finish"`

	// Additional fields.
	Ext map[string]any `json:"ext,omitempty"`
}

// WorkflowEventType represents the type of workflow event
type WorkflowEventType string

const (
	// The output message from the workflow node, such as the output message from the message node or
	// end node. You can view the specific message content in data.
	WorkflowEventTypeMessage WorkflowEventType = "Message"

	// An error has occurred. You can view the error_code and error_message in data to troubleshoot
	// the issue.
	WorkflowEventTypeError WorkflowEventType = "Error"

	// End. Indicates the end of the workflow execution, where data is empty.
	WorkflowEventTypeDone WorkflowEventType = "Done"

	// Interruption. Indicates the workflow has been interrupted, where the data field contains
	// specific interruption information.
	WorkflowEventTypeInterrupt WorkflowEventType = "Interrupt"
)

// RunWorkflowReq represents request for running workflow
type RunWorkflowReq struct {
	// The ID of the workflow, which should have been published.
	WorkflowID string `json:"workflow_id"`

	// Input parameters and their values for the starting node of the workflow.
	Parameters map[string]any `json:"parameters,omitempty"`

	// The associated Bot ID required for some workflow executions.
	BotID string `json:"bot_id,omitempty"`

	// Used to specify some additional fields.
	Ext map[string]string `json:"ext,omitempty"`

	// Whether to run asynchronously.
	IsAsync bool `json:"is_async,omitempty"`
}

// ResumeRunReq represents request for resuming workflow run
type ResumeRunReq struct {
	// The ID of the workflow, which should have been published.
	WorkflowID string `json:"workflow_id"`

	// Event ID
	EventID string `json:"event_id"`

	// Resume data
	ResumeData string `json:"resume_data"`

	// Interrupt type
	InterruptType int `json:"interrupt_type"`
}

type workflowRun struct {
	client    *internal.Client
	Histories *workflowRunHistories
}

func newWorkflowRun(client *internal.Client) *workflowRun {
	return &workflowRun{
		client:    client,
		Histories: newWorkflowRunHistories(client),
	}
}

func (r *workflowRun) Run(ctx context.Context, req RunWorkflowReq) (*RunWorkflowResp, error) {
	method := http.MethodPost
	uri := "/v1/workflows/run"
	resp := &RunWorkflowResp{}
	err := r.client.Request(ctx, method, uri, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

type WorkflowEventReader struct {
	*streamReader[WorkflowEvent]
}

func (r *workflowRun) Resume(ctx context.Context, req ResumeRunReq) (*WorkflowEventReader, error) {
	method := http.MethodPost
	uri := "/v1/workflow/stream_resume"
	resp, err := r.client.RowRequest(ctx, method, uri, req)
	if err != nil {
		return nil, err
	}

	return &WorkflowEventReader{
		streamReader: &streamReader[WorkflowEvent]{
			response:  resp,
			reader:    bufio.NewReader(resp.Body),
			logID:     internal.GetLogID(resp.Header),
			processor: parseWorkflowEvent,
		},
	}, nil
}

func (r *workflowRun) Stream(ctx context.Context, req RunWorkflowReq) (*WorkflowEventReader, error) {
	method := http.MethodPost
	uri := "/v1/workflow/stream_run"
	resp, err := r.client.RowRequest(ctx, method, uri, req)
	if err != nil {
		return nil, err
	}

	return &WorkflowEventReader{
		streamReader: &streamReader[WorkflowEvent]{
			response:  resp,
			reader:    bufio.NewReader(resp.Body),
			logID:     internal.GetLogID(resp.Header),
			processor: parseWorkflowEvent,
		},
	}, nil
}

func parseWorkflowEvent(lineBytes []byte, reader *bufio.Reader, logID string) (*WorkflowEvent, bool, error) {
	line := string(lineBytes)
	if strings.HasPrefix(line, "id:") {
		id := strings.TrimSpace(line[3:])
		data, err := reader.ReadString('\n')
		if err != nil {
			return nil, false, err
		}
		event := strings.TrimSpace(data[6:])
		data, err = reader.ReadString('\n')
		if err != nil {
			return nil, false, err
		}
		data = strings.TrimSpace(data[5:])

		eventLine := map[string]string{
			"id":    id,
			"event": event,
			"data":  data,
		}

		eventData, err := ParseWorkflowEvent(eventLine, logID)
		if err != nil {
			return nil, false, err
		}

		return eventData, eventData.IsDone(), nil
	}
	return nil, false, nil
}
