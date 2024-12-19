package coze

import "github.com/coze/coze/internal"

type workflows struct {
	Run *workflowRun
}

func newWorkflows(client *internal.Client) *workflows {
	return &workflows{
		Run: newWorkflowRun(client),
	}
}
