package coze

type workflows struct {
	Runs *workflowRuns
}

func newWorkflows(client *httpClient) *workflows {
	return &workflows{
		Runs: newWorkflowRun(client),
	}
}
