package coze

import (
	"net/http"

	"github.com/coze-dev/coze-go/coze/internal"
)

type CozeAPI struct {
	Audio         *audio
	Bots          *bots
	Chats         *chats
	Conversations *conversations
	Workflows     *workflows
	Workspaces    *workspace
	Datasets      *dataset
	Files         *files

	baseURL string
}

type newCozeAPIOpt struct {
	baseURL string
	client  *http.Client
}

type CozeAPIOption func(*newCozeAPIOpt)

// WithBaseURL 添加基准url
func WithBaseURL(baseURL string) CozeAPIOption {
	return func(opt *newCozeAPIOpt) {
		opt.baseURL = baseURL
	}
}

func WithHttpClient(client *http.Client) CozeAPIOption {
	return func(opt *newCozeAPIOpt) {
		opt.client = client
	}
}

func NewCozeAPI(auth Auth, opts ...CozeAPIOption) CozeAPI {
	opt := &newCozeAPIOpt{
		baseURL: COZE_COM_BASE_URL,
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

	cozeClient := CozeAPI{
		Audio:         newAudio(httpClient),
		Bots:          newBots(httpClient),
		Chats:         newChats(httpClient),
		Conversations: newConversations(httpClient),
		Workflows:     newWorkflows(httpClient),
		Workspaces:    newWorkspace(httpClient),
		Datasets:      newDataset(httpClient),
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
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	return h.next.RoundTrip(req)
}
