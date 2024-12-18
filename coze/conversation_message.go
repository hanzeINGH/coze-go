package coze

// CreateMessageReq represents request for creating message
type CreateMessageReq struct {
	// The ID of the conversation.
	ConversationID string `json:"conversation_id"`

	// The entity that sent this message.
	Role MessageRole `json:"role"`

	// The content of the message, supporting pure text, multimodal (mixed input of text, images, files),
	// cards, and various types of content.
	Content string `json:"content"`

	// The type of message content.
	ContentType MessageContentType `json:"content_type"`

	// Additional information when creating a message, and this additional information will also be
	// returned when retrieving messages.
	MetaData map[string]string `json:"meta_data,omitempty"`
}

// ConversationListMessageReq represents request for listing messages
type ConversationListMessageReq struct {
	// The ID of the conversation.
	ConversationID string `json:"conversation_id"`

	// The sorting method for the message list.
	Order string `json:"order,omitempty"`

	// The ID of the Chat.
	ChatID string `json:"chat_id,omitempty"`

	// Get messages before the specified position.
	BeforeID string `json:"before_id,omitempty"`

	// Get messages after the specified position.
	AfterID string `json:"after_id,omitempty"`

	// The amount of data returned per query. Default is 50, with a range of 1 to 50.
	Limit int `json:"limit,omitempty"`

	BotID string `json:"bot_id,omitempty"`
}

// RetrieveMessageReq represents request for retrieving message
type RetrieveMessageReq struct {
	ConversationID string `json:"conversation_id"`
	MessageID      string `json:"message_id"`
}

// UpdateMessageReq represents request for updating message
type UpdateMessageReq struct {
	// The ID of the conversation.
	ConversationID string `json:"conversation_id"`

	// The ID of the message.
	MessageID string `json:"message_id"`

	// The content of the message, supporting pure text, multimodal (mixed input of text, images, files),
	// cards, and various types of content.
	Content string `json:"content,omitempty"`

	MetaData map[string]string `json:"meta_data,omitempty"`

	// The type of message content.
	ContentType MessageContentType `json:"content_type,omitempty"`
}

// DeleteMessageReq represents request for deleting message
type DeleteMessageReq struct {
	// The ID of the conversation.
	ConversationID string `json:"conversation_id"`

	// message id
	MessageID string `json:"message_id"`
}

// CreateMessageResp represents response for creating message
type CreateMessageResp struct {
	Message *Message `json:"message"`
}

// ConversationListMessageResp represents response for listing messages
type ConversationListMessageResp struct {
	HasMore  bool      `json:"has_more"`
	FirstID  string    `json:"first_id"`
	LastID   string    `json:"last_id"`
	Messages []Message `json:"messages"`
}

// RetrieveMessageResp represents response for retrieving message
type RetrieveMessageResp struct {
	Message *Message `json:"message"`
}

// UpdateMessageResp represents response for updating message
type UpdateMessageResp struct {
	Message *Message `json:"message"`
}

// DeleteMessageResp represents response for deleting message
type DeleteMessageResp struct {
	Message *Message `json:"message"`
}

type conversationMessage struct {
}
