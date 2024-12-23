package coze

import (
	"context"
	"net/http"

	"github.com/coze-dev/coze-go/internal"
)

// AudioCodec 音频编解码器
type AudioCodec string

const (
	AudioCodecAACLC AudioCodec = "AACLC"
	AudioCodecG711A AudioCodec = "G711A"
	AudioCodecOPUS  AudioCodec = "OPUS"
	AudioCodecG722  AudioCodec = "G722"
)

// RoomAudioConfig 房间音频配置
type RoomAudioConfig struct {
	Codec AudioCodec `json:"codec"`
}

// RoomConfig 房间配置
type RoomConfig struct {
	AudioConfig *RoomAudioConfig `json:"audio_config"`
}

// CreateAudioRoomsReq 创建房间请求
type CreateAudioRoomsReq struct {
	BotID          string      `json:"bot_id"`
	ConversationID string      `json:"conversation_id,omitempty"`
	VoiceID        string      `json:"voice_id,omitempty"`
	Config         *RoomConfig `json:"config,omitempty"`
}

// CreateAudioRoomsResp 创建房间响应
type CreateAudioRoomsResp struct {
	internal.BaseResponse
	Data struct {
		RoomID string `json:"room_id"`
		AppID  string `json:"app_id"`
		Token  string `json:"token"`
		UID    string `json:"uid"`
	} `json:"data"`
}

type audioRooms struct {
	client *internal.Client
}

func newRooms(client *internal.Client) *audioRooms {
	return &audioRooms{client: client}
}

func (r *audioRooms) Create(ctx context.Context, req *CreateAudioRoomsReq) (*CreateAudioRoomsResp, error) {
	method := http.MethodPost
	uri := "/v1/audio/audioRooms"
	resp := &CreateAudioRoomsResp{}
	err := r.client.Request(ctx, method, uri, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
