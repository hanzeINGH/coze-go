package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/coze/coze"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get an access_token through personal access token or oauth.
	token := os.Getenv("COZE_API_TOKEN")
	botID := os.Getenv("PUBLISHED_BOT_ID")
	userID := os.Getenv("USER_ID")

	authCli := coze.NewTokenAuth(token)

	// Init the Coze client through the access_token.
	cozeCli := coze.NewCozeAPI(authCli, coze.WithBaseURL(os.Getenv("COZE_API_BASE")))

	/*
	 * Step one, create chat
	 * Call the coze.chat().stream() method to create a chat. The create method is a streaming
	 * chat and will return a Flowable ChatEvent. Developers should iterate the iterator to get
	 * chat event and handle them.
	 * */
	req := coze.CreateChatReq{
		BotID:  botID,
		UserID: userID,
		Messages: []coze.Message{
			coze.BuildUserQuestionText("What can you do?", nil),
		},
	}

	resp, err := cozeCli.Chats.Stream(ctx, req)
	if err != nil {
		log.Fatalf("Error starting chat: %v", err)
	}

	defer resp.Close()
	for {
		event, err := resp.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("Stream finished")
			break
		}
		if err != nil {
			fmt.Println(err)
			break
		}
		if event.Event == coze.ChatEventConversationMessageDelta {
			fmt.Print(event.Message.Content)
		} else if event.Event == coze.ChatEventConversationChatCompleted {
			fmt.Printf("Token usage:%d\n", event.Chat.Usage.TokenCount)
		} else {
			fmt.Printf("\n")
		}
	}

	fmt.Println("done")
}
