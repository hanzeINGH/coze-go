package coze

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/coze-dev/coze-go/internal"
	"github.com/coze-dev/coze-go/internal/log"
)

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

func (r *chats) Create(ctx context.Context, req *CreateChatsReq) (*CreateChatsResp, error) {
	method := http.MethodPost
	uri := "/v3/chats"
	resp := &CreateChatsResp{}
	req.Stream = internal.Ptr(false)
	req.AutoSaveHistory = internal.Ptr(true)
	err := r.client.Request(ctx, method, uri, req, resp, internal.WithQuery("conversation_id", req.ConversationID))
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *chats) CreateAndPoll(ctx context.Context, req *CreateChatsReq, timeout *int) (*ChatPoll, error) {
	req.Stream = internal.Ptr(false)
	req.AutoSaveHistory = internal.Ptr(true)

	chatResp, err := r.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	chat := chatResp.Chat
	conversationID := chat.ConversationID
	now := time.Now()
	for {
		time.Sleep(time.Second)
		if timeout != nil && time.Since(now) > time.Duration(*timeout)*time.Second {
			log.Infof("Create timeout: ", *timeout, " seconds, cancel Create")
			cancelResp, err := r.Cancel(ctx, &CancelChatsReq{
				ConversationID: conversationID,
				ChatID:         chat.ID,
			})
			if err != nil {
				log.Warnf("Cancel chats failed, err:%v", err)
				return nil, err
			}
			chat = cancelResp.Chat
			break
		}
		retrieveChat, err := r.Retrieve(ctx, &RetrieveChatsReq{
			ConversationID: conversationID,
			ChatID:         chat.ID,
		})
		if err != nil {
			return nil, err
		}
		if retrieveChat.Chat.Status == ChatStatusCompleted {
			chat = retrieveChat.Chat
			log.Infof("Create completed, spend: %v", time.Since(now))
			break
		}
	}
	messages, err := r.Messages.List(ctx, &ListChatsMessagesReq{
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

func (r *chats) Stream(ctx context.Context, req *CreateChatsReq) (*ChatEventReader, error) {
	method := http.MethodPost
	uri := "/v3/chats"
	req.Stream = internal.Ptr(true)
	resp, err := r.client.RawRequest(ctx, method, uri, req, internal.WithQuery("conversation_id", req.ConversationID))
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

func parseChatEvent(lineBytes []byte, reader *bufio.Reader) (*ChatEvent, bool, error) {
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

		eventData, err := doParseChatEvent(eventLine)
		if err != nil {
			return nil, false, err
		}

		return eventData, eventData.IsDone(), nil
	}
	return nil, false, nil
}

func (r *chats) Cancel(ctx context.Context, req *CancelChatsReq) (*CancelChatsResp, error) {
	method := http.MethodPost
	uri := "/v3/chats/cancel"
	resp := &CancelChatsResp{}
	err := r.client.Request(ctx, method, uri, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *chats) Retrieve(ctx context.Context, req *RetrieveChatsReq) (*RetrieveChatsResp, error) {
	method := http.MethodGet
	uri := "/v3/chats/retrieve"
	resp := &RetrieveChatsResp{}
	err := r.client.Request(ctx, method, uri, nil, resp,
		internal.WithQuery("conversation_id", req.ConversationID),
		internal.WithQuery("chat_id", req.ChatID),
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *chats) SubmitToolOutputs(ctx context.Context, req *SubmitToolOutputsChatReq) (*SubmitToolOutputsChatResp, error) {
	method := http.MethodPost
	uri := "/v3/chats/submit_tool_outputs"
	resp := &SubmitToolOutputsChatResp{}
	req.Stream = internal.Ptr(false)
	err := r.client.Request(ctx, method, uri, req, resp,
		internal.WithQuery("conversation_id", req.ConversationID),
		internal.WithQuery("chat_id", req.ChatID),
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *chats) StreamSubmitToolOutputs(ctx context.Context, req *SubmitToolOutputsChatReq) (*ChatEventReader, error) {
	method := http.MethodPost
	req.Stream = internal.Ptr(true)
	uri := "/v3/chats/submit_tool_outputs"
	resp, err := r.client.RawRequest(ctx, method, uri, req,
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
	// The session is user cancelled chats.
	ChatStatusCancelled ChatStatus = "canceled"
)

// ChatEventType Event types for chats.
type ChatEventType string

const (
	// Event for creating a conversation, indicating the start of the conversation.
	ChatEventConversationChatCreated ChatEventType = "conversation.chats.created"
	// The server is processing the conversation.
	ChatEventConversationChatInProgress ChatEventType = "conversation.chats.in_progress"
	// Incremental message, usually an incremental message when type=answer.
	ChatEventConversationMessageDelta ChatEventType = "conversation.message.delta"
	// The message has been completely replied to.
	ChatEventConversationMessageCompleted ChatEventType = "conversation.message.completed"
	// The conversation is completed.
	ChatEventConversationChatCompleted ChatEventType = "conversation.chats.completed"
	// This event is used to mark a failed conversation.
	ChatEventConversationChatFailed ChatEventType = "conversation.chats.failed"
	// The conversation is interrupted and requires the user to report the execution results of the tool.
	ChatEventConversationChatRequiresAction ChatEventType = "conversation.chats.requires_action"
	// Audio delta event
	ChatEventConversationAudioDelta ChatEventType = "conversation.audio.delta"
	// Error events during the streaming response process.
	ChatEventError ChatEventType = "error"
	// The streaming response for this session ended normally.
	ChatEventDone ChatEventType = "done"
)

// Chat represents chats information
type Chat struct {
	// The ID of the chats.
	ID string `json:"id"`
	// The ID of the conversation.
	ConversationID string `json:"conversation_id"`
	// The ID of the bot.
	BotID string `json:"bot_id"`
	// Indicates the create time of the chats. The value format is Unix timestamp in seconds.
	CreatedAt int `json:"created_at"`
	// Indicates the end time of the chats. The value format is Unix timestamp in seconds.
	CompletedAt int `json:"completed_at,omitempty"`
	// Indicates the failure time of the chats. The value format is Unix timestamp in seconds.
	FailedAt int `json:"failed_at,omitempty"`
	// Additional information when creating a message, and this additional information will also be
	// returned when retrieving messages.
	MetaData map[string]string `json:"meta_data,omitempty"`
	// When the chats encounters an auth_error, this field returns detailed error information.
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
	// The total number of Tokens consumed in this chats, including the consumption for both the input
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
	// chats can continue afterward.
	SubmitToolOutputs *ChatSubmitToolOutputs `json:"submit_tool_outputs,omitempty"`
}

// ChatSubmitToolOutputs represents tool outputs that need to be submitted
type ChatSubmitToolOutputs struct {
	// Details of the specific reported information.
	ToolCalls []*ChatToolCall `json:"tool_calls"`
}

// ChatToolCall represents a tool call
type ChatToolCall struct {
	// The ID for reporting the running results.
	ID string `json:"id"`
	// The type of tool, with the enum value of function.
	Type string `json:"type"`
	// The definition of the execution method function.
	Function *ChatToolCallFunction `json:"function"`
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

// CreateChatsReq represents the request to create a chats
type CreateChatsReq struct {
	// Indicate which conversation the chats is taking place in.
	ConversationID string `json:"-"`

	// The ID of the bot that the API interacts with.
	BotID string `json:"bot_id"`

	// The user who calls the API to chats with the bot.
	UserID string `json:"user_id"`

	// Additional information for the conversation. You can pass the user's query for this
	// conversation through this field. The array length is limited to 100, meaning up to 100 messages can be input.
	Messages []*Message `json:"additional_messages,omitempty"`

	// developer can ignore this param
	Stream *bool `json:"stream,omitempty"`

	// The customized variable in a key-value pair.
	CustomVariables map[string]string `json:"custom_variables,omitempty"`

	// Whether to automatically save the history of conversation records.
	AutoSaveHistory *bool `json:"auto_save_history,omitempty"`

	// Additional information, typically used to encapsulate some business-related fields.
	MetaData map[string]string `json:"meta_data,omitempty"`
}

// CancelChatsReq represents the request to cancel a chats
type CancelChatsReq struct {
	// The Conversation ID can be viewed in the 'conversation_id' field of the Response when
	// initiating a conversation through the Chat API.
	ConversationID string `json:"conversation_id"`

	// The Chat ID can be viewed in the 'id' field of the Response when initiating a chats through the
	// Chat API. If it is a streaming response, check the 'id' field in the chats event of the Response.
	ChatID string `json:"chat_id"`
}

// RetrieveChatsReq represents the request to retrieve a chats
type RetrieveChatsReq struct {
	// The Conversation ID can be viewed in the 'conversation_id' field of the Response when
	// initiating a conversation through the Chat API.
	ConversationID string `json:"conversation_id"`

	// The Chat ID can be viewed in the 'id' field of the Response when initiating a chats through the
	// Chat API. If it is a streaming response, check the 'id' field in the chats event of the Response.
	ChatID string `json:"chat_id"`
}

// SubmitToolOutputsChatReq represents the request to submit tool outputs
type SubmitToolOutputsChatReq struct {
	// The Conversation ID can be viewed in the 'conversation_id' field of the Response when
	// initiating a conversation through the Chat API.
	ConversationID string `json:"-"`

	// The Chat ID can be viewed in the 'id' field of the Response when initiating a chats through the
	// Chat API. If it is a streaming response, check the 'id' field in the chats event of the Response.
	ChatID string `json:"-"`

	// The execution result of the tool. For detailed instructions, refer to the ToolOutput Object
	ToolOutputs []*ToolOutput `json:"tool_outputs"`

	Stream *bool `json:"stream,omitempty"`
}

// CreateChatsResp represents the response to create a chats
type CreateChatsResp struct {
	internal.BaseResponse
	Chat *Chat `json:"data"`
}

// CancelChatsResp represents the response to cancel a chats
type CancelChatsResp struct {
	internal.BaseResponse
	Chat *Chat `json:"data"`
}

// RetrieveChatsResp represents the response to retrieve a chats
type RetrieveChatsResp struct {
	internal.BaseResponse
	Chat *Chat `json:"data"`
}

// SubmitToolOutputsChatResp represents the response to submit tool outputs
type SubmitToolOutputsChatResp struct {
	internal.BaseResponse
	Chat *Chat `json:"data"`
}

// ChatEvent represents a chats event in the streaming response
type ChatEvent struct {
	Event   ChatEventType `json:"event"`
	Chat    *Chat         `json:"chats,omitempty"`
	Message *Message      `json:"message,omitempty"`
}

func doParseChatEvent(eventLine map[string]string) (*ChatEvent, error) {
	eventType := ChatEventType(eventLine["event"])
	data := eventLine["data"]
	switch eventType {
	case ChatEventDone:
		return &ChatEvent{Event: eventType}, nil
	case ChatEventError:
		return nil, errors.New(data)
	case ChatEventConversationMessageDelta, ChatEventConversationMessageCompleted, ChatEventConversationAudioDelta:
		message := &Message{}
		if err := json.Unmarshal([]byte(data), message); err != nil {
			return nil, err
		}
		return &ChatEvent{Event: eventType, Message: message}, nil
	case ChatEventConversationChatCreated, ChatEventConversationChatInProgress, ChatEventConversationChatCompleted, ChatEventConversationChatFailed, ChatEventConversationChatRequiresAction:
		chat := &Chat{}
		if err := json.Unmarshal([]byte(data), chat); err != nil {
			return nil, err
		}
		return &ChatEvent{Event: eventType, Chat: chat}, nil
	default:
		return &ChatEvent{Event: eventType}, nil
	}
}

func (c *ChatEvent) IsDone() bool {
	return c.Event == ChatEventDone || c.Event == ChatEventError
}

// ChatPoll represents polling information for a chats
type ChatPoll struct {
	Chat     *Chat      `json:"chats"`
	Messages []*Message `json:"messages"`
}
