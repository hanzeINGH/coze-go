package coze

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/coze-dev/coze-go/internal"
	"github.com/coze-dev/coze-go/pagination"
)

type audioVoice struct {
	client *internal.Client
}

func newVoice(client *internal.Client) *audioVoice {
	return &audioVoice{client: client}
}

func (r *audioVoice) Clone(ctx context.Context, req *CloneAudioVoicesReq) (*CloneAudioVoicesResp, error) {
	path := "/v1/audio/voices/clone"
	if req.File == nil {
		return nil, fmt.Errorf("file is required")
	}

	fields := map[string]string{
		"voice_name":   req.VoiceName,
		"audio_format": req.AudioFormat.String(),
	}

	// Add other fields
	if req.Language != nil {
		fields["language"] = req.Language.String()
	}
	if req.VoiceID != nil {
		fields["voice_id"] = *req.VoiceID
	}
	if req.PreviewText != nil {
		fields["preview_text"] = *req.PreviewText
	}
	if req.Text != nil {
		fields["text"] = *req.Text
	}
	resp := &CloneAudioVoicesResp{}
	err := r.client.UploadFile(ctx, path, req.File, req.VoiceName, fields, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *audioVoice) List(ctx context.Context, req *ListAudioVoicesReq) (*pagination.NumberPaged[Voice], error) {
	if req.PageSize == 0 {
		req.PageSize = 20
	}
	if req.PageNum == 0 {
		req.PageNum = 1
	}
	return pagination.NewNumberPaged[Voice](
		func(request *pagination.PageRequest) (*pagination.PageResponse[Voice], error) {
			uri := "/v1/audio/voices"
			resp := &ListAudioVoicesResp{}
			err := r.client.Request(ctx, http.MethodGet, uri, nil, resp,
				internal.WithQuery("page_num", strconv.Itoa(request.PageNum)),
				internal.WithQuery("page_size", strconv.Itoa(request.PageSize)),
				internal.WithQuery("filter_system_voice", strconv.FormatBool(req.FilterSystemVoice)))
			if err != nil {
				return nil, err
			}
			return &pagination.PageResponse[Voice]{
				HasMore: len(resp.Data.VoiceList) >= request.PageSize,
				Data:    resp.Data.VoiceList,
				LogID:   resp.LogID,
			}, nil
		}, req.PageSize, req.PageNum)
}

// Voice represents the voice model
type Voice struct {
	VoiceID                string `json:"voice_id"`
	Name                   string `json:"name"`
	IsSystemVoice          bool   `json:"is_system_voice"`
	LanguageCode           string `json:"language_code"`
	LanguageName           string `json:"language_name"`
	PreviewText            string `json:"preview_text"`
	PreviewAudio           string `json:"preview_audio"`
	AvailableTrainingTimes int    `json:"available_training_times"`
	CreateTime             int    `json:"create_time"`
	UpdateTime             int    `json:"update_time"`
}

// CloneAudioVoicesReq represents the request for cloning a voice
type CloneAudioVoicesReq struct {
	VoiceName   string
	File        io.Reader
	AudioFormat *AudioFormat
	Language    *LanguageCode
	VoiceID     *string
	PreviewText *string
	Text        *string
}

// CloneAudioVoicesResp represents the response for cloning a voice
type CloneAudioVoicesResp struct {
	internal.BaseResponse
	Data struct {
		VoiceID string `json:"voice_id"`
	} `json:"data"`
}

// ListAudioVoicesReq represents the request for listing voices
type ListAudioVoicesReq struct {
	FilterSystemVoice bool `json:"filter_system_voice,omitempty"`
	PageNum           int  `json:"page_num"`
	PageSize          int  `json:"page_size"`
}

// ListAudioVoicesResp represents the response for listing voices
type ListAudioVoicesResp struct {
	internal.BaseResponse
	Data struct {
		VoiceList []*Voice `json:"voice_list"`
	} `json:"data"`
}
