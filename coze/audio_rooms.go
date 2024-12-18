package coze

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

// CreateRoomReq 创建房间请求
type CreateRoomReq struct {
	BotID          string      `json:"bot_id"`
	ConversationID string      `json:"conversation_id,omitempty"`
	VoiceID        string      `json:"voice_id,omitempty"`
	Config         *RoomConfig `json:"config,omitempty"`
}

// CreateRoomResp 创建房间响应
type CreateRoomResp struct {
	RoomID string `json:"room_id"`
	AppID  string `json:"app_id"`
	Token  string `json:"token"`
	UID    string `json:"uid"`
}

type rooms struct {
}
