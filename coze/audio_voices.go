package coze

import (
	"context"
	"net/http"
	"os"
	"strconv"

	"github.com/coze-dev/coze-go/coze/internal"
	"github.com/coze-dev/coze-go/coze/pagination"
)

// Voice 语音模型
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

// CloneVoiceReq 克隆语音请求
type CloneVoiceReq struct {
	VoiceName   string        `json:"voice_name"`
	FilePath    string        `json:"file_path"`
	AudioFormat *AudioFormat  `json:"audio_format"`
	Language    *LanguageCode `json:"language,omitempty"`
	VoiceID     *string       `json:"voice_id,omitempty"`
	PreviewText *string       `json:"preview_text,omitempty"`
	Text        *string       `json:"text,omitempty"`
}

// CloneVoiceResp 克隆语音响应
type CloneVoiceResp struct {
	internal.BaseResponse
	Data struct {
		VoiceID string `json:"voice_id"`
	} `json:"data"`
}

// ListVoiceReq 列出语音请求
type ListVoiceReq struct {
	FilterSystemVoice bool `json:"filter_system_voice,omitempty"`
	PageNum           int  `json:"page_num"`
	PageSize          int  `json:"page_size"`
}

// ListVoiceResp 列出语音响应
type ListVoiceResp struct {
	internal.BaseResponse
	Data struct {
		VoiceList []Voice `json:"voice_list"`
	} `json:"data"`
}

type voice struct {
	client *internal.Client
}

func newVoice(client *internal.Client) *voice {
	return &voice{client: client}
}

func (r *voice) Clone(ctx context.Context, req CloneVoiceReq) (*CloneVoiceResp, error) {
	path := "/v1/audio/voices/clone"
	file, err := os.Open(req.FilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fields := map[string]string{
		"voiceName":   req.VoiceName,
		"audioFormat": req.AudioFormat.String(),
	}

	// 添加其他字段
	if req.Language != nil {
		fields["language"] = req.Language.String()
	}
	if req.VoiceID != nil {
		fields["voiceID"] = *req.VoiceID
	}
	if req.PreviewText != nil {
		fields["previewText"] = *req.PreviewText
	}
	if req.Text != nil {
		fields["text"] = *req.Text
	}
	resp := &CloneVoiceResp{}
	err = r.client.UploadFile(ctx, path, file, file.Name(), fields, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *voice) List(ctx context.Context, req ListVoiceReq) (*pagination.PageNumBasedPager[Voice], error) {
	if req.PageSize == 0 {
		req.PageSize = 20
	}
	if req.PageNum == 0 {
		req.PageNum = 1
	}
	return pagination.NewPageNumBasedPager[Voice](
		func(request *pagination.PageRequest) (*pagination.PageResponse[Voice], error) {
			uri := "/v1/audio/voices"
			resp := &ListVoiceResp{}
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
