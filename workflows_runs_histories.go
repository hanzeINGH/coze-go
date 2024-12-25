package coze

import (
	"context"
	"fmt"
	"net/http"
)

type workflowRunHistories struct {
	client *httpClient
}

func newWorkflowRunHistories(client *httpClient) *workflowRunHistories {
	return &workflowRunHistories{client: client}
}

func (r *workflowRunHistories) Retrieve(ctx context.Context, req *RetrieveWorkflowsRunHistoriesReq) (*RetrieveWorkflowRunHistoriesResp, error) {
	method := http.MethodGet
	uri := fmt.Sprintf("/v1/workflows/%s/run_histories/%s", req.WorkflowID, req.ExecuteID)
	resp := &retrieveWorkflowRunHistoriesResp{}
	err := r.client.Request(ctx, method, uri, nil, resp)
	if err != nil {
		return nil, err
	}
	resp.RetrieveWorkflowRunHistoriesResp.setHTTPResponse(resp.httpResponse)
	return resp.RetrieveWorkflowRunHistoriesResp, nil
}

// WorkflowRunMode represents how the workflow runs
type WorkflowRunMode int

const (
	// Synchronous operation.
	WorkflowRunModeSynchronous WorkflowRunMode = 0

	// Streaming operation.
	WorkflowRunModeStreaming WorkflowRunMode = 1

	// Asynchronous operation.
	WorkflowRunModeAsynchronous WorkflowRunMode = 2
)

// WorkflowExecuteStatus represents the execution status of a workflow
type WorkflowExecuteStatus string

const (
	// Execution succeeded.
	WorkflowExecuteStatusSuccess WorkflowExecuteStatus = "Success"

	// Execution in progress.
	WorkflowExecuteStatusRunning WorkflowExecuteStatus = "Running"

	// Execution failed.
	WorkflowExecuteStatusFail WorkflowExecuteStatus = "Fail"
)

// RetrieveWorkflowsRunHistoriesReq represents request for retrieving workflow runs history
type RetrieveWorkflowsRunHistoriesReq struct {
	// The ID of the workflow.
	ExecuteID string `json:"execute_id"`

	// The ID of the workflow async execute.
	WorkflowID string `json:"workflow_id"`
}

// runWorkflowsResp represents response for running workflow
type runWorkflowsResp struct {
	baseResponse
	*RunWorkflowsResp
}

// RunWorkflowsResp represents response for running workflow
type RunWorkflowsResp struct {
	baseModel
	// Execution ID of asynchronous execution.
	ExecuteID string `json:"execute_id,omitempty"`

	// Workflow execution result.
	Data string `json:"data,omitempty"`

	DebugURL string `json:"debug_url,omitempty"`
	Token    int    `json:"token,omitempty"`
	Cost     string `json:"cost,omitempty"`
}

// retrieveWorkflowRunHistoriesResp represents response for retrieving workflow runs history
type retrieveWorkflowRunHistoriesResp struct {
	baseResponse
	*RetrieveWorkflowRunHistoriesResp
}

// RetrieveWorkflowRunHistoriesResp represents response for retrieving workflow runs history
type RetrieveWorkflowRunHistoriesResp struct {
	baseModel
	Histories []*WorkflowRunHistory `json:"data"`
}

// WorkflowRunHistory represents the history of a workflow runs
type WorkflowRunHistory struct {
	// The ID of execute.
	ExecuteID string `json:"execute_id"`

	// Execute status: success: Execution succeeded. running: Execution in progress. fail: Execution failed.
	ExecuteStatus WorkflowExecuteStatus `json:"execute_status"`

	// The Bot ID specified when executing the workflow. Returns 0 if no Bot ID was specified.
	BotID string `json:"bot_id"`

	// The release connector ID of the agent. By default, only the Agent as API connector is
	// displayed, and the connector ID is 1024.
	ConnectorID string `json:"connector_id"`

	// User ID, the user_id specified by the ext field when executing the workflow. If not specified,
	// the token applicant's button ID is returned.
	ConnectorUid string `json:"connector_uid"`

	// How the workflow runs: 0: Synchronous operation. 1: Streaming operation. 2: Asynchronous operation.
	RunMode WorkflowRunMode `json:"run_mode"`

	// The Log ID of the asynchronously running workflow. If the workflow is executed abnormally, you
	// can contact the service team to troubleshoot the problem through the Log ID.
	LogID string `json:"logid"`

	// The start time of the workflow, in Unix time timestamp format, in seconds.
	CreateTime int `json:"create_time"`

	// The workflow resume running time, in Unix time timestamp format, in seconds.
	UpdateTime int `json:"update_time"`

	// The output of the workflow is usually a JSON serialized string, but it may also be a non-JSON
	// structured string.
	Output string `json:"output"`

	// Status code. 0 represents a successful API call. Other values indicate that the call has
	// failed. You can determine the detailed reason for the error through the error_message field.
	ErrorCode string `json:"error_code"`

	// Status message. You can get detailed error information when the API call fails.
	ErrorMessage string `json:"error_message"`

	// Workflow trial runs debugging page. Visit this page to view the running results, input and
	// output information of each workflow node.
	DebugUrl string `json:"debug_url"`
}
