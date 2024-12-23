package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/coze-dev/coze-go"
)

func main() {
	// Get an access_token through personal access token or oauth.
	token := os.Getenv("COZE_API_TOKEN")
	authCli := coze.NewTokenAuth(token)

	// Init the Coze client through the access_token.
	cozeCli := coze.NewCozeAPI(authCli, coze.WithBaseURL(os.Getenv("COZE_API_BASE")))

	ctx := context.Background()
	workflowID := os.Getenv("WORKFLOW_ID")

	// if your workflow need input params, you can send them by map
	data := map[string]interface{}{
		"date": "param values",
	}

	req := &coze.RunWorkflowsReq{
		WorkflowID: workflowID,
		Parameters: data,
		IsAsync:    true, // if you want the workflow run asynchronously, you must set isAsync to true.
	}

	resp, err := cozeCli.Workflows.Runs.Create(ctx, req)
	if err != nil {
		fmt.Println("Error running workflow:", err)
		return
	}
	fmt.Println("Start async workflow run:", resp.ExecuteID)

	executeID := resp.ExecuteID
	isFinished := false

	for !isFinished {
		historyResp, err := cozeCli.Workflows.Runs.Histories.Retrieve(ctx, &coze.RetrieveWorkflowsRunHistoriesReq{
			WorkflowID: workflowID,
			ExecuteID:  executeID,
		})
		if err != nil {
			fmt.Println("Error retrieving history:", err)
			return
		}
		fmt.Println(historyResp)

		history := historyResp.Histories[0]
		switch history.ExecuteStatus {
		case coze.WorkflowExecuteStatusFail:
			fmt.Println("Workflow run failed, reason:", history.ErrorMessage)
			isFinished = true
		case coze.WorkflowExecuteStatusRunning:
			fmt.Println("Workflow run is running")
			time.Sleep(time.Second)
		default:
			fmt.Println("Workflow run success:", history.Output)
			isFinished = true
		}
	}
}
