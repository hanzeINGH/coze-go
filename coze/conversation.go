package coze

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
	Conversation *Conversation `json:"conversation"`
}

// ListConversationResp represents response for listing conversations
type ListConversationResp struct {
	HasMore       bool           `json:"has_more"`
	Conversations []Conversation `json:"conversations"`
}

// RetrieveConversationResp represents response for retrieving conversation
type RetrieveConversationResp struct {
	Conversation *Conversation `json:"conversation"`
}

// ClearConversationResp represents response for clearing conversation
type ClearConversationResp struct {
	ConversationID string `json:"conversation_id"`
}

type conversations struct {
	Messages *conversationMessage
}
