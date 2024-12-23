package internal

type BaseResponse struct {
	Code  int    `json:"code"`
	Msg   string `json:"msg"`
	LogID string
}

func (b *BaseResponse) SetLogID(logID string) {
	b.LogID = logID
}

func (b *BaseResponse) SetCode(code int) {
	b.Code = code
}

func (b *BaseResponse) SetMsg(msg string) {
	b.Msg = msg
}

func (b *BaseResponse) GetCode() int {
	return b.Code
}

func (b *BaseResponse) GetMsg() string {
	return b.Msg
}

type BaseResp interface {
	SetLogID(logID string)
	SetCode(code int)
	SetMsg(msg string)
	GetMsg() string
	GetCode() int
}
