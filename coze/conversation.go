package coze

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/coze-dev/coze-go/coze/internal"
	"github.com/coze-dev/coze-go/coze/pagination"
)

// Conversation represents conversation information
type Conversation struct {
	// The ID of the conversation
	ID string `json:"id"`

	// Indicates the create time of the conversation. The value format is Unix timestamp in seconds.
	CreatedAt int `json:"created_at"`

	// Additional information when creating a message, and this additional information will also be
	// returned when retrieving messages.
	MetaData map[string]string `json:"meta_data,omitempty"`

	// section_id is used to distinguish the context sections of the session history.
	// The same section is one context.
	LastSectionID string `json:"last_section_id"`
}

// CreateConversationReq represents request for creating conversation
type CreateConversationReq struct {
	// Messages in the conversation. For more information, see EnterMessage object.
	Messages []Message `json:"messages,omitempty"`

	// Additional information when creating a message, and this additional information will also be
	// returned when retrieving messages.
	MetaData map[string]string `json:"meta_data,omitempty"`

	// Bind and isolate conversation on different bots.
	BotID string `json:"bot_id,omitempty"`
}

// ListConversationReq represents request for listing conversations
type ListConversationReq struct {
	// The ID of the bot.
	BotID string `json:"bot_id"`

	// The page number.
	PageNum int `json:"page_num,omitempty"`

	// The page size.
	PageSize int `json:"page_size,omitempty"`
}

// RetrieveConversationReq represents request for retrieving conversation
type RetrieveConversationReq struct {
	// The ID of the conversation.
	ConversationID string `json:"conversation_id"`
}

// ClearConversationReq represents request for clearing conversation
type ClearConversationReq struct {
	// The ID of the conversation.
	ConversationID string `json:"conversation_id"`
}

// CreateConversationResp represents response for creating conversation
type CreateConversationResp struct {
	internal.BaseResponse
	Conversation *Conversation `json:"data"`
}

// ListConversationResp represents response for listing conversations
type ListConversationResp struct {
	internal.BaseResponse
	Data struct {
		HasMore       bool           `json:"has_more"`
		Conversations []Conversation `json:"conversations"`
	} `json:"data"`
}

// RetrieveConversationResp represents response for retrieving conversation
type RetrieveConversationResp struct {
	internal.BaseResponse
	Conversation *Conversation `json:"data"`
}

// ClearConversationResp represents response for clearing conversation
type ClearConversationResp struct {
	internal.BaseResponse
	Data struct {
		ConversationID string `json:"conversation_id"`
	} `json:"data"`
}

type conversations struct {
	client   *internal.Client
	Messages *conversationMessage
}

func newConversations(client *internal.Client) *conversations {
	return &conversations{
		client:   client,
		Messages: newConversationMessage(client),
	}
}

func (r *conversations) List(ctx context.Context, req ListConversationReq) (*pagination.PageNumBasedPager[Conversation], error) {
	if req.PageSize == 0 {
		req.PageSize = 20
	}
	if req.PageNum == 0 {
		req.PageNum = 1
	}
	return pagination.NewPageNumBasedPager[Conversation](
		func(request *pagination.PageRequest) (*pagination.PageResponse[Conversation], error) {
			uri := "/v1/conversations"
			resp := &ListConversationResp{}
			err := r.client.Request(ctx, http.MethodGet, uri, nil, resp,
				internal.WithQuery("bot_id", req.BotID),
				internal.WithQuery("page_num", strconv.Itoa(request.PageNum)),
				internal.WithQuery("page_size", strconv.Itoa(request.PageSize)))
			if err != nil {
				return nil, err
			}
			return &pagination.PageResponse[Conversation]{
				HasMore: resp.Data.HasMore,
				Data:    resp.Data.Conversations,
				LogID:   resp.LogID,
			}, nil
		}, req.PageSize, req.PageNum)
}

func (r *conversations) Create(ctx context.Context, req CreateConversationReq) (*CreateConversationResp, error) {
	uri := "/v1/conversation/create"
	resp := &CreateConversationResp{}
	err := r.client.Request(ctx, http.MethodPost, uri, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *conversations) Retrieve(ctx context.Context, req RetrieveConversationReq) (*RetrieveConversationResp, error) {
	uri := "/v1/conversation/retrieve"
	resp := &RetrieveConversationResp{}
	err := r.client.Request(ctx, http.MethodGet, uri, nil, resp, internal.WithQuery("conversation_id", req.ConversationID))
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *conversations) Clear(ctx context.Context, req ClearConversationReq) (*ClearConversationResp, error) {
	uri := fmt.Sprintf("/v1/conversations/%s/clear", req.ConversationID)
	resp := &ClearConversationResp{}
	err := r.client.Request(ctx, http.MethodPost, uri, nil, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
