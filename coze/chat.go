package coze

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/chyroc/go-ptr"
	"github.com/coze-dev/coze-go/coze/internal"
)

// ChatStatus The running status of the session.
type ChatStatus string

const (
	// The session has been created.
	ChatStatusCreated ChatStatus = "created"
	// The Bot is processing.
	ChatStatusInProgress ChatStatus = "in_progress"
	// The Bot has finished processing, and the session has ended.
	ChatStatusCompleted ChatStatus = "completed"
	// The session has failed.
	ChatStatusFailed ChatStatus = "failed"
	// The session is interrupted and requires further processing.
	ChatStatusRequiresAction ChatStatus = "requires_action"
	// The session is user cancelled chat.
	ChatStatusCancelled ChatStatus = "canceled"
)

// ChatEventType Event types for chat.
type ChatEventType string

const (
	// Event for creating a conversation, indicating the start of the conversation.
	ChatEventConversationChatCreated ChatEventType = "conversation.chat.created"
	// The server is processing the conversation.
	ChatEventConversationChatInProgress ChatEventType = "conversation.chat.in_progress"
	// Incremental message, usually an incremental message when type=answer.
	ChatEventConversationMessageDelta ChatEventType = "conversation.message.delta"
	// The message has been completely replied to.
	ChatEventConversationMessageCompleted ChatEventType = "conversation.message.completed"
	// The conversation is completed.
	ChatEventConversationChatCompleted ChatEventType = "conversation.chat.completed"
	// This event is used to mark a failed conversation.
	ChatEventConversationChatFailed ChatEventType = "conversation.chat.failed"
	// The conversation is interrupted and requires the user to report the execution results of the tool.
	ChatEventConversationChatRequiresAction ChatEventType = "conversation.chat.requires_action"
	// Audio delta event
	ChatEventConversationAudioDelta ChatEventType = "conversation.audio.delta"
	// Error events during the streaming response process.
	ChatEventError ChatEventType = "error"
	// The streaming response for this session ended normally.
	ChatEventDone ChatEventType = "done"
)

// Chat represents chat information
type Chat struct {
	// The ID of the chat.
	ID string `json:"id"`
	// The ID of the conversation.
	ConversationID string `json:"conversation_id"`
	// The ID of the bot.
	BotID string `json:"bot_id"`
	// Indicates the create time of the chat. The value format is Unix timestamp in seconds.
	CreatedAt int `json:"created_at"`
	// Indicates the end time of the chat. The value format is Unix timestamp in seconds.
	CompletedAt int `json:"completed_at,omitempty"`
	// Indicates the failure time of the chat. The value format is Unix timestamp in seconds.
	FailedAt int `json:"failed_at,omitempty"`
	// Additional information when creating a message, and this additional information will also be
	// returned when retrieving messages.
	MetaData map[string]string `json:"meta_data,omitempty"`
	// When the chat encounters an auth_error, this field returns detailed error information.
	LastError *ChatError `json:"last_error,omitempty"`
	// The running status of the session.
	Status ChatStatus `json:"status"`
	// Details of the information needed for execution.
	RequiredAction *ChatRequiredAction `json:"required_action,omitempty"`
	// Detailed information about Token consumption.
	Usage *ChatUsage `json:"usage,omitempty"`
}

// ChatError represents error information
type ChatError struct {
	// The error code. An integer type. 0 indicates success, other values indicate failure.
	Code int `json:"code"`
	// The error message. A string type.
	Msg string `json:"msg"`
}

// ChatUsage represents token usage information
type ChatUsage struct {
	// The total number of Tokens consumed in this chat, including the consumption for both the input
	// and output parts.
	TokenCount int `json:"token_count"`
	// The total number of Tokens consumed for the output part.
	OutputCount int `json:"output_count"`
	// The total number of Tokens consumed for the input part.
	InputCount int `json:"input_count"`
}

