package coze

// CreateSpeechReq 创建语音请求
type CreateSpeechReq struct {
	Input          string      `json:"input"`
	VoiceID        string      `json:"voice_id"`
	ResponseFormat AudioFormat `json:"response_format"`
	Speed          float32     `json:"speed"`
}

// CreateSpeechResp 创建语音响应
type CreateSpeechResp struct {
	Response []byte
}

type speech struct {
}
