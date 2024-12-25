package coze

import "net/http"

type HTTPResponse interface {
	LogID() string
}

type httpResponse struct {
	Status        int
	Header        http.Header
	ContentLength int64

	logid string
}

func (r *httpResponse) LogID() string {
	if r.logid == "" {
		r.logid = r.Header.Get(logIDHeader)
	}
	return r.logid
}

type baseResponse struct {
	Code         int           `json:"code"`
	Msg          string        `json:"msg"`
	HTTPResponse *httpResponse `json:"http_response"`
}

func (b *baseResponse) SetHTTPResponse(httpResponse *httpResponse) {
	b.HTTPResponse = httpResponse
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
	SetHTTPResponse(httpResponse *httpResponse)
	SetCode(code int)
	SetMsg(msg string)
	GetMsg() string
	GetCode() int
}

type baseModel struct {
	httpResponse *httpResponse
}

func (b *baseModel) setHTTPResponse(httpResponse *httpResponse) {
	b.httpResponse = httpResponse
}

func (b *baseModel) HTTPResponse() HTTPResponse {
	return b.httpResponse
}

func (b *baseModel) LogID() string {
	return b.httpResponse.LogID()
}

func newHTTPResponse(resp *http.Response) *httpResponse {
	return &httpResponse{
		Status:        resp.StatusCode,
		Header:        resp.Header,
		ContentLength: resp.ContentLength,
	}
}
