package coze

import (
	"context"
	"io"
	"net/http"
)

type audioSpeech struct {
	client *httpClient
}

func newSpeech(client *httpClient) *audioSpeech {
	return &audioSpeech{client: client}
}

func (r *audioSpeech) Create(ctx context.Context, req *CreateAudioSpeechReq) (*CreateAudioSpeechResp, error) {
	uri := "/v1/audio/speech"
	resp, err := r.client.RawRequest(ctx, http.MethodPost, uri, req)
	if err != nil {
		return nil, err
	}
	res := &CreateAudioSpeechResp{
		Data: resp.Body,
	}
	res.SetHTTPResponse(newHTTPResponse(resp))
	return res, nil
}

// CreateAudioSpeechReq represents the request for creating speech
type CreateAudioSpeechReq struct {
	Input          string      `json:"input"`
	VoiceID        string      `json:"voice_id"`
	ResponseFormat AudioFormat `json:"response_format"`
	Speed          float32     `json:"speed"`
}

// CreateAudioSpeechResp represents the response for creating speech
type CreateAudioSpeechResp struct {
	baseResponse
	Data io.ReadCloser
}
