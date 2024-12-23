package main

import (
	"context"
	"fmt"
	"os"

	"github.com/coze-dev/coze-go"
)

func main() {
	// Get an access_token through personal access token or oauth.
	token := os.Getenv("COZE_API_TOKEN")
	authCli := coze.NewTokenAuth(token)

	// Init the Coze client through the access_token.
	cozeCli := coze.NewCozeAPI(authCli, coze.WithBaseURL(os.Getenv("COZE_API_BASE")))

	file, err := os.Open(os.Getenv("VOICE_FILE_PATH"))
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	ctx := context.Background()
	resp, err := cozeCli.Audio.Voices.Clone(ctx, &coze.CloneAudioVoicesReq{
		File:        file,
		VoiceName:   "your voice name",
		AudioFormat: coze.AudioFormatM4A.Ptr(),
	})
	if err != nil {
		fmt.Println("Error cloning voice:", err)
		return
	}

	fmt.Println(resp)
}