// ChatRequiredAction represents required action information
type ChatRequiredAction struct {
	// The type of additional operation, with the enum value of submit_tool_outputs.
	Type string `json:"type"`
	// Details of the results that need to be submitted, uploaded through the submission API, and the
	// chat can continue afterward.
	SubmitToolOutputs *ChatSubmitToolOutputs `json:"submit_tool_outputs,omitempty"`
}

// ChatSubmitToolOutputs represents tool outputs that need to be submitted
type ChatSubmitToolOutputs struct {
	// Details of the specific reported information.
	ToolCalls []ChatToolCall `json:"tool_calls"`
}

// ChatToolCall represents a tool call
type ChatToolCall struct {
	// The ID for reporting the running results.
	ID string `json:"id"`
	// The type of tool, with the enum value of function.
	Type string `json:"type"`
	// The definition of the execution method function.
	Function ChatToolCallFunction `json:"function"`
}

// ChatToolCallFunction represents a function call in a tool
type ChatToolCallFunction struct {
	// The name of the method.
	Name string `json:"name"`
	// The parameters of the method.
	Arguments string `json:"arguments"`
}

// ToolOutput represents the output of a tool
type ToolOutput struct {
	// The ID for reporting the running results. You can get this ID under the tool_calls field in
	// response of the Chat API.
	ToolCallID string `json:"tool_call_id"`
	// The execution result of the tool.
	Output string `json:"output"`
}

// CreateChatReq represents the request to create a chat
type CreateChatReq struct {
	// Indicate which conversation the chat is taking place in.
	ConversationID string `json:"-"`

	// The ID of the bot that the API interacts with.
	BotID string `json:"bot_id"`

	// The user who calls the API to chat with the bot.
	UserID string `json:"user_id"`

	// Additional information for the conversation. You can pass the user's query for this
	// conversation through this field. The array length is limited to 100, meaning up to 100 messages can be input.
	Messages []Message `json:"additional_messages,omitempty"`

	// developer can ignore this param
	Stream *bool `json:"stream,omitempty"`

	// The customized variable in a key-value pair.
	CustomVariables map[string]string `json:"custom_variables,omitempty"`

	// Whether to automatically save the history of conversation records.
	AutoSaveHistory *bool `json:"auto_save_history,omitempty"`

	// Additional information, typically used to encapsulate some business-related fields.
	MetaData map[string]string `json:"meta_data,omitempty"`
}

// CancelChatReq represents the request to cancel a chat
type CancelChatReq struct {
	// The Conversation ID can be viewed in the 'conversation_id' field of the Response when
	// initiating a conversation through the Chat API.
	ConversationID string `json:"conversation_id"`

	// The Chat ID can be viewed in the 'id' field of the Response when initiating a chat through the
	// Chat API. If it is a streaming response, check the 'id' field in the chat event of the Response.
	ChatID string `json:"chat_id"`
}

// RetrieveChatReq represents the request to retrieve a chat
type RetrieveChatReq struct {
	// The Conversation ID can be viewed in the 'conversation_id' field of the Response when
	// initiating a conversation through the Chat API.
	ConversationID string `json:"conversation_id"`

	// The Chat ID can be viewed in the 'id' field of the Response when initiating a chat through the
	// Chat API. If it is a streaming response, check the 'id' field in the chat event of the Response.
	ChatID string `json:"chat_id"`
}

// SubmitToolOutputsReq represents the request to submit tool outputs
type SubmitToolOutputsReq struct {
	// The Conversation ID can be viewed in the 'conversation_id' field of the Response when
	// initiating a conversation through the Chat API.
	ConversationID string `json:"-"`

	// The Chat ID can be viewed in the 'id' field of the Response when initiating a chat through the
	// Chat API. If it is a streaming response, check the 'id' field in the chat event of the Response.
	ChatID string `json:"-"`

	// The execution result of the tool. For detailed instructions, refer to the ToolOutput Object
	ToolOutputs []ToolOutput `json:"tool_outputs"`

	Stream *bool `json:"stream,omitempty"`
}

