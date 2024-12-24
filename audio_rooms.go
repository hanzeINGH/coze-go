package coze

import (
	"context"
	"net/http"

	"github.com/coze-dev/coze-go/internal"
)

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

// AudioCodec represents the audio codec
type AudioCodec string

const (
	AudioCodecAACLC AudioCodec = "AACLC"
	AudioCodecG711A AudioCodec = "G711A"
	AudioCodecOPUS  AudioCodec = "OPUS"
	AudioCodecG722  AudioCodec = "G722"
)

// RoomAudioConfig represents the room audio configuration
type RoomAudioConfig struct {
	Codec AudioCodec `json:"codec"`
}

// RoomConfig represents the room configuration
type RoomConfig struct {
	AudioConfig *RoomAudioConfig `json:"audio_config"`
}

// CreateAudioRoomsReq represents the request for creating an audio room
type CreateAudioRoomsReq struct {
	BotID          string      `json:"bot_id"`
	ConversationID string      `json:"conversation_id,omitempty"`
	VoiceID        string      `json:"voice_id,omitempty"`
	Config         *RoomConfig `json:"config,omitempty"`
}

// CreateAudioRoomsResp represents the response for creating an audio room
type CreateAudioRoomsResp struct {
	internal.BaseResponse
	Data struct {
		RoomID string `json:"room_id"`
		AppID  string `json:"app_id"`
		Token  string `json:"token"`
		UID    string `json:"uid"`
	} `json:"data"`
}
