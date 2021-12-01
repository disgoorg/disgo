package discord

import (
	"time"

	"github.com/DisgoOrg/disgo/json"
)

type AccessTokenExchange struct {
	AccessToken  string           `json:"access_token"`
	TokenType    TokenType        `json:"token_type"`
	ExpiresIn    time.Duration    `json:"expires_in"`
	RefreshToken string           `json:"refresh_token"`
	Scope        string           `json:"scope"`
	Webhook      *IncomingWebhook `json:"webhook"`
}

func (e *AccessTokenExchange) UnmarshalJSON(data []byte) error {
	type accessTokenExchange AccessTokenExchange
	var v struct {
		ExpiresIn int64 `json:"expires_in"`
		accessTokenExchange
	}
	if err := json.Unmarshal(data, v); err != nil {
		return err
	}

	*e = AccessTokenExchange(v.accessTokenExchange)
	e.ExpiresIn = time.Duration(v.ExpiresIn) * time.Second
	return nil
}

type GrantType string

const (
	GrantTypeAuthorizationCode GrantType = "authorization_code"
	GrantTypeRefreshToken      GrantType = "refresh_token"
)

func (t GrantType) String() string {
	return string(t)
}
