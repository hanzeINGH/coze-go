package main

import (
	"context"
	"fmt"
	"os"

	"github.com/coze/coze"
)

func main() {
	// Get an access_token through personal access token or oauth.
	token := os.Getenv("COZE_API_TOKEN")
	authCli := coze.NewTokenAuth(token)
	botID := os.Getenv("COZE_BOT_ID")

	// Init the Coze client through the access_token.
	cozeCli := coze.NewCozeAPI(authCli, coze.WithBaseURL(os.Getenv("COZE_API_BASE")))

	ctx := context.Background()
	// Create a new conversation
	resp, err := cozeCli.Conversations.Create(ctx, coze.CreateConversationReq{BotID: botID})
	if err != nil {
		fmt.Println("Error creating conversation:", err)
		return
	}
	fmt.Println("create conversations:", resp.Conversation)

	conversationID := resp.Conversation.ID

	// Retrieve the conversation
	getResp, err := cozeCli.Conversations.Retrieve(ctx, coze.RetrieveConversationReq{ConversationID: conversationID})
	if err != nil {
		fmt.Println("Error retrieving conversation:", err)
		return
	}
	fmt.Println("retrieve conversations:", getResp)

	// you can manually create message for conversation
	createMessageReq := coze.CreateMessageReq{}
	createMessageReq.Role = coze.MessageRoleAssistant
	createMessageReq.ConversationID = conversationID
	createMessageReq.SetObjectContext([]*coze.MessageObjectString{
		coze.NewFileMessageObjectByURL(os.Getenv("FILE_URL")),
		coze.NewTextMessageObject("hello"),
		coze.NewImageMessageObjectByURL(os.Getenv("IMAGE_FILE_PATH")),
	})

	msgs, err := cozeCli.Conversations.Messages.Create(ctx, createMessageReq)
	if err != nil {
		fmt.Println("Error creating message:", err)
		return
	}
	fmt.Println(msgs)

	// Clear the conversation
	clearResp, err := cozeCli.Conversations.Clear(ctx, coze.ClearConversationReq{ConversationID: conversationID})
	if err != nil {
		fmt.Println("Error clearing conversation:", err)
		return
	}
	fmt.Println(clearResp)
}
