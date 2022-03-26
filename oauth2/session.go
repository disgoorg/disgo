package oauth2

import (
	"time"

	"github.com/disgoorg/disgo/discord"
)

var _ Session = (*sessionImpl)(nil)

// Session represents a discord access token response (https://discord.com/developers/docs/topics/oauth2#authorization-code-grant-access-token-response)
type Session interface {
	// AccessToken allows requesting user information
	AccessToken() string

	// RefreshToken allows refreshing the AccessToken
	RefreshToken() string

	// Scopes returns the discord.ApplicationScope(s) of the Session
	Scopes() []discord.ApplicationScope

	// TokenType returns the discord.TokenType of the AccessToken
	TokenType() discord.TokenType

	// Expiration returns the time.Time when the AccessToken expires and needs to be refreshed
	Expiration() time.Time

	// Webhook returns the discord.IncomingWebhook when the discord.ApplicationScopeWebhookIncoming is set
	Webhook() *discord.IncomingWebhook
}

type sessionImpl struct {
	accessToken  string
	refreshToken string
	scopes       []discord.ApplicationScope
	tokenType    discord.TokenType
	expiration   time.Time
	webhook      *discord.IncomingWebhook
}

func (s *sessionImpl) AccessToken() string {
	return s.accessToken
}

func (s *sessionImpl) RefreshToken() string {
	return s.refreshToken
}

func (s *sessionImpl) Scopes() []discord.ApplicationScope {
	return s.scopes
}

func (s *sessionImpl) TokenType() discord.TokenType {
	return s.tokenType
}

func (s *sessionImpl) Expiration() time.Time {
	return s.expiration
}

func (s *sessionImpl) Webhook() *discord.IncomingWebhook {
	return s.webhook
}
