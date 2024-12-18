package coze

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
func (a *TokenAuth) Token() string {
	return a.accessToken
}
