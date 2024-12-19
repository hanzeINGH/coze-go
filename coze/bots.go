package coze

import (
	"context"
	"net/http"
	"strconv"

	"github.com/coze/coze/internal"
	"github.com/coze/coze/pagination"
)

// BotMode 机器人模式
type BotMode int

const (
	BotModeSingleAgent BotMode = iota
	BotModeMultiAgent
	BotModeSingleAgentWorkflow
)

// Bot 完整的机器人信息
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
	PluginInfoList []BotPluginInfo    `json:"plugin_info_list,omitempty"`
	ModelInfo      *BotModelInfo      `json:"model_info,omitempty"`
}

// SimpleBot 简化的机器人信息
type SimpleBot struct {
	BotID       string `json:"bot_id"`
	BotName     string `json:"bot_name"`
	Description string `json:"description,omitempty"`
	IconURL     string `json:"icon_url,omitempty"`
	PublishTime string `json:"publish_time,omitempty"`
}

// BotKnowledge 机器人知识库配置
type BotKnowledge struct {
	DatasetIDs     []string `json:"dataset_ids"`
	AutoCall       bool     `json:"auto_call"`
	SearchStrategy int      `json:"search_strategy"`
}

// BotModelInfo 机器人模型信息
type BotModelInfo struct {
	ModelID   string `json:"model_id"`
	ModelName string `json:"model_name"`
}

// BotOnboardingInfo 机器人引导信息
type BotOnboardingInfo struct {
	Prologue           string   `json:"prologue,omitempty"`
	SuggestedQuestions []string `json:"suggested_questions,omitempty"`
}

// BotPluginAPIInfo 机器人插件API信息
type BotPluginAPIInfo struct {
	APIID       string `json:"api_id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// BotPluginInfo 机器人插件信息
type BotPluginInfo struct {
	PluginID    string             `json:"plugin_id"`
	Name        string             `json:"name"`
	Description string             `json:"description,omitempty"`
	IconURL     string             `json:"icon_url,omitempty"`
	APIInfoList []BotPluginAPIInfo `json:"api_info_list,omitempty"`
}

// BotPromptInfo 机器人提示信息
type BotPromptInfo struct {
	Prompt string `json:"prompt"`
}

type CreateBotReq struct {
	SpaceID        string            `json:"space_id"`        // 空间 ID
	Name           string            `json:"name"`            // 名称
	Description    string            `json:"description"`     // 描述
	IconFileID     string            `json:"icon_file_id"`    // 图标文件 ID
	PromptInfo     BotPromptInfo     `json:"prompt_info"`     // 提示信息
	OnboardingInfo BotOnboardingInfo `json:"onboarding_info"` // 上线信息
}

// CreateBotResp 创建机器人响应
type CreateBotResp struct {
	internal.BaseResponse
	Data struct {
		BotID string `json:"bot_id"`
	} `json:"data"`
}

// PublishBotReq 表示发布机器人请求的结构体
type PublishBotReq struct {
	BotID        string   `json:"bot_id"`        // 机器人 ID
	ConnectorIDs []string `json:"connector_ids"` // 连接器 ID 列表
}

// PublishBotResp 发布机器人响应
type PublishBotResp struct {
	internal.BaseResponse
	Data struct {
		BotID      string `json:"bot_id"`
		BotVersion string `json:"version"`
	} `json:"data"`
}

// ListBotReq 表示列出机器人请求的结构体
type ListBotReq struct {
	SpaceID  string `json:"space_id"`  // 空间 ID
	PageNum  int    `json:"page_num"`  // 页码
	PageSize int    `json:"page_size"` // 每页大小
}

// ListBotResp 列出机器人响应
type ListBotResp struct {
	internal.BaseResponse
	Data struct {
		Bots  []SimpleBot `json:"space_bots"`
		Total int         `json:"total"`
	} `json:"data"`
}

// RetrieveBotReq 表示检索机器人请求的结构体
type RetrieveBotReq struct {
	BotID string `json:"bot_id"` // 机器人 ID
}

// RetrieveBotResp 获取机器人响应
type RetrieveBotResp struct {
	internal.BaseResponse
	Data struct {
		Bot *Bot `json:"bot"`
	} `json:"data"`
}

// UpdateBotReq 表示更新机器人请求的结构体
type UpdateBotReq struct {
	BotID          string            `json:"bot_id"`          // 机器人 ID
	Name           string            `json:"name"`            // 名称
	Description    string            `json:"description"`     // 描述
	IconFileID     string            `json:"icon_file_id"`    // 图标文件 ID
	PromptInfo     BotPromptInfo     `json:"prompt_info"`     // 提示信息
	OnboardingInfo BotOnboardingInfo `json:"onboarding_info"` // 上线信息
	Knowledge      BotKnowledge      `json:"knowledge"`       // 知识
}

// UpdateBotResp 更新机器人响应
type UpdateBotResp struct {
	internal.BaseResponse
}

type bots struct {
	client *internal.Client
}

func newBots(client *internal.Client) *bots {
	return &bots{client: client}
}

func (r *bots) Create(ctx context.Context, req CreateBotReq) (*CreateBotResp, error) {
	method := http.MethodPost
	uri := "/v1/bot/create"
	resp := &CreateBotResp{}
	err := r.client.Request(ctx, method, uri, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *bots) Update(ctx context.Context, req UpdateBotReq) (*UpdateBotResp, error) {
	method := http.MethodPost
	uri := "/v1/bot/update"
	resp := &UpdateBotResp{}
	err := r.client.Request(ctx, method, uri, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *bots) Publish(ctx context.Context, req PublishBotReq) (*PublishBotResp, error) {
	method := http.MethodPost
	uri := "/v1/bot/publish"
	resp := &PublishBotResp{}
	err := r.client.Request(ctx, method, uri, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *bots) Retrieve(ctx context.Context, req RetrieveBotReq) (*RetrieveBotResp, error) {
	method := http.MethodGet
	uri := "/v1/bot/get_online_info"
	resp := &RetrieveBotResp{}
	err := r.client.Request(ctx, method, uri, nil, resp, internal.WithQuery("bot_id", req.BotID))
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *bots) List(ctx context.Context, req ListBotReq) (*pagination.PageNumBasedPager[SimpleBot], error) {
	if req.PageSize == 0 {
		req.PageSize = 20
	}
	if req.PageNum == 0 {
		req.PageNum = 1
	}
	return pagination.NewPageNumBasedPager[SimpleBot](
		func(request *pagination.PageRequest) (*pagination.PageResponse[SimpleBot], error) {
			uri := "/v1/space/published_bots_list"
			resp := &ListBotResp{}
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
