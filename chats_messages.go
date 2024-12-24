package coze

import (
	"context"
	"net/http"

	"github.com/coze-dev/coze-go/internal"
)

// ListChatsMessagesReq represents the request to list messages
type ListChatsMessagesReq struct {
	// The Conversation ID can be viewed in the 'conversation_id' field of the Response when
	// initiating a conversation through the Chat API.
	ConversationID string `json:"conversation_id"`

	// The Chat ID can be viewed in the 'id' field of the Response when initiating a chats through the
	// Chat API. If it is a streaming response, check the 'id' field in the chats event of the Response.
	ChatID string `json:"chat_id"`
}

// ListChatsMessagesResp represents the response to list messages
type ListChatsMessagesResp struct {
	internal.BaseResponse
	Messages []*Message `json:"data"`
}

type chatMessages struct {
	client *internal.Client
}

func newChatMessages(client *internal.Client) *chatMessages {
	return &chatMessages{client: client}
}

func (r *chatMessages) List(ctx context.Context, req *ListChatsMessagesReq) (*ListChatsMessagesResp, error) {
	method := http.MethodGet
	uri := "/v3/chats/message/list"
	resp := &ListChatsMessagesResp{}
	err := r.client.Request(ctx, method, uri, nil, resp,
		internal.WithQuery("conversation_id", req.ConversationID),
		internal.WithQuery("chat_id", req.ChatID),
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
