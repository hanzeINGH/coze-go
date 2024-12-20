package coze

import (
	"context"
	"io"
	"net/http"

	"github.com/coze-dev/coze-go/coze/internal"
)

// CreateSpeechReq 创建语音请求
type CreateSpeechReq struct {
	Input          string      `json:"input"`
	VoiceID        string      `json:"voice_id"`
	ResponseFormat AudioFormat `json:"response_format"`
	Speed          float32     `json:"speed"`
}

// CreateSpeechResp 创建语音响应
type CreateSpeechResp struct {
	internal.BaseResponse
	Data io.ReadCloser
}

type speech struct {
	client *internal.Client
}

func newSpeech(client *internal.Client) *speech {
	return &speech{client: client}
}

func (r *speech) Create(ctx context.Context, req CreateSpeechReq) (*CreateSpeechResp, error) {
	uri := "/v1/audio/speech"
	resp, err := r.client.RowRequest(ctx, http.MethodPost, uri, req)
	if err != nil {
		return nil, err
	}
	logID := internal.GetLogID(resp.Header)

	return &CreateSpeechResp{
		BaseResponse: internal.BaseResponse{LogID: logID},
		Data:         resp.Body,
	}, nil
}
