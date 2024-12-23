package coze

import (
	"net/http"

	"github.com/coze-dev/coze-go/internal"
	"github.com/coze-dev/coze-go/internal/log"
)

type CozeAPI struct {
	Audio         *audio
	Bots          *bots
	Chats         *chats
	Conversations *conversations
	Workflows     *workflows
	Workspaces    *workspace
	Datasets      *datasets
	Files         *files

	baseURL string
}

type newCozeAPIOpt struct {
	baseURL  string
	client   *http.Client
	logLevel log.LogLevel
}

type CozeAPIOption func(*newCozeAPIOpt)

// WithBaseURL 添加基准url
func WithBaseURL(baseURL string) CozeAPIOption {
	return func(opt *newCozeAPIOpt) {
		opt.baseURL = baseURL
	}
}

// WithHttpClient 设置自定义的 HTTP 客户端
func WithHttpClient(client *http.Client) CozeAPIOption {
	return func(opt *newCozeAPIOpt) {
		opt.client = client
	}
}

// WithLogLevel 设置日志级别
func WithLogLevel(level log.LogLevel) CozeAPIOption {
	return func(opt *newCozeAPIOpt) {
		opt.logLevel = level
	}
}

func WithLogger(logger log.Logger) CozeAPIOption {
	return func(opt *newCozeAPIOpt) {
		log.SetLogger(logger)
	}
}

func NewCozeAPI(auth Auth, opts ...CozeAPIOption) CozeAPI {
	opt := &newCozeAPIOpt{
		baseURL:  CozeComBaseURL,
		logLevel: log.LogInfo, // 默认日志级别为 Info
	}
	for _, option := range opts {
		option(opt)
	}
	if opt.client == nil {
		opt.client = http.DefaultClient
	}
	saveTransport := opt.client.Transport
	if saveTransport == nil {
		saveTransport = http.DefaultTransport
	}
	opt.client.Transport = &authTransport{
		auth: auth,
		next: saveTransport,
	}
	httpClient := internal.NewClient(opt.client, opt.baseURL)

	// 设置日志级别
	log.SetLevel(opt.logLevel)

	cozeClient := CozeAPI{
		Audio:         newAudio(httpClient),
		Bots:          newBots(httpClient),
		Chats:         newChats(httpClient),
		Conversations: newConversations(httpClient),
		Workflows:     newWorkflows(httpClient),
		Workspaces:    newWorkspace(httpClient),
		Datasets:      newDatasets(httpClient),
		Files:         newFiles(httpClient),

		baseURL: opt.baseURL,
	}
	return cozeClient
}

type authTransport struct {
	auth Auth
	next http.RoundTripper
}

func (h *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	accessToken, err := h.auth.Token(req.Context())
	if err != nil {
		log.Errorf("Failed to get access token: %v", err)
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	return h.next.RoundTrip(req)
}
