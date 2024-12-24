package coze

import (
	"context"
	"net/http"
	"strconv"

	"github.com/coze-dev/coze-go/internal"
	"github.com/coze-dev/coze-go/pagination"
)

func (r *bots) Create(ctx context.Context, req *CreateBotsReq) (*CreateBotsResp, error) {
	method := http.MethodPost
	uri := "/v1/bot/create"
	resp := &CreateBotsResp{}
	err := r.client.Request(ctx, method, uri, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *bots) Update(ctx context.Context, req *UpdateBotsReq) (*UpdateBotsResp, error) {
	method := http.MethodPost
	uri := "/v1/bot/update"
	resp := &UpdateBotsResp{}
	err := r.client.Request(ctx, method, uri, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *bots) Publish(ctx context.Context, req *PublishBotsReq) (*PublishBotsResp, error) {
	method := http.MethodPost
	uri := "/v1/bot/publish"
	resp := &PublishBotsResp{}
	err := r.client.Request(ctx, method, uri, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *bots) Retrieve(ctx context.Context, req *RetrieveBotsReq) (*RetrieveBotsResp, error) {
	method := http.MethodGet
	uri := "/v1/bot/get_online_info"
	resp := &RetrieveBotsResp{}
	err := r.client.Request(ctx, method, uri, nil, resp, internal.WithQuery("bot_id", req.BotID))
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *bots) List(ctx context.Context, req *ListBotsReq) (*pagination.NumberPaged[SimpleBot], error) {
	if req.PageSize == 0 {
		req.PageSize = 20
	}
	if req.PageNum == 0 {
		req.PageNum = 1
	}
	return pagination.NewNumberPaged[SimpleBot](
		func(request *pagination.PageRequest) (*pagination.PageResponse[SimpleBot], error) {
			uri := "/v1/space/published_bots_list"
			resp := &ListBotsResp{}
			err := r.client.Request(ctx, http.MethodGet, uri, nil, resp,
				internal.WithQuery("space_id", req.SpaceID),
				internal.WithQuery("page_num", strconv.Itoa(request.PageNum)),
				internal.WithQuery("page_size", strconv.Itoa(request.PageSize)))
			if err != nil {
				return nil, err
			}
			return &pagination.PageResponse[SimpleBot]{
				Total:   resp.Data.Total,
				HasMore: len(resp.Data.Bots) >= request.PageSize,
				Data:    resp.Data.Bots,
				LogID:   resp.LogID,
			}, nil
		}, req.PageSize, req.PageNum)
}

// BotMode represents the bot mode
type BotMode int

const (
	BotModeMultiAgent          BotMode = 1
	BotModeSingleAgentWorkflow BotMode = 0
)

// Bot represents complete bot information
type Bot struct {
	BotID          string             `json:"bot_id"`
	Name           string             `json:"name"`
	Description    string             `json:"description,omitempty"`
	IconURL        string             `json:"icon_url,omitempty"`
	CreateTime     int64              `json:"create_time"`
	UpdateTime     int64              `json:"update_time"`
	Version        string             `json:"version,omitempty"`
	PromptInfo     *BotPromptInfo     `json:"prompt_info,omitempty"`
	OnboardingInfo *BotOnboardingInfo `json:"onboarding_info,omitempty"`
	BotMode        BotMode            `json:"bot_mode"`
	PluginInfoList []*BotPluginInfo   `json:"plugin_info_list,omitempty"`
	ModelInfo      *BotModelInfo      `json:"model_info,omitempty"`
}

// SimpleBot represents simplified bot information
type SimpleBot struct {
	BotID       string `json:"bot_id"`
	BotName     string `json:"bot_name"`
	Description string `json:"description,omitempty"`
	IconURL     string `json:"icon_url,omitempty"`
	PublishTime string `json:"publish_time,omitempty"`
}

// BotKnowledge represents bot knowledge base configuration
type BotKnowledge struct {
	DatasetIDs     []string `json:"dataset_ids"`
	AutoCall       bool     `json:"auto_call"`
	SearchStrategy int      `json:"search_strategy"`
}

// BotModelInfo represents bot model information
type BotModelInfo struct {
	ModelID   string `json:"model_id"`
	ModelName string `json:"model_name"`
}

// BotOnboardingInfo represents bot onboarding information
type BotOnboardingInfo struct {
	Prologue           string   `json:"prologue,omitempty"`
	SuggestedQuestions []string `json:"suggested_questions,omitempty"`
}

// BotPluginAPIInfo represents bot plugin API information
type BotPluginAPIInfo struct {
	APIID       string `json:"api_id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// BotPluginInfo represents bot plugin information
type BotPluginInfo struct {
	PluginID    string              `json:"plugin_id"`
	Name        string              `json:"name"`
	Description string              `json:"description,omitempty"`
	IconURL     string              `json:"icon_url,omitempty"`
	APIInfoList []*BotPluginAPIInfo `json:"api_info_list,omitempty"`
}

// BotPromptInfo represents bot prompt information
type BotPromptInfo struct {
	Prompt string `json:"prompt"`
}

type CreateBotsReq struct {
	SpaceID        string             `json:"space_id"`        // Space ID
	Name           string             `json:"name"`            // Name
	Description    string             `json:"description"`     // Description
	IconFileID     string             `json:"icon_file_id"`    // Icon file ID
	PromptInfo     *BotPromptInfo     `json:"prompt_info"`     // Prompt information
	OnboardingInfo *BotOnboardingInfo `json:"onboarding_info"` // Onboarding information
}

// CreateBotsResp 创建机器人响应
type CreateBotsResp struct {
	internal.BaseResponse
	Data struct {
		BotID string `json:"bot_id"`
	} `json:"data"`
}

// PublishBotsReq represents the request structure for publishing a bot
type PublishBotsReq struct {
	BotID        string   `json:"bot_id"`        // Bot ID
	ConnectorIDs []string `json:"connector_ids"` // Connector ID list
}

// PublishBotsResp 发布机器人响应
type PublishBotsResp struct {
	internal.BaseResponse
	Data struct {
		BotID      string `json:"bot_id"`
		BotVersion string `json:"version"`
	} `json:"data"`
}

// ListBotsReq represents the request structure for listing bots
type ListBotsReq struct {
	SpaceID  string `json:"space_id"`  // Space ID
	PageNum  int    `json:"page_num"`  // Page number
	PageSize int    `json:"page_size"` // Page size
}

// ListBotsResp 列出机器人响应
type ListBotsResp struct {
	internal.BaseResponse
	Data struct {
		Bots  []*SimpleBot `json:"space_bots"`
		Total int          `json:"total"`
	} `json:"data"`
}

// RetrieveBotsReq represents the request structure for retrieving a bot
type RetrieveBotsReq struct {
	BotID string `json:"bot_id"` // Bot ID
}

// RetrieveBotsResp 获取机器人响应
type RetrieveBotsResp struct {
	internal.BaseResponse
	Bot *Bot `json:"data"`
}

// UpdateBotsReq represents the request structure for updating a bot
type UpdateBotsReq struct {
	BotID          string             `json:"bot_id"`          // Bot ID
	Name           string             `json:"name"`            // Name
	Description    string             `json:"description"`     // Description
	IconFileID     string             `json:"icon_file_id"`    // Icon file ID
	PromptInfo     *BotPromptInfo     `json:"prompt_info"`     // Prompt information
	OnboardingInfo *BotOnboardingInfo `json:"onboarding_info"` // Onboarding information
	Knowledge      *BotKnowledge      `json:"knowledge"`       // Knowledge
}

// UpdateBotsResp 更新机器人响应
type UpdateBotsResp struct {
	internal.BaseResponse
}

type bots struct {
	client *internal.Client
}

func newBots(client *internal.Client) *bots {
	return &bots{client: client}
}
