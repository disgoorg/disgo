package discord

import (
	"time"

	"github.com/DisgoOrg/disgo/json"
)

type AccessTokenExchange struct {
	AccessToken  string        `json:"access_token"`
	TokenType    TokenType     `json:"token_type"`
	ExpiresIn    time.Duration `json:"expires_in"`
	RefreshToken string        `json:"refresh_token"`
	Scope        string        `json:"scope"`
	Webhook      Webhook       `json:"webhook"`
}

func (e *AccessTokenExchange) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &e)
	if err != nil {
		return err
	}
	e.ExpiresIn = e.ExpiresIn * time.Second
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
