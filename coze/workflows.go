package coze

import "github.com/coze-dev/coze-go/coze/internal"

type workflows struct {
	Run *workflowRun
}

func newWorkflows(client *internal.Client) *workflows {
	return &workflows{
		Run: newWorkflowRun(client),
	}
}
