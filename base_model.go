package coze

type baseResponse struct {
	Code  int    `json:"code"`
	Msg   string `json:"msg"`
	LogID string
}

func (b *baseResponse) SetLogID(logID string) {
	b.LogID = logID
}

func (b *baseResponse) SetCode(code int) {
	b.Code = code
}

func (b *baseResponse) SetMsg(msg string) {
	b.Msg = msg
}

func (b *baseResponse) GetCode() int {
	return b.Code
}

func (b *baseResponse) GetMsg() string {
	return b.Msg
}

type baseRespInterface interface {
	SetLogID(logID string)
	SetCode(code int)
	SetMsg(msg string)
	GetMsg() string
	GetCode() int
}

type baseModel struct {
	LogID string `json:"log_id"`
}

func (b *baseModel) SetLogID(logID string) {
	b.LogID = logID
}
