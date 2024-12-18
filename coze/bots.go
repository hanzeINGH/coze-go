package coze

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

// CreateBotResp 创建机器人响应
type CreateBotResp struct {
	BotID string `json:"bot_id"`
}

// ListBotResp 列出机器人响应
type ListBotResp struct {
	Bots  []SimpleBot `json:"space_bots"`
	Total int         `json:"total"`
}

// PublishBotResp 发布机器人响应
type PublishBotResp struct {
	BotID      string `json:"bot_id"`
	BotVersion string `json:"version"`
}

// RetrieveBotResp 获取机器人响应
type RetrieveBotResp struct {
	Bot *Bot `json:"bot"`
}

// UpdateBotResp 更新机器人响应
type UpdateBotResp struct{}

type bots struct {
}
