package coze

type datasets struct {
	Documents *datasetsDocuments
}

func newDatasets(client *httpClient) *datasets {
	return &datasets{
		Documents: newDocuments(client),
	}
}
