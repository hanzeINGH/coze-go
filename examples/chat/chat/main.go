package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/coze-dev/coze-go"
	"github.com/coze-dev/coze-go/internal"
)

//
// This examples describes how to use the chat interface to initiate conversations,
// poll the status of the conversation, and obtain the messages after the conversation is completed.
//

func main() {
	// Get an access_token through personal access token or oauth.
	token := os.Getenv("COZE_API_TOKEN")
	botID := os.Getenv("PUBLISHED_BOT_ID")
	uid := os.Getenv("USER_ID")

	authCli := coze.NewTokenAuth(token)

	// Init the Coze client through the access_token.
	cozeCli := coze.NewCozeAPI(authCli, coze.WithBaseURL(os.Getenv("COZE_API_BASE")))

	ctx := context.Background()

	//
	// Step one, create chat
	// Call the coze.Create.Create() method to create a chat. The create method is a non-streaming
	// chat and will return a Create class. Developers should periodically check the status of the
	// chat and handle them separately according to different states.
	//
	req := &coze.CreateChatsReq{
		BotID:  botID,
		UserID: uid,
		Messages: []*coze.Message{
			coze.BuildUserQuestionText("What can you do?", nil),
		},
	}

	chatResp, err := cozeCli.Chats.Create(ctx, req)
	if err != nil {
		fmt.Println("Error creating chat:", err)
		return
	}
	fmt.Println(chatResp)
	chat := chatResp.Chat
	chatID := chat.ID
	conversationID := chat.ConversationID

	//
	// Step two, poll the result of chat
	// Assume the development allows at most one chat to run for 10 seconds. If it exceeds 10 seconds,
	// the chat will be cancelled.
	// And when the chat status is not completed, poll the status of the chat once every second.
	// After the chat is completed, retrieve all messages in the chat.
	//
	timeout := time.After(10) // time.Second
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for chat.Status == coze.ChatStatusInProgress {
		select {
		case <-timeout:
			// The chat can be cancelled before its completed.
			cancelResp, err := cozeCli.Chats.Cancel(ctx, &coze.CancelChatsReq{
				ConversationID: conversationID,
				ChatID:         chatID,
			})
			if err != nil {
				fmt.Println("Error cancelling chat:", err)
			}
			fmt.Println(cancelResp)
			break
		case <-ticker.C:
			resp, err := cozeCli.Chats.Retrieve(ctx, &coze.RetrieveChatsReq{
				ConversationID: conversationID,
				ChatID:         chatID,
			})
			if err != nil {
				fmt.Println("Error retrieving chat:", err)
				continue
			}
			fmt.Println(resp)
			chat = resp.Chat
			if chat.Status == coze.ChatStatusCompleted {
				break
			}
		}
	}

	// The sdk provide an automatic polling method.
	chat2, err := cozeCli.Chats.CreateAndPoll(ctx, req, nil)
	if err != nil {
		fmt.Println("Error in CreateAndPoll:", err)
		return
	}
	fmt.Println(chat2)

	// the developer can also set the timeout.
	chat3, err := cozeCli.Chats.CreateAndPoll(ctx, req, internal.Ptr(10))
	if err != nil {
		fmt.Println("Error in CreateAndPollWithTimeout:", err)
		return
	}
	fmt.Println(chat3)
}
