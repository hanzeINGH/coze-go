package coze

import (
	"context"
	"io"
	"net/http"

	"github.com/coze-dev/coze-go/internal"
)

type audioSpeech struct {
	client *internal.Client
}

func newSpeech(client *internal.Client) *audioSpeech {
	return &audioSpeech{client: client}
}

func (r *audioSpeech) Create(ctx context.Context, req *CreateAudioSpeechReq) (*CreateAudioSpeechResp, error) {
	uri := "/v1/audio/speech"
	resp, err := r.client.RawRequest(ctx, http.MethodPost, uri, req)
	if err != nil {
		return nil, err
	}
	logID := internal.GetLogID(resp.Header)

	return &CreateAudioSpeechResp{
		BaseResponse: internal.BaseResponse{LogID: logID},
		Data:         resp.Body,
	}, nil
}

// CreateAudioSpeechReq 创建语音请求
type CreateAudioSpeechReq struct {
	Input          string      `json:"input"`
	VoiceID        string      `json:"voice_id"`
	ResponseFormat AudioFormat `json:"response_format"`
	Speed          float32     `json:"speed"`
}

// CreateAudioSpeechResp 创建语音响应
type CreateAudioSpeechResp struct {
	internal.BaseResponse
	Data io.ReadCloser
}
