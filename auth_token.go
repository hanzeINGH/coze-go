package coze

import (
	"context"
	"time"
)

type Auth interface {
	Token(ctx context.Context) (string, error)
}

var (
	_ Auth = &tokenAuthImpl{}
	_ Auth = &jwtOAuthImpl{}
)

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

func NewJWTAuth(client *JWTOAuthClient, opt *GetJWTAccessTokenReq) Auth {
	ttl := 900
	// default refresh token before expire in 30 seconds
	refreshBefore := int64(30)
	if opt == nil {
		return &jwtOAuthImpl{
			TTL:           ttl,
			client:        client,
			refreshBefore: refreshBefore,
		}
	}
	if opt.TTL > 0 {
		ttl = opt.TTL
	}
	if opt.RefreshBefore > 0 {
		refreshBefore = opt.RefreshBefore
	}
	return &jwtOAuthImpl{
		TTL:           ttl,
		Scope:         opt.Scope,
		SessionName:   opt.SessionName,
		refreshBefore: refreshBefore,
		client:        client,
		accountID:     opt.AccountID,
	}
}

// Token returns the access token.
func (r *tokenAuthImpl) Token(ctx context.Context) (string, error) {
	return r.accessToken, nil
}

type jwtOAuthImpl struct {
	TTL           int
	SessionName   *string
	Scope         *Scope
	client        *JWTOAuthClient
	accessToken   *string
	expireIn      int64
	refreshBefore int64 // refresh moment before expireIn, unit second
	refreshAt     int64
	accountID     *int64
}

func (r *jwtOAuthImpl) needRefresh() bool {
	return r.accessToken == nil || time.Now().Unix() > r.refreshAt
}

func (r *jwtOAuthImpl) Token(ctx context.Context) (string, error) {
	if !r.needRefresh() {
		return ptrValue(r.accessToken), nil
	}
	resp, err := r.client.GetAccessToken(ctx, &GetJWTAccessTokenReq{
		TTL:         r.TTL,
		SessionName: r.SessionName,
		Scope:       r.Scope,
		AccountID:   r.accountID,
	})
	if err != nil {
		return "", err
	}
	r.accessToken = ptr(resp.AccessToken)
	r.expireIn = resp.ExpiresIn
	r.refreshAt = resp.ExpiresIn - r.refreshBefore
	return resp.AccessToken, nil
}
