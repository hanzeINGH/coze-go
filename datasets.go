package coze

import (
	"github.com/coze-dev/coze-go/internal"
)

type datasets struct {
	Documents *datasetsDocuments
}

func newDatasets(client *internal.Client) *datasets {
	return &datasets{
		Documents: newDocuments(client),
	}
}
