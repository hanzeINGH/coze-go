package coze

import (
	"github.com/coze-dev/coze-go/internal"
)

type workflows struct {
	Runs *workflowRuns
}

func newWorkflows(client *internal.Client) *workflows {
	return &workflows{
		Runs: newWorkflowRun(client),
	}
}
