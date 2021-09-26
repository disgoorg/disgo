package oauth2

import (
	"time"

	"github.com/DisgoOrg/disgo/discord"
)

var _ Session = (*sessionImpl)(nil)

type Session interface {
	AccessToken() string
	RefreshToken() string
	Scopes() []discord.ApplicationScope
	TokenType() discord.TokenType
	Expiration() time.Time
}

type sessionImpl struct {
	accessToken  string
	refreshToken string
	scopes       []discord.ApplicationScope
	tokenType    discord.TokenType
	expiration   time.Time
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
