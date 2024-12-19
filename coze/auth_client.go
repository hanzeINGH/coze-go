package coze

import (
	"context"
	"time"

	"github.com/chyroc/go-ptr"
)

type Auth interface {
	Token(ctx context.Context) (string, error)
}

// TokenAuth 表示基于令牌的认证
type TokenAuth struct {
	accessToken string
}

// NewTokenAuth 创建一个新的令牌认证实例
func NewTokenAuth(accessToken string) *TokenAuth {
	return &TokenAuth{
		accessToken: accessToken,
	}
}

// Token 获取访问令牌
func (a *TokenAuth) Token(ctx context.Context) (string, error) {
	return a.accessToken, nil
}

type JWTOAuth struct {
	TTL         int
	SessionName *string
	Scope       *Scope
	client      JWTOAuthClient
	accessToken *string
	expireIn    int64
}

func (j *JWTOAuth) needRefresh() bool {
	return j.accessToken == nil || time.Now().Unix() > j.expireIn
}

func (j *JWTOAuth) Token(ctx context.Context) (string, error) {
	if !j.needRefresh() {
		return ptr.Value(j.accessToken), nil
	}
	resp, err := j.client.GetAccessToken(ctx, &JWTGetAccessTokenOptions{
		TTL:         j.TTL,
		SessionName: j.SessionName,
		Scope:       j.Scope,
	})
	if err != nil {
		return "", err
	}
	j.accessToken = ptr.Ptr(resp.AccessToken)
	j.expireIn = resp.ExpiresIn
	return resp.AccessToken, nil
}
