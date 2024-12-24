package internal

type BaseResponse struct {
	LogID string
}

func (b *BaseResponse) SetLogID(logID string) {
	b.LogID = logID
}

type BaseResp interface {
	SetLogID(logID string)
}