// CreateChatResp represents the response to create a chat
type CreateChatResp struct {
	internal.BaseResponse
	Chat *Chat `json:"data"`
}

// CancelChatResp represents the response to cancel a chat
type CancelChatResp struct {
	internal.BaseResponse
	Chat *Chat `json:"data"`
}

// RetrieveChatResp represents the response to retrieve a chat
type RetrieveChatResp struct {
	internal.BaseResponse
	Chat *Chat `json:"data"`
}

// SubmitToolOutputsResp represents the response to submit tool outputs
type SubmitToolOutputsResp struct {
	internal.BaseResponse
	Chat *Chat `json:"data"`
}

// ChatEvent represents a chat event in the streaming response
type ChatEvent struct {
	Event   ChatEventType `json:"event"`
	Chat    *Chat         `json:"chat,omitempty"`
	Message *Message      `json:"message,omitempty"`
	LogID   string        `json:"log_id,omitempty"`
}

func ChatEventParseParse(eventLine map[string]string, logID string) (*ChatEvent, error) {
	eventType := ChatEventType(eventLine["event"])
	data := eventLine["data"]
	switch eventType {
	case ChatEventDone:
		return &ChatEvent{Event: eventType, LogID: logID}, nil
	case ChatEventError:
		// todo
		return nil, errors.New(data)
	case ChatEventConversationMessageDelta, ChatEventConversationMessageCompleted, ChatEventConversationAudioDelta:
		message := &Message{}
		if err := json.Unmarshal([]byte(data), message); err != nil {
			return nil, err
		}
		return &ChatEvent{Event: eventType, Message: message, LogID: logID}, nil
	case ChatEventConversationChatCreated, ChatEventConversationChatInProgress, ChatEventConversationChatCompleted, ChatEventConversationChatFailed, ChatEventConversationChatRequiresAction:
		chat := &Chat{}
		if err := json.Unmarshal([]byte(data), chat); err != nil {
			return nil, err
		}
		return &ChatEvent{Event: eventType, Chat: chat, LogID: logID}, nil
	default:
		return &ChatEvent{Event: eventType, LogID: logID}, nil
	}
}

func (c *ChatEvent) IsDone() bool {
	return c.Event == ChatEventDone || c.Event == ChatEventError
}

// ChatPoll represents polling information for a chat
type ChatPoll struct {
	Chat     *Chat     `json:"chat"`
	Messages []Message `json:"messages"`
}

type chats struct {
	client   *internal.Client
	Messages *chatMessages
}

func newChats(client *internal.Client) *chats {
	return &chats{
		client:   client,
		Messages: newChatMessages(client),
	}
}

func (r *chats) Chat(ctx context.Context, req CreateChatReq) (*CreateChatResp, error) {
	method := http.MethodPost
	uri := "/v3/chat"
	resp := &CreateChatResp{}
	req.Stream = ptr.Ptr(false)
	req.AutoSaveHistory = ptr.Ptr(true)
	err := r.client.Request(ctx, method, uri, req, resp, internal.WithQuery("conversation_id", req.ConversationID))
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *chats) CreateAndPoll(ctx context.Context, req CreateChatReq, timeout *int) (*ChatPoll, error) {
	req.Stream = ptr.Ptr(false)
	req.AutoSaveHistory = ptr.Ptr(true)

	chatResp, err := r.Chat(ctx, req)
	if err != nil {
		return nil, err
	}
	chat := chatResp.Chat
	conversationID := chat.ConversationID
	now := time.Now()
	for {
		time.Sleep(time.Second)
		if timeout != nil && time.Since(now) > time.Duration(*timeout)*time.Second {
			log.Println("Chat timeout: ", *timeout, " seconds, cancel Chat")
			cancelResp, err := r.Cancel(ctx, CancelChatReq{
				ConversationID: conversationID,
				ChatID:         chat.ID,
			})
			if err != nil {
				return nil, err
			}
			chat = cancelResp.Chat
			break
		}
		retrieveChat, err := r.Retrieve(ctx, RetrieveChatReq{
			ConversationID: conversationID,
			ChatID:         chat.ID,
		})
		if err != nil {
			return nil, err
		}
		if retrieveChat.Chat.Status == ChatStatusCompleted {
			chat = retrieveChat.Chat
			log.Println("Chat completed, spend: ", time.Since(now))
			break
		}
	}
	messages, err := r.Messages.List(ctx, ChatListMessageReq{
		ConversationID: conversationID,
		ChatID:         chat.ID,
	})
	if err != nil {
		return nil, err
	}
	return &ChatPoll{
		Chat:     chat,
		Messages: messages.Messages,
	}, nil
}

