package coze

import "github.com/coze-dev/coze-go/coze/internal"

type dataset struct {
	Documents *documents
}

func newDataset(client *internal.Client) *dataset {
	return &dataset{
		Documents: newDocuments(client),
	}
}
