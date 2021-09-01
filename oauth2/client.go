package oauth2

import (
	"time"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/disgo/rest/route"
	"github.com/DisgoOrg/log"
)

// Client lets you edit/send WebhookMessage(s) or update/delete the Webhook
type Client interface {
	Logger() log.Logger

	RestClient() rest.Client
	OAuth2Service() rest.OAuth2Service
	EntityBuilder() EntityBuilder
	StateController() StateController
	SessionController() SessionController

	GenerateAuthorizationURL(redirectURI string, scopes ...discord.ApplicationScope) string
	StartSession(code string, state string, identifier string, opts ...rest.RequestOpt) (Session, rest.Error)

	GetUser(session Session, opts ...rest.RequestOpt) (*User, rest.Error)
	GetGuilds(session Session, opts ...rest.RequestOpt) ([]*Guild, rest.Error)

	ID() discord.Snowflake
	Secret() string
}

func New(id discord.Snowflake, secret string, opts ...ConfigOpt) Client {
	config := &DefaultConfig
	config.Apply(opts)

	return &clientImpl{id: id, secret: secret, config: *config}
}

type clientImpl struct {
	id     discord.Snowflake
	secret string
	config Config
}

func (c *clientImpl) Logger() log.Logger {
	return c.config.Logger
}

func (c *clientImpl) RestClient() rest.Client {
	return c.config.RestClient
}
func (c *clientImpl) OAuth2Service() rest.OAuth2Service {
	return c.config.OAuth2Service
}
func (c *clientImpl) EntityBuilder() EntityBuilder {
	return c.config.EntityBuilder
}
func (c *clientImpl) StateController() StateController {
	return c.config.StateController
}
func (c *clientImpl) SessionController() SessionController {
	return c.config.SessionController
}

func (c *clientImpl) GenerateAuthorizationURL(redirectURI string, scopes ...discord.ApplicationScope) string {
	values := route.QueryValues{
		"client_id":     c.ID(),
		"redirect_uri":  redirectURI,
		"response_type": "code",
		"scope":         discord.JoinScopes(scopes...),
		"state":         c.StateController().GenerateNewState(redirectURI),
	}
	compiledRoute, _ := route.Authorize.Compile(values)
	return compiledRoute.URL()
}

func (c *clientImpl) StartSession(code string, state string, identifier string, opts ...rest.RequestOpt) (Session, rest.Error) {
	redirectURI := c.StateController().ConsumeState(state)
	if redirectURI == nil {
		return nil, rest.NewError(nil, ErrStateNotFound)
	}
	exchange, err := c.OAuth2Service().GetAccessToken(c.ID(), c.Secret(), code, *redirectURI, opts...)
	if err != nil {
		return nil, err
	}

	return c.SessionController().CreateSession(identifier, exchange.AccessToken, exchange.RefreshToken, discord.SplitScopes(exchange.Scope), exchange.TokenType, time.Now().Add(exchange.ExpiresIn*time.Second)), nil
}

func (c *clientImpl) GetUser(session Session, opts ...rest.RequestOpt) (*User, rest.Error) {
	user, err := c.OAuth2Service().GetCurrentUser(session.AccessToken(), opts...)
	if err != nil {
		return nil, err
	}
	return c.EntityBuilder().CreateUser(*user), nil
}

func (c *clientImpl) GetGuilds(session Session, opts ...rest.RequestOpt) ([]*Guild, rest.Error) {
	partialGuilds, err := c.OAuth2Service().GetCurrentUserGuilds(session.AccessToken(), opts...)
	if err != nil {
		return nil, err
	}
	guilds := make([]*Guild, len(partialGuilds))
	for i, guild := range partialGuilds {
		guilds[i] = c.EntityBuilder().CreateGuild(guild)
	}
	return guilds, nil
}

func (c *clientImpl) ID() discord.Snowflake {
	return c.id
}
func (c *clientImpl) Secret() string {
	return c.secret
}
