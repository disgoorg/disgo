package oauth2

import (
	"time"

	"github.com/DisgoOrg/disgo/discord"
)

var _ SessionController = (*sessionControllerImpl)(nil)

type SessionController interface {
	GetSession(identifier string) Session
	CreateSession(identifier string, accessToken string, refreshToken string, scopes []discord.ApplicationScope, tokenType discord.TokenType, expiration time.Time) Session
}

func NewSessionController() SessionController {
	return &sessionControllerImpl{sessions: map[string]Session{}}
}

type sessionControllerImpl struct {
	sessions map[string]Session
}

func (c *sessionControllerImpl) GetSession(identifier string) Session {
	return c.sessions[identifier]
}
func (c *sessionControllerImpl) CreateSession(identifier string, accessToken string, refreshToken string, scopes []discord.ApplicationScope, tokenType discord.TokenType, expiration time.Time) Session {
	session := &sessionImpl{
		accessToken:  accessToken,
		refreshToken: refreshToken,
		scopes:       scopes,
		tokenType:    tokenType,
		expiration:   expiration,
	}

	c.sessions[identifier] = session

	return session
}
