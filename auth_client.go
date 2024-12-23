package coze

import (
	"context"
	"time"

	"github.com/coze-dev/coze-go/internal"
)

type Auth interface {
	Token(ctx context.Context) (string, error)
}

var _ Auth = &tokenAuthImpl{}
var _ Auth = &jwtOAuthImpl{}

// tokenAuthImpl implements the Auth interface with fixed access token.
type tokenAuthImpl struct {
	accessToken string
}

// NewTokenAuth creates a new token authentication instance.
func NewTokenAuth(accessToken string) Auth {
	return &tokenAuthImpl{
		accessToken: accessToken,
	}
}

func NewJWTAuth(client *JWTOAuthClient, opt *JWTGetAccessTokenOptions) Auth {
	ttl := 900
	if opt == nil {
		return &jwtOAuthImpl{
			TTL:    ttl,
			client: client,
		}
	}
	if opt.TTL > 0 {
		ttl = opt.TTL
	}
	return &jwtOAuthImpl{
		TTL:         ttl,
		Scope:       opt.Scope,
		SessionName: opt.SessionName,
		client:      client,
	}
}

// Token returns the access token.
func (a *tokenAuthImpl) Token(ctx context.Context) (string, error) {
	return a.accessToken, nil
}

type jwtOAuthImpl struct {
	TTL         int
	SessionName *string
	Scope       *Scope
	client      *JWTOAuthClient
	accessToken *string
	expireIn    int64
}

func (j *jwtOAuthImpl) needRefresh() bool {
	return j.accessToken == nil || time.Now().Unix() > j.expireIn
}

func (j *jwtOAuthImpl) Token(ctx context.Context) (string, error) {
	if !j.needRefresh() {
		return internal.Value(j.accessToken), nil
	}
	resp, err := j.client.GetAccessToken(ctx, &JWTGetAccessTokenOptions{
		TTL:         j.TTL,
		SessionName: j.SessionName,
		Scope:       j.Scope,
	})
	if err != nil {
		return "", err
	}
	j.accessToken = internal.Ptr(resp.AccessToken)
	j.expireIn = resp.ExpiresIn
	return resp.AccessToken, nil
}
