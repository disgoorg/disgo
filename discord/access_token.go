package discord

import "time"

type AccessTokenExchange struct {
	AccessToken  string        `json:"access_token"`
	TokenType    TokenType     `json:"token_type"`
	ExpiresIn    time.Duration `json:"expires_in"`
	RefreshToken string        `json:"refresh_token"`
	Scope        string        `json:"scope"`
}

type GrantType string

const (
	GrantTypeAuthorizationCode GrantType = "authorization_code"
	GrantTypeRefreshToken      GrantType = "refresh_token"
)

func (t GrantType) String() string {
	return string(t)
}
