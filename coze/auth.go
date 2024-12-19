package coze

import (
	"context"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/chyroc/go-ptr"
	"github.com/coze/coze/auth_error"
	"github.com/coze/coze/internal"
	"github.com/coze/coze/utils"
	"github.com/golang-jwt/jwt"
)

// DeviceAuthReq represents the device authorization request
type DeviceAuthReq struct {
	ClientID string `json:"client_id"`
	LogID    string `json:"log_id,omitempty"`
}

// DeviceAuthResp represents the device authorization response
type DeviceAuthResp struct {
	internal.BaseResponse
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	VerificationURL string `json:"verification_url"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
	LogID           string `json:"log_id,omitempty"`
}

// GetAccessTokenReq represents the access token request
type GetAccessTokenReq struct {
	ClientID        string `json:"client_id"`
	Code            string `json:"code,omitempty"`
	GrantType       string `json:"grant_type"`
	RedirectURI     string `json:"redirect_uri,omitempty"`
	RefreshToken    string `json:"refresh_token,omitempty"`
	CodeVerifier    string `json:"code_verifier,omitempty"`
	DeviceCode      string `json:"device_code,omitempty"`
	DurationSeconds int    `json:"duration_seconds,omitempty"`
	Scope           *Scope `json:"scope,omitempty"`
	LogID           string `json:"log_id,omitempty"`
}

// GetPKCEAuthURLResp represents the PKCE authorization URL response
type GetPKCEAuthURLResp struct {
	internal.BaseResponse
	CodeVerifier     string `json:"code_verifier"`
	AuthorizationURL string `json:"authorization_url"`
}

// GrantType represents the OAuth grant type
type GrantType string

const (
	GrantTypeAuthorizationCode GrantType = "authorization_code"
	GrantTypeDeviceCode        GrantType = "urn:ietf:params:oauth:grant-type:device_code"
	GrantTypeJWTCode           GrantType = "urn:ietf:params:oauth:grant-type:jwt-bearer"
	GrantTypeRefreshToken      GrantType = "refresh_token"
)

func (g GrantType) String() string {
	return string(g)
}

// OAuthToken represents the OAuth token response
type OAuthToken struct {
	internal.BaseResponse
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// Scope represents the OAuth scope
type Scope struct {
	AccountPermission   *ScopeAccountPermission   `json:"account_permission"`
	AttributeConstraint *ScopeAttributeConstraint `json:"attribute_constraint,omitempty"`
}

func (s *Scope) ToMap() (map[string]any, error) {
	data, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	var result map[string]any
	err = json.Unmarshal(data, &result)
	return result, err
}

func BuildBotChat(botIDList []string, permissionList []string) *Scope {
	if len(permissionList) == 0 {
		permissionList = []string{"Connector.botChat"}
	}

	var attributeConstraint *ScopeAttributeConstraint
	if len(botIDList) > 0 {
		chatAttribute := &ScopeAttributeConstraintConnectorBotChatAttribute{
			BotIDList: botIDList,
		}
		attributeConstraint = &ScopeAttributeConstraint{
			ConnectorBotChatAttribute: chatAttribute,
		}
	}

	return &Scope{
		AccountPermission:   &ScopeAccountPermission{PermissionList: permissionList},
		AttributeConstraint: attributeConstraint,
	}
}

// ScopeAccountPermission represents the account permissions in the scope
type ScopeAccountPermission struct {
	PermissionList []string `json:"permission_list"`
}

// ScopeAttributeConstraint represents the attribute constraints in the scope
type ScopeAttributeConstraint struct {
	ConnectorBotChatAttribute *ScopeAttributeConstraintConnectorBotChatAttribute `json:"connector_bot_chat_attribute"`
}

// ScopeAttributeConstraintConnectorBotChatAttribute represents the bot chat attributes
type ScopeAttributeConstraintConnectorBotChatAttribute struct {
	BotIDList []string `json:"bot_id_list"`
}

// CodeChallengeMethod 代码挑战方法
type CodeChallengeMethod string

const (
	CodeChallengeMethodPlain CodeChallengeMethod = "plain"
	CodeChallengeMethodS256  CodeChallengeMethod = "S256"
)

func (m CodeChallengeMethod) String() string {
	return string(m)
}

func (m CodeChallengeMethod) Ptr() *CodeChallengeMethod {
	return &m
}

// OAuthClient OAuth客户端基础结构
type OAuthClient struct {
	clientID     string
	clientSecret string
	baseURL      string
	httpClient   *internal.Client
	hostName     string
}

const (
	getTokenPath               = "/api/permission/oauth2/token"
	getDeviceCodePath          = "/api/permission/oauth2/device/code"
	getWorkspaceDeviceCodePath = "/api/permission/oauth2/workspace_id/%s/device/code"
)

type oauthOpt struct {
	baseURL    string
	httpClient *http.Client
}

type OAuthClientOption func(*oauthOpt)

// WithAuthBaseURL 添加基准url
func WithAuthBaseURL(baseURL string) OAuthClientOption {
	return func(opt *oauthOpt) {
		opt.baseURL = baseURL
	}
}

func WithAuthHttpClient(client *http.Client) OAuthClientOption {
	return func(opt *oauthOpt) {
		opt.httpClient = client
	}
}

// newOAuthClient 创建新的OAuth客户端
func newOAuthClient(clientID, clientSecret string, opts ...OAuthClientOption) (*OAuthClient, error) {
	initSettings := &oauthOpt{
		baseURL: COZE_COM_BASE_URL,
	}

	for _, opt := range opts {
		opt(initSettings)
	}

	var hostName string
	if initSettings.baseURL != "" {
		parsedURL, err := url.Parse(initSettings.baseURL)
		if err != nil {
			return nil, fmt.Errorf("invalid base URL %s: %w", initSettings.baseURL, err)
		}
		hostName = parsedURL.Host
	} else {
		return nil, errors.New("base URL is required")
	}
	var httpClient *http.Client
	if initSettings.httpClient != nil {
		httpClient = initSettings.httpClient
	} else {
		httpClient = http.DefaultClient
	}

	return &OAuthClient{
		clientID:     clientID,
		clientSecret: clientSecret,
		baseURL:      initSettings.baseURL,
		hostName:     hostName,
		httpClient:   internal.NewClient(httpClient, initSettings.baseURL),
	}, nil
}

// getOAuthURL 生成OAuth URL
func (c *OAuthClient) getOAuthURL(redirectURI, state string, opts ...urlOption) string {
	params := url.Values{}
	params.Set("response_type", "code")
	if c.clientID != "" {
		params.Set("client_id", c.clientID)
	}
	if redirectURI != "" {
		params.Set("redirect_uri", redirectURI)
	}
	if state != "" {
		params.Set("state", state)
	}

	for _, opt := range opts {
		opt(&params)
	}

	uri := c.baseURL + "/api/permission/oauth2/authorize"
	return uri + "?" + params.Encode()
}

// getWorkspaceOAuthURL 生成带workspace的OAuth URL
func (c *OAuthClient) getWorkspaceOAuthURL(redirectURI, state, workspaceID string, opts ...urlOption) string {
	params := url.Values{}
	params.Set("response_type", "code")
	if c.clientID != "" {
		params.Set("client_id", c.clientID)
	}
	if redirectURI != "" {
		params.Set("redirect_uri", redirectURI)
	}
	if state != "" {
		params.Set("state", state)
	}

	for _, opt := range opts {
		opt(&params)
	}

	uri := fmt.Sprintf("%s/api/permission/oauth2/workspace_id/%s/authorize", c.baseURL, workspaceID)
	return uri + "?" + params.Encode()
}

type getAccessTokenParams struct {
	Type         GrantType
	Code         string
	Secret       string
	RedirectURI  string
	RefreshToken string
	Request      *GetAccessTokenReq
}

func (c *OAuthClient) getAccessToken(ctx context.Context, params getAccessTokenParams) (*OAuthToken, error) {
	// 如果提供了 Request，直接使用它
	result := &OAuthToken{}
	var req *GetAccessTokenReq
	if params.Request != nil {
		req = params.Request
	} else {
		req = &GetAccessTokenReq{
			ClientID:     c.clientID,
			GrantType:    params.Type.String(),
			Code:         params.Code,
			RefreshToken: params.RefreshToken,
			RedirectURI:  params.RedirectURI,
		}
	}

	opt := make([]internal.RequestOption, 0)
	if params.Secret != "" {
		opt = append(opt, internal.WithHeader(AuthorizeHeader, fmt.Sprintf("Bearer %s", params.Secret)))
	}
	if err := c.httpClient.Request(ctx, http.MethodPost, getTokenPath, req, result, opt...); err != nil {
		return nil, err
	}
	return result, nil
}

// refreshAccessToken 是一个便捷方法，内部调用 getAccessToken
func (c *OAuthClient) refreshAccessToken(ctx context.Context, refreshToken string) (*OAuthToken, error) {
	return c.getAccessToken(ctx, getAccessTokenParams{
		Type:         GrantTypeRefreshToken,
		RefreshToken: refreshToken,
	})
}

// refreshAccessToken 是一个便捷方法，内部调用 getAccessToken
func (c *OAuthClient) refreshAccessTokenWithClientSecret(ctx context.Context, refreshToken string) (*OAuthToken, error) {
	return c.getAccessToken(ctx, getAccessTokenParams{
		Secret:       c.clientSecret,
		Type:         GrantTypeRefreshToken,
		RefreshToken: refreshToken,
	})
}

// PKCEOAuthClient PKCE OAuth客户端
type PKCEOAuthClient struct {
	*OAuthClient
}

// NewPKCEOAuthClient 创建新的PKCE OAuth客户端
func NewPKCEOAuthClient(clientID string, opts ...OAuthClientOption) (*PKCEOAuthClient, error) {
	client, err := newOAuthClient(clientID, "", opts...)
	if err != nil {
		return nil, err
	}
	return &PKCEOAuthClient{
		OAuthClient: client,
	}, err
}

// GenOAuthURL 生成OAuth URL
func (c *PKCEOAuthClient) GenOAuthURL(redirectURI, state string, method *CodeChallengeMethod) (*GetPKCEAuthURLResp, error) {
	if method == nil {
		method = CodeChallengeMethodS256.Ptr()
	}
	return c.GenOAuthURLWithMethod(redirectURI, state, *method)
}

// GenOAuthURLWithWorkspace 生成带workspace的OAuth URL
func (c *PKCEOAuthClient) GenOAuthURLWithWorkspace(redirectURI, state, workspaceID string, method *CodeChallengeMethod) (*GetPKCEAuthURLResp, error) {
	if method == nil {
		method = CodeChallengeMethodS256.Ptr()
	}
	return c.GenOAuthURLWithMethodAndWorkspace(redirectURI, state, *method, workspaceID)
}

// GenOAuthURLWithMethod 使用指定方法生成OAuth URL
func (c *PKCEOAuthClient) GenOAuthURLWithMethod(redirectURI, state string, method CodeChallengeMethod) (*GetPKCEAuthURLResp, error) {
	codeVerifier, err := utils.GenerateRandomString(16)
	if err != nil {
		return nil, err
	}
	code, err := c.getCode(codeVerifier, method)
	if err != nil {
		return nil, err
	}

	authorizationURL := c.getOAuthURL(redirectURI, state,
		withCodeChallenge(code),
		withCodeChallengeMethod(string(method)))

	return &GetPKCEAuthURLResp{
		CodeVerifier:     codeVerifier,
		AuthorizationURL: authorizationURL,
	}, nil
}

// GenOAuthURLWithMethodAndWorkspace 使用指定方法生成带workspace的OAuth URL
func (c *PKCEOAuthClient) GenOAuthURLWithMethodAndWorkspace(redirectURI, state string, method CodeChallengeMethod, workspaceID string) (*GetPKCEAuthURLResp, error) {
	codeVerifier, err := utils.GenerateRandomString(16)
	if err != nil {
		return nil, err
	}
	code, err := c.getCode(codeVerifier, method)
	if err != nil {
		return nil, err
	}

	authorizationURL := c.getWorkspaceOAuthURL(redirectURI, state, workspaceID,
		withCodeChallenge(code),
		withCodeChallengeMethod(string(method)))

	return &GetPKCEAuthURLResp{
		CodeVerifier:     codeVerifier,
		AuthorizationURL: authorizationURL,
	}, nil
}

// getCode 获取验证码
func (c *PKCEOAuthClient) getCode(codeVerifier string, method CodeChallengeMethod) (string, error) {
	if method == CodeChallengeMethodPlain {
		return codeVerifier, nil
	}
	return genS256CodeChallenge(codeVerifier)
}

func (c *PKCEOAuthClient) GetAccessToken(ctx context.Context, code, redirectURI, codeVerifier string) (*OAuthToken, error) {
	req := &GetAccessTokenReq{
		ClientID:     c.clientID,
		GrantType:    string(GrantTypeAuthorizationCode),
		Code:         code,
		RedirectURI:  redirectURI,
		CodeVerifier: codeVerifier,
	}
	return c.getAccessToken(ctx, getAccessTokenParams{
		Request: req,
	})
}

// RefreshToken 刷新令牌
func (c *PKCEOAuthClient) RefreshToken(ctx context.Context, refreshToken string) (*OAuthToken, error) {
	return c.refreshAccessToken(ctx, refreshToken)
}

// genS256CodeChallenge 生成S256验证码挑战
func genS256CodeChallenge(codeVerifier string) (string, error) {
	hash := sha256.New()
	hash.Write([]byte(codeVerifier))
	b64 := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(hash.Sum(nil))
	return strings.ReplaceAll(b64, "=", ""), nil
}

// urlOption URL选项函数类型
type urlOption func(*url.Values)

// withCodeChallenge 添加code_challenge参数
func withCodeChallenge(challenge string) urlOption {
	return func(v *url.Values) {
		v.Set("code_challenge", challenge)
	}
}

// withCodeChallengeMethod 添加code_challenge_method参数
func withCodeChallengeMethod(method string) urlOption {
	return func(v *url.Values) {
		v.Set("code_challenge_method", method)
	}
}

// DeviceOAuthClient 设备OAuth客户端
type DeviceOAuthClient struct {
	*OAuthClient
}

// NewDeviceOAuthClient 创建新的设备OAuth客户端
func NewDeviceOAuthClient(clientID string, opts ...OAuthClientOption) (*DeviceOAuthClient, error) {
	client, err := newOAuthClient(clientID, "", opts...)
	if err != nil {
		return nil, err
	}
	return &DeviceOAuthClient{
		OAuthClient: client,
	}, err
}

// GetDeviceCode 获取设备码
func (c *DeviceOAuthClient) GetDeviceCode(ctx context.Context) (*DeviceAuthResp, error) {
	return c.doGetDeviceCode(ctx, nil)
}

// GetDeviceCodeWithWorkspace 获取带workspace的设备码
func (c *DeviceOAuthClient) GetDeviceCodeWithWorkspace(ctx context.Context, workspaceID string) (*DeviceAuthResp, error) {
	return c.doGetDeviceCode(ctx, &workspaceID)
}

func (c *DeviceOAuthClient) doGetDeviceCode(ctx context.Context, workspaceID *string) (*DeviceAuthResp, error) {
	urlPath := ""
	if workspaceID == nil {
		urlPath = getDeviceCodePath
	} else {
		urlPath = fmt.Sprintf(getWorkspaceDeviceCodePath, *workspaceID)
	}
	req := DeviceAuthReq{
		ClientID: c.clientID,
	}
	result := &DeviceAuthResp{}
	err := c.httpClient.Request(ctx, http.MethodPost, urlPath, req, result)
	if err != nil {
		return nil, err
	}
	result.VerificationURL = fmt.Sprintf("%s?user_code=%s", result.VerificationURI, result.UserCode)
	return result, nil
}

func (c *DeviceOAuthClient) GetAccessToken(ctx context.Context, deviceCode string, poll bool) (*OAuthToken, error) {
	req := &GetAccessTokenReq{
		ClientID:   c.clientID,
		GrantType:  string(GrantTypeDeviceCode),
		DeviceCode: deviceCode,
	}

	if !poll {
		return c.doGetAccessToken(ctx, req)
	}

	interval := 5
	for {
		var resp *OAuthToken
		var err error
		// todo log
		if resp, err = c.doGetAccessToken(ctx, req); err == nil {
			return resp, nil
		}
		authErr, ok := auth_error.AsCozeAuthError(err)
		if !ok {
			return nil, err
		}
		switch authErr.Code {
		case auth_error.AuthorizationPending:
			fmt.Printf("pending, sleep:%ds\n", interval)
		case auth_error.SlowDown:
			if interval < 30 {
				interval += 5
			}
			fmt.Printf("slow down, sleep:%ds\n", interval)
		default:
			return nil, err
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

func (c *DeviceOAuthClient) doGetAccessToken(ctx context.Context, req *GetAccessTokenReq) (*OAuthToken, error) {
	result := &OAuthToken{}
	if err := c.httpClient.Request(ctx, http.MethodPost, getTokenPath, req, result); err != nil {
		return nil, err
	}
	return result, nil
}

// RefreshToken 刷新令牌
func (c *DeviceOAuthClient) RefreshToken(ctx context.Context, refreshToken string) (*OAuthToken, error) {
	return c.refreshAccessToken(ctx, refreshToken)
}

// JWTOAuthClient JWT OAuth客户端
type JWTOAuthClient struct {
	*OAuthClient
	ttl        int
	privateKey *rsa.PrivateKey
	publicKey  string
}

// NewJWTOAuthClient 创建新的JWT OAuth客户端
func NewJWTOAuthClient(clientID, publicKey, privateKeyPEM string, ttl *int, opts ...OAuthClientOption) (*JWTOAuthClient, error) {
	privateKey, err := parsePrivateKey(privateKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}
	client, err := newOAuthClient(clientID, "", opts...)
	if err != nil {
		return nil, err
	}
	if ttl == nil {
		ttl = ptr.Ptr(900) // 默认15分钟
	}
	jwtClient := &JWTOAuthClient{
		OAuthClient: client,
		ttl:         *ttl,
		privateKey:  privateKey,
		publicKey:   publicKey,
	}

	return jwtClient, nil
}

// JWTGetAccessTokenOptions JWT OAuth获取token的选项
type JWTGetAccessTokenOptions struct {
	TTL         int     `json:"ttl,omitempty"`          // token有效期（秒）
	Scope       *Scope  `json:"scope,omitempty"`        // 权限范围
	SessionName *string `json:"session_name,omitempty"` // 会话名称
}

// GetAccessToken 获取访问令牌，使用选项模式
func (c *JWTOAuthClient) GetAccessToken(ctx context.Context, opts *JWTGetAccessTokenOptions) (*OAuthToken, error) {
	if opts == nil {
		opts = &JWTGetAccessTokenOptions{}
	}

	ttl := c.ttl
	if opts.TTL > 0 {
		ttl = opts.TTL
	}

	jwtCode, err := c.generateJWT(ttl, opts.SessionName)
	if err != nil {
		return nil, err
	}

	req := getAccessTokenParams{
		Type:   GrantTypeJWTCode,
		Secret: jwtCode,
		Request: &GetAccessTokenReq{
			ClientID:  c.clientID,
			GrantType: string(GrantTypeJWTCode),
			Scope:     opts.Scope,
		},
	}
	return c.getAccessToken(ctx, req)
}

func (c *JWTOAuthClient) generateJWT(ttl int, sessionName *string) (string, error) {
	now := time.Now()
	jti, err := utils.GenerateRandomString(16)
	if err != nil {
		return "", err
	}

	// 构建 claims
	claims := jwt.MapClaims{
		"iss": c.clientID,
		"aud": c.hostName,
		"iat": now.Unix(),
		"exp": now.Add(time.Duration(ttl) * time.Second).Unix(),
		"jti": jti,
	}

	// 如果有 session_name,添加到 claims 中
	if sessionName != nil {
		claims["session_name"] = *sessionName
	}

	// 创建 token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// 设置 header
	token.Header["kid"] = c.publicKey
	token.Header["typ"] = "JWT"
	token.Header["alg"] = "RS256"

	// 签名并获取完整的令牌字符串
	tokenString, err := token.SignedString(c.privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// WebOAuthClient Web OAuth客户端
type WebOAuthClient struct {
	*OAuthClient
}

// NewWebOAuthClient 创建新的Web OAuth客户端
func NewWebOAuthClient(clientID, clientSecret string, opts ...OAuthClientOption) (*WebOAuthClient, error) {
	client, err := newOAuthClient(clientID, clientSecret, opts...)
	if err != nil {
		return nil, err
	}
	return &WebOAuthClient{
		OAuthClient: client,
	}, err
}

// GetAccessToken 获取访问令牌
func (c *WebOAuthClient) GetAccessToken(ctx context.Context, code, redirectURI string) (*OAuthToken, error) {
	req := &GetAccessTokenReq{
		ClientID:    c.clientID,
		GrantType:   string(GrantTypeAuthorizationCode),
		Code:        code,
		RedirectURI: redirectURI,
	}
	return c.getAccessToken(ctx, getAccessTokenParams{
		Secret:  c.clientSecret,
		Request: req,
	})
}

// RefreshToken 刷新令牌
func (c *WebOAuthClient) GetOAuthURL(redirectURI, state string) string {
	return c.getOAuthURL(redirectURI, state)
}

// GetAccessToken 获取访问令牌
func (c *WebOAuthClient) GetOAuthURLWithWorkspace(redirectURI, state, workspaceID string) string {
	return c.getWorkspaceOAuthURL(redirectURI, state, workspaceID)
}

// RefreshToken 刷新令牌
func (c *WebOAuthClient) RefreshToken(ctx context.Context, refreshToken string) (*OAuthToken, error) {
	return c.refreshAccessTokenWithClientSecret(ctx, refreshToken)
}

// 工具函数
func parsePrivateKey(privateKeyPEM string) (*rsa.PrivateKey, error) {
	// 移除PEM头尾和空白字符
	privateKeyPEM = strings.ReplaceAll(privateKeyPEM, "-----BEGIN PRIVATE KEY-----", "")
	privateKeyPEM = strings.ReplaceAll(privateKeyPEM, "-----END PRIVATE KEY-----", "")
	privateKeyPEM = strings.ReplaceAll(privateKeyPEM, "\n", "")
	privateKeyPEM = strings.ReplaceAll(privateKeyPEM, "\r", "")
	privateKeyPEM = strings.ReplaceAll(privateKeyPEM, " ", "")

	// 解码Base64
	block, err := base64.StdEncoding.DecodeString(privateKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key: %w", err)
	}

	// 解析PKCS8私钥
	key, err := x509.ParsePKCS8PrivateKey(block)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key is not RSA")
	}

	return rsaKey, nil
}
