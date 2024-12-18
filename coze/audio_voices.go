package coze

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
	VoiceName   string       `json:"voice_name"`
	FilePath    string       `json:"file_path"`
	AudioFormat AudioFormat  `json:"audio_format"`
	Language    LanguageCode `json:"language,omitempty"`
	VoiceID     string       `json:"voice_id,omitempty"`
	PreviewText string       `json:"preview_text,omitempty"`
	Text        string       `json:"text,omitempty"`
}

// CloneVoiceResp 克隆语音响应
type CloneVoiceResp struct {
	VoiceID string `json:"voice_id"`
}

// ListVoiceReq 列出语音请求
type ListVoiceReq struct {
	FilterSystemVoice bool `json:"filter_system_voice,omitempty"`
	PageNum           int  `json:"page_num"`
	PageSize          int  `json:"page_size"`
}

// ListVoiceResp 列出语音响应
type ListVoiceResp struct {
	VoiceList []Voice `json:"voice_list"`
}

type voice struct {
}
