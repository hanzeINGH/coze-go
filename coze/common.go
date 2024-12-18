package coze

// Message represents a message in conversation
type Message struct {
	// The entity that sent this message.
	Role MessageRole `json:"role"`

	// The type of message.
	Type MessageType `json:"type"`

	// The content of the message. It supports various types of content, including plain text,
	// multimodal (a mix of text, images, and files), message cards, and more.
	Content string `json:"content"`

	// The type of message content.
	ContentType MessageContentType `json:"content_type"`

	// Additional information when creating a message, and this additional information will also be
	// returned when retrieving messages. Custom key-value pairs should be specified in Map object
	// format, with a length of 16 key-value pairs. The length of the key should be between 1 and 64
	// characters, and the length of the value should be between 1 and 512 characters.
	MetaData map[string]string `json:"meta_data,omitempty"`

	ID             string `json:"id"`
	ConversationID string `json:"conversation_id"`

	// section_id is used to distinguish the context sections of the session history. The same section
	// is one context.
	SectionID string `json:"section_id"`
	BotID     string `json:"bot_id"`
	ChatID    string `json:"chat_id"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

// MessageRole represents the role of message sender
type MessageRole string

const (
	MessageRoleUnknown MessageRole = "unknown"
	// Indicates that the content of the message is sent by the user.
	MessageRoleUser MessageRole = "user"
	// Indicates that the content of the message is sent by the bot.
	MessageRoleAssistant MessageRole = "assistant"
)

// MessageType represents the type of message
type MessageType string

const (
	// User input content.
	MessageTypeQuestion MessageType = "question"

	// The message content returned by the Bot to the user, supporting incremental return.
	MessageTypeAnswer MessageType = "answer"

	// Intermediate results of the function (function call) called during the Bot conversation
	// process.
	MessageTypeFunctionCall MessageType = "function_call"

	// Results returned after calling the tool (function call).
	MessageTypeToolOutput MessageType = "tool_output"

	// Results returned after calling the tool (function call).
	MessageTypeToolResponse MessageType = "tool_response"

	// If the user question suggestion switch is turned on in the Bot configuration, the reply content
	// related to the recommended questions will be returned.
	MessageTypeFollowUp MessageType = "follow_up"

	// In the scenario of multiple answers, the server will return a verbose package.
	MessageTypeVerbose MessageType = "verbose"

	MessageTypeUnknown MessageType = ""
)

// MessageContentType represents the type of message content
type MessageContentType string

const (
	MessageContentTypeUnknown MessageContentType = "unknown"
	// Text.
	MessageContentTypeText MessageContentType = "text"

	// Multimodal content, that is, a combination of text and files, or a combination of text and
	// images.
	MessageContentTypeObjectString MessageContentType = "object_string"

	// Message card. This enum value only appears in the interface response and is not supported as an
	// input parameter.
	MessageContentTypeCard MessageContentType = "card"

	// If there is a voice message in the input message, the conversation.audio.delta event will be
	// returned in the streaming response event.
	MessageContentTypeAudio MessageContentType = "audio"
)

// MessageObjectString represents a multimodal message object
type MessageObjectString struct {
	// The content type of the multimodal message.
	Type MessageObjectStringType `json:"type"`

	// Text content. Required when type is text.
	Text string `json:"text,omitempty"`

	// The ID of the file or image content.
	FileID string `json:"file_id,omitempty"`

	// The online address of the file or image content. Must be a valid address that is publicly
	// accessible. file_id or file_url must be specified when type is file or image.
	FileURL string `json:"file_url,omitempty"`
}

// MessageObjectStringType represents the type of multimodal message content
type MessageObjectStringType string

const (
	MessageObjectStringTypeUnknown MessageObjectStringType = "unknown"
	MessageObjectStringTypeText    MessageObjectStringType = "text"
	MessageObjectStringTypeFile    MessageObjectStringType = "file"
	MessageObjectStringTypeImage   MessageObjectStringType = "image"
	MessageObjectStringTypeAudio   MessageObjectStringType = "audio"
)

// Helper functions for creating MessageObjectString
func NewTextMessageObject(text string) *MessageObjectString {
	return &MessageObjectString{
		Type: MessageObjectStringTypeText,
		Text: text,
	}
}

func NewImageMessageObject(fileID, fileURL string) *MessageObjectString {
	if fileID == "" && fileURL == "" {
		panic("file_id or file_url must be specified")
	}
	return &MessageObjectString{
		Type:    MessageObjectStringTypeImage,
		FileID:  fileID,
		FileURL: fileURL,
	}
}

func NewFileMessageObject(fileID, fileURL string) *MessageObjectString {
	if fileID == "" && fileURL == "" {
		panic("file_id or file_url must be specified")
	}
	return &MessageObjectString{
		Type:    MessageObjectStringTypeFile,
		FileID:  fileID,
		FileURL: fileURL,
	}
}

func NewAudioMessageObject(fileID, fileURL string) *MessageObjectString {
	if fileID == "" && fileURL == "" {
		panic("file_id or file_url must be specified")
	}
	return &MessageObjectString{
		Type:    MessageObjectStringTypeAudio,
		FileID:  fileID,
		FileURL: fileURL,
	}
}
