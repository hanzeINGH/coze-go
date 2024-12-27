package coze

import (
	"context"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

// DeviceAuthReq represents the device authorization request
type DeviceAuthReq struct {
	ClientID string `json:"client_id"`
	LogID    string `json:"log_id,omitempty"`
}

// GetDeviceAuthResp represents the device authorization response
type GetDeviceAuthResp struct {
	baseResponse
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	VerificationURL string `json:"verification_url"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

// getAccessTokenReq represents the access token request
type getAccessTokenReq struct {
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

// GetPKCEOAuthURLResp represents the PKCE authorization URL response
type GetPKCEOAuthURLResp struct {
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

func (r GrantType) String() string {
	return string(r)
}

type getOAuthTokenResp struct {
	baseResponse
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// OAuthToken represents the OAuth token response
type OAuthToken struct {
	baseResponse
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// Scope represents the OAuth scope
type Scope struct {
	AccountPermission   *ScopeAccountPermission   `json:"account_permission"`
	AttributeConstraint *ScopeAttributeConstraint `json:"attribute_constraint,omitempty"`
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

// CodeChallengeMethod represents the code challenge method
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

// OAuthClient represents the base OAuth core structure
type OAuthClient struct {
	core *core

	clientID     string
	clientSecret string
	baseURL      string
	hostName     string
}

const (
	getTokenPath               = "/api/permission/oauth2/token"
	getDeviceCodePath          = "/api/permission/oauth2/device/code"
	getWorkspaceDeviceCodePath = "/api/permission/oauth2/workspace_id/%s/device/code"
)

type oauthOption struct {
	baseURL    string
	httpClient HTTPClient
}

type OAuthClientOption func(*oauthOption)

// WithAuthBaseURL adds base URL
func WithAuthBaseURL(baseURL string) OAuthClientOption {
	return func(opt *oauthOption) {
		opt.baseURL = baseURL
	}
}

func WithAuthHttpClient(client HTTPClient) OAuthClientOption {
	return func(opt *oauthOption) {
		opt.httpClient = client
	}
}

// newOAuthClient creates a new OAuth core
func newOAuthClient(clientID, clientSecret string, opts ...OAuthClientOption) (*OAuthClient, error) {
	initSettings := &oauthOption{
		baseURL: ComBaseURL,
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
	var httpClient HTTPClient
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
		core:         newCore(httpClient, initSettings.baseURL),
	}, nil
}

// getOAuthURL generates OAuth URL
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

// getWorkspaceOAuthURL generates OAuth URL with workspace
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
	Request      *getAccessTokenReq
}

func (c *OAuthClient) getAccessToken(ctx context.Context, params getAccessTokenParams) (*OAuthToken, error) {
	// If Request is provided, use it directly
	result := &OAuthToken{}
	var req *getAccessTokenReq
	if params.Request != nil {
		req = params.Request
	} else {
		req = &getAccessTokenReq{
			ClientID:     c.clientID,
			GrantType:    params.Type.String(),
			Code:         params.Code,
			RefreshToken: params.RefreshToken,
			RedirectURI:  params.RedirectURI,
		}
	}

	opt := make([]RequestOption, 0)
	if params.Secret != "" {
		opt = append(opt, withHTTPHeader(authorizeHeader, fmt.Sprintf("Bearer %s", params.Secret)))
	}
	if err := c.core.Request(ctx, http.MethodPost, getTokenPath, req, result, opt...); err != nil {
		return nil, err
	}
	return result, nil
}

// refreshAccessToken is a convenience method that internally calls getAccessToken
func (c *OAuthClient) refreshAccessToken(ctx context.Context, refreshToken string) (*OAuthToken, error) {
	return c.getAccessToken(ctx, getAccessTokenParams{
		Type:         GrantTypeRefreshToken,
		RefreshToken: refreshToken,
	})
}

// refreshAccessToken is a convenience method that internally calls getAccessToken
func (c *OAuthClient) refreshAccessTokenWithClientSecret(ctx context.Context, refreshToken string) (*OAuthToken, error) {
	return c.getAccessToken(ctx, getAccessTokenParams{
		Secret:       c.clientSecret,
		Type:         GrantTypeRefreshToken,
		RefreshToken: refreshToken,
	})
}

// PKCEOAuthClient PKCE OAuth core
type PKCEOAuthClient struct {
	*OAuthClient
}

// NewPKCEOAuthClient creates a new PKCE OAuth core
func NewPKCEOAuthClient(clientID string, opts ...OAuthClientOption) (*PKCEOAuthClient, error) {
	client, err := newOAuthClient(clientID, "", opts...)
	if err != nil {
		return nil, err
	}
	return &PKCEOAuthClient{
		OAuthClient: client,
	}, err
}

type GetPKCEOAuthURLReq struct {
	RedirectURI string
	State       string
	Method      *CodeChallengeMethod
	WorkspaceID *string
}

// GetOAuthURL generates OAuth URL
func (c *PKCEOAuthClient) GetOAuthURL(ctx context.Context, req *GetPKCEOAuthURLReq) (*GetPKCEOAuthURLResp, error) {
	if req == nil {
		return nil, errors.New("request is required")
	}
	if len(req.RedirectURI) == 0 {
		return nil, errors.New("redirectURI is required")
	}
	method := CodeChallengeMethodS256
	if req.Method != nil {
		method = *req.Method
	}
	codeVerifier, err := generateRandomString(16)
	if err != nil {
		return nil, err
	}
	code, err := c.getCode(codeVerifier, ptrValue(req.Method))
	if err != nil {
		return nil, err
	}
	var authorizationURL string
	if req.WorkspaceID != nil {
		authorizationURL = c.getWorkspaceOAuthURL(req.RedirectURI, req.State, *req.WorkspaceID,
			withCodeChallenge(code),
			withCodeChallengeMethod(string(method)))
	} else {
		authorizationURL = c.getOAuthURL(req.RedirectURI, req.State,
			withCodeChallenge(code),
			withCodeChallengeMethod(string(method)))
	}

	return &GetPKCEOAuthURLResp{
		CodeVerifier:     codeVerifier,
		AuthorizationURL: authorizationURL,
	}, nil
}

// getCode gets the verification code
func (c *PKCEOAuthClient) getCode(codeVerifier string, method CodeChallengeMethod) (string, error) {
	if method == CodeChallengeMethodPlain {
		return codeVerifier, nil
	}
	return genS256CodeChallenge(codeVerifier)
}

type GetPKCEAccessTokenReq struct {
	Code, RedirectURI, CodeVerifier string
}

func (c *PKCEOAuthClient) GetAccessToken(ctx context.Context, req *GetPKCEAccessTokenReq) (*OAuthToken, error) {
	return c.getAccessToken(ctx, getAccessTokenParams{
		Request: &getAccessTokenReq{
			ClientID:     c.clientID,
			GrantType:    string(GrantTypeAuthorizationCode),
			Code:         req.Code,
			RedirectURI:  req.RedirectURI,
			CodeVerifier: req.CodeVerifier,
		},
	})
}

// RefreshToken refreshes the access token
func (c *PKCEOAuthClient) RefreshToken(ctx context.Context, refreshToken string) (*OAuthToken, error) {
	return c.refreshAccessToken(ctx, refreshToken)
}

// genS256CodeChallenge generates S256 code challenge
func genS256CodeChallenge(codeVerifier string) (string, error) {
	hash := sha256.New()
	hash.Write([]byte(codeVerifier))
	b64 := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(hash.Sum(nil))
	return strings.ReplaceAll(b64, "=", ""), nil
}

// urlOption represents URL option function type
type urlOption func(*url.Values)

// withCodeChallenge adds code_challenge parameter
func withCodeChallenge(challenge string) urlOption {
	return func(v *url.Values) {
		v.Set("code_challenge", challenge)
	}
}

// withCodeChallengeMethod adds code_challenge_method parameter
func withCodeChallengeMethod(method string) urlOption {
	return func(v *url.Values) {
		v.Set("code_challenge_method", method)
	}
}

// DeviceOAuthClient represents the device OAuth core
type DeviceOAuthClient struct {
	*OAuthClient
}

// NewDeviceOAuthClient creates a new device OAuth core
func NewDeviceOAuthClient(clientID string, opts ...OAuthClientOption) (*DeviceOAuthClient, error) {
	client, err := newOAuthClient(clientID, "", opts...)
	if err != nil {
		return nil, err
	}
	return &DeviceOAuthClient{
		OAuthClient: client,
	}, err
}

type GetDeviceOAuthCodeReq struct {
	WorkspaceID *string
}

// GetDeviceCode gets the device code
func (c *DeviceOAuthClient) GetDeviceCode(ctx context.Context, req *GetDeviceOAuthCodeReq) (*GetDeviceAuthResp, error) {
	return c.doGetDeviceCode(ctx, req.WorkspaceID)
}

func (c *DeviceOAuthClient) doGetDeviceCode(ctx context.Context, workspaceID *string) (*GetDeviceAuthResp, error) {
	urlPath := ""
	if workspaceID == nil {
		urlPath = getDeviceCodePath
	} else {
		urlPath = fmt.Sprintf(getWorkspaceDeviceCodePath, *workspaceID)
	}
	req := DeviceAuthReq{
		ClientID: c.clientID,
	}
	result := &GetDeviceAuthResp{}
	err := c.core.Request(ctx, http.MethodPost, urlPath, req, result)
	if err != nil {
		return nil, err
	}
	result.VerificationURL = fmt.Sprintf("%s?user_code=%s", result.VerificationURI, result.UserCode)
	return result, nil
}

type GetDeviceOAuthAccessTokenReq struct {
	DeviceCode string
	Poll       bool
}

func (c *DeviceOAuthClient) GetAccessToken(ctx context.Context, dReq *GetDeviceOAuthAccessTokenReq) (*OAuthToken, error) {
	req := &getAccessTokenReq{
		ClientID:   c.clientID,
		GrantType:  string(GrantTypeDeviceCode),
		DeviceCode: dReq.DeviceCode,
	}

	if !dReq.Poll {
		return c.doGetAccessToken(ctx, req)
	}

	logger.Infof(ctx, "polling get access token\n")
	interval := 5
	for {
		var resp *OAuthToken
		var err error
		if resp, err = c.doGetAccessToken(ctx, req); err == nil {
			return resp, nil
		}
		authErr, ok := AsCozeAuthError(err)
		if !ok {
			return nil, err
		}
		switch authErr.Code {
		case AuthorizationPending:
			logger.Infof(ctx, "pending, sleep:%ds\n", interval)
		case SlowDown:
			if interval < 30 {
				interval += 5
			}
			logger.Infof(ctx, "slow down, sleep:%ds\n", interval)
		default:
			logger.Warnf(ctx, "get access token error:%s, return\n", err.Error())
			return nil, err
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

func (c *DeviceOAuthClient) doGetAccessToken(ctx context.Context, req *getAccessTokenReq) (*OAuthToken, error) {
	resp := &getOAuthTokenResp{}
	if err := c.core.Request(ctx, http.MethodPost, getTokenPath, req, resp); err != nil {
		return nil, err
	}
	res := &OAuthToken{
		AccessToken:  resp.AccessToken,
		ExpiresIn:    resp.ExpiresIn,
		RefreshToken: resp.RefreshToken,
	}
	return res, nil
}

// RefreshToken refreshes the access token
func (c *DeviceOAuthClient) RefreshToken(ctx context.Context, refreshToken string) (*OAuthToken, error) {
	return c.refreshAccessToken(ctx, refreshToken)
}

// JWTOAuthClient represents the JWT OAuth core
type JWTOAuthClient struct {
	*OAuthClient
	ttl        int
	privateKey *rsa.PrivateKey
	publicKey  string
}

type NewJWTOAuthClientParam struct {
	ClientID      string
	PublicKey     string
	PrivateKeyPEM string
	TTL           *int
}

// NewJWTOAuthClient creates a new JWT OAuth core
func NewJWTOAuthClient(param NewJWTOAuthClientParam, opts ...OAuthClientOption) (*JWTOAuthClient, error) {
	privateKey, err := parsePrivateKey(param.PrivateKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}
	client, err := newOAuthClient(param.ClientID, "", opts...)
	if err != nil {
		return nil, err
	}
	ttl := param.TTL
	if ttl == nil {
		ttl = ptr(900) // Default 15 minutes
	}
	jwtClient := &JWTOAuthClient{
		OAuthClient: client,
		ttl:         *ttl,
		privateKey:  privateKey,
		publicKey:   param.PublicKey,
	}

	return jwtClient, nil
}

// GetJWTAccessTokenReq represents options for getting JWT OAuth token
type GetJWTAccessTokenReq struct {
	TTL         int     `json:"ttl,omitempty"`          // Token validity period (in seconds)
	Scope       *Scope  `json:"scope,omitempty"`        // Permission scope
	SessionName *string `json:"session_name,omitempty"` // Session name
}

// GetAccessToken gets the access token, using options pattern
func (c *JWTOAuthClient) GetAccessToken(ctx context.Context, opts *GetJWTAccessTokenReq) (*OAuthToken, error) {
	if opts == nil {
		opts = &GetJWTAccessTokenReq{}
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
		Request: &getAccessTokenReq{
			ClientID:  c.clientID,
			GrantType: string(GrantTypeJWTCode),
			Scope:     opts.Scope,
		},
	}
	return c.getAccessToken(ctx, req)
}

func (c *JWTOAuthClient) generateJWT(ttl int, sessionName *string) (string, error) {
	now := time.Now()
	jti, err := generateRandomString(16)
	if err != nil {
		return "", err
	}

	// Build claims
	claims := jwt.MapClaims{
		"iss": c.clientID,
		"aud": c.hostName,
		"iat": now.Unix(),
		"exp": now.Add(time.Duration(ttl) * time.Second).Unix(),
		"jti": jti,
	}

	// If session_name is provided, add it to claims
	if sessionName != nil {
		claims["session_name"] = *sessionName
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Set header
	token.Header["kid"] = c.publicKey
	token.Header["typ"] = "JWT"
	token.Header["alg"] = "RS256"

	// Sign and get full token string
	tokenString, err := token.SignedString(c.privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// WebOAuthClient Web OAuth core
type WebOAuthClient struct {
	*OAuthClient
}

// NewWebOAuthClient creates a new Web OAuth core
func NewWebOAuthClient(clientID, clientSecret string, opts ...OAuthClientOption) (*WebOAuthClient, error) {
	client, err := newOAuthClient(clientID, clientSecret, opts...)
	if err != nil {
		return nil, err
	}
	return &WebOAuthClient{
		OAuthClient: client,
	}, err
}

type GetWebOAuthAccessTokenReq struct {
	Code, RedirectURI string
}

// GetAccessToken gets the access token
func (c *WebOAuthClient) GetAccessToken(ctx context.Context, req *GetWebOAuthAccessTokenReq) (*OAuthToken, error) {
	return c.getAccessToken(ctx, getAccessTokenParams{
		Secret: c.clientSecret,
		Request: &getAccessTokenReq{
			ClientID:    c.clientID,
			GrantType:   string(GrantTypeAuthorizationCode),
			Code:        req.Code,
			RedirectURI: req.RedirectURI,
		},
	})
}

type GetWebOAuthURLReq struct {
	RedirectURI, State string
	WorkspaceID        *string
}

// GetOAuthURL Get OAuth URL
func (c *WebOAuthClient) GetOAuthURL(ctx context.Context, req *GetWebOAuthURLReq) string {
	if req.WorkspaceID != nil {
		return c.getWorkspaceOAuthURL(req.RedirectURI, req.State, *req.WorkspaceID)
	}
	return c.getOAuthURL(req.RedirectURI, req.State)
}

// RefreshToken refreshes the access token
func (c *WebOAuthClient) RefreshToken(ctx context.Context, refreshToken string) (*OAuthToken, error) {
	return c.refreshAccessTokenWithClientSecret(ctx, refreshToken)
}

// 工具函数
func parsePrivateKey(privateKeyPEM string) (*rsa.PrivateKey, error) {
	// Remove PEM header and footer and whitespace
	privateKeyPEM = strings.ReplaceAll(privateKeyPEM, "-----BEGIN PRIVATE KEY-----", "")
	privateKeyPEM = strings.ReplaceAll(privateKeyPEM, "-----END PRIVATE KEY-----", "")
	privateKeyPEM = strings.ReplaceAll(privateKeyPEM, "\n", "")
	privateKeyPEM = strings.ReplaceAll(privateKeyPEM, "\r", "")
	privateKeyPEM = strings.ReplaceAll(privateKeyPEM, " ", "")

	// Decode Base64
	block, err := base64.StdEncoding.DecodeString(privateKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key: %w", err)
	}

	// Parse PKCS8 private key
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
