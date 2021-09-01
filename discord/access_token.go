package discord

import "time"

type AccessTokenExchange struct {
	AccessToken  string
	TokenType    TokenType
	ExpiresIn    time.Duration
	RefreshToken string
	Scope        string
}

type GrantType string

const (
	GrantTypeAuthorizationCode GrantType = "authorization_code"
	GrantTypeRefreshToken      GrantType = "refresh_token"
)
