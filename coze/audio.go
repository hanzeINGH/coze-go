package coze

import "github.com/coze-dev/coze-go/coze/internal"

// AudioFormat 音频格式类型
type AudioFormat string

const (
	AudioFormatWAV     AudioFormat = "wav"
	AudioFormatPCM     AudioFormat = "pcm"
	AudioFormatOGGOPUS AudioFormat = "ogg_opus"
	AudioFormatM4A     AudioFormat = "m4a"
	AudioFormatAAC     AudioFormat = "aac"
	AudioFormatMP3     AudioFormat = "mp3"
)

func (f AudioFormat) String() string {
	return string(f)
}

// LanguageCode 语言代码
type LanguageCode string

const (
	LanguageCodeZH LanguageCode = "zh"
	LanguageCodeEN LanguageCode = "en"
	LanguageCodeJA LanguageCode = "ja"
	LanguageCodeES LanguageCode = "es"
	LanguageCodeID LanguageCode = "id"
	LanguageCodePT LanguageCode = "pt"
)

func (l LanguageCode) String() string {
	return string(l)
}

type audio struct {
	Rooms  *rooms
	Speech *speech
	Voice  *voice
}

func newAudio(client *internal.Client) *audio {
	return &audio{
		Rooms:  newRooms(client),
		Speech: newSpeech(client),
		Voice:  newVoice(client),
	}
}