func (r *chats) Stream(ctx context.Context, req CreateChatReq) (*ChatEventReader, error) {
	method := http.MethodPost
	uri := "/v3/chat"
	req.Stream = ptr.Ptr(true)
	resp, err := r.client.RowRequest(ctx, method, uri, req, internal.WithQuery("conversation_id", req.ConversationID))
	if err != nil {
		return nil, err
	}

	return &ChatEventReader{
		streamReader: &streamReader[ChatEvent]{
			response:  resp,
			reader:    bufio.NewReader(resp.Body),
			logID:     internal.GetLogID(resp.Header),
			processor: parseChatEvent,
		},
	}, nil
}

type ChatEventReader struct {
	*streamReader[ChatEvent]
}

func parseChatEvent(lineBytes []byte, reader *bufio.Reader, logID string) (*ChatEvent, bool, error) {
	line := string(lineBytes)
	if strings.HasPrefix(line, "event:") {
		event := strings.TrimSpace(line[6:])
		data, err := reader.ReadString('\n')
		if err != nil {
			return nil, false, err
		}
		data = strings.TrimSpace(data[5:])

		eventLine := map[string]string{
			"event": event,
			"data":  data,
		}

		eventData, err := ChatEventParseParse(eventLine, logID)
		if err != nil {
			return nil, false, err
		}

		return eventData, eventData.IsDone(), nil
	}
	return nil, false, nil
}

func (r *chats) Cancel(ctx context.Context, req CancelChatReq) (*CancelChatResp, error) {
	method := http.MethodPost
	uri := "/v3/chat/cancel"
	resp := &CancelChatResp{}
	err := r.client.Request(ctx, method, uri, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *chats) Retrieve(ctx context.Context, req RetrieveChatReq) (*RetrieveChatResp, error) {
	method := http.MethodGet
	uri := "/v3/chat/retrieve"
	resp := &RetrieveChatResp{}
	err := r.client.Request(ctx, method, uri, nil, resp,
		internal.WithQuery("conversation_id", req.ConversationID),
		internal.WithQuery("chat_id", req.ChatID),
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *chats) SubmitToolOutputs(ctx context.Context, req SubmitToolOutputsReq) (*SubmitToolOutputsResp, error) {
	method := http.MethodPost
	uri := "/v3/chat/submit_tool_outputs"
	resp := &SubmitToolOutputsResp{}
	req.Stream = ptr.Ptr(false)
	err := r.client.Request(ctx, method, uri, req, resp,
		internal.WithQuery("conversation_id", req.ConversationID),
		internal.WithQuery("chat_id", req.ChatID),
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *chats) StreamSubmitToolOutputs(ctx context.Context, req SubmitToolOutputsReq) (*ChatEventReader, error) {
	method := http.MethodPost
	req.Stream = ptr.Ptr(true)
	uri := "/v3/chat/submit_tool_outputs"
	resp, err := r.client.RowRequest(ctx, method, uri, req,
		internal.WithQuery("conversation_id", req.ConversationID),
		internal.WithQuery("chat_id", req.ChatID),
	)
	if err != nil {
		return nil, err
	}

	return &ChatEventReader{
		streamReader: &streamReader[ChatEvent]{
			response:  resp,
			reader:    bufio.NewReader(resp.Body),
			logID:     internal.GetLogID(resp.Header),
			processor: parseChatEvent,
		},
	}, nil
}
