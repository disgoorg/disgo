package oauth2

import (
	"time"

	"github.com/disgoorg/disgo/discord"
)

var _ SessionController = (*sessionControllerImpl)(nil)

// SessionController lets you manage your Session(s)
type SessionController interface {
	// GetSession returns the Session for the given identifier or nil if none was found
	GetSession(identifier string) Session

	// CreateSession creates a new Session from the given identifier, access token, refresh token, scope, token type, expiration and webhook
	CreateSession(identifier string, accessToken string, refreshToken string, scopes []discord.ApplicationScope, tokenType discord.TokenType, expiration time.Time, webhook *discord.IncomingWebhook) Session

	// CreateSessionFromExchange creates a new Session from the given identifier and discord.AccessTokenExchange payload
	CreateSessionFromExchange(identifier string, exchange discord.AccessTokenExchange) Session
}

// NewSessionController returns a new empty SessionController
func NewSessionController() SessionController {
	return NewSessionControllerWithSessions(map[string]Session{})
}

// NewSessionControllerWithSessions returns a new SessionController with the given Session(s)
func NewSessionControllerWithSessions(sessions map[string]Session) SessionController {
	return &sessionControllerImpl{sessions: sessions}
}

type sessionControllerImpl struct {
	sessions map[string]Session
}

func (c *sessionControllerImpl) GetSession(identifier string) Session {
	return c.sessions[identifier]
}

func (c *sessionControllerImpl) CreateSession(identifier string, accessToken string, refreshToken string, scopes []discord.ApplicationScope, tokenType discord.TokenType, expiration time.Time, webhook *discord.IncomingWebhook) Session {
	session := &sessionImpl{
		accessToken:  accessToken,
		refreshToken: refreshToken,
		scopes:       scopes,
		tokenType:    tokenType,
		expiration:   expiration,
		webhook:      webhook,
	}

	c.sessions[identifier] = session

	return session
}

func (c *sessionControllerImpl) CreateSessionFromExchange(identifier string, exchange discord.AccessTokenExchange) Session {
	return c.CreateSession(identifier, exchange.AccessToken, exchange.RefreshToken, discord.SplitScopes(exchange.Scope), exchange.TokenType, time.Now().Add(exchange.ExpiresIn*time.Second), exchange.Webhook)
}
