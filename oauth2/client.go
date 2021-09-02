package oauth2

import (
	"errors"
	"fmt"
	"time"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/disgo/rest/route"
	"github.com/DisgoOrg/log"
)

var (
	ErrStateNotFound      = errors.New("state could not be found")
	ErrAccessTokenExpired = errors.New("access token expired. refresh the session")
	ErrMissingOAuth2Scope = func(scope discord.ApplicationScope) error {
		return fmt.Errorf("missing '%s' scope", scope)
	}
)

// Client lets you edit/send WebhookMessage(s) or update/delete the Webhook
type Client interface {
	Logger() log.Logger

	RestClient() rest.Client
	OAuth2Service() rest.OAuth2Service
	StateController() StateController
	SessionController() SessionController

	GenerateAuthorizationURL(redirectURI string, scopes ...discord.ApplicationScope) string
	StartSession(code string, state string, identifier string, opts ...rest.RequestOpt) (Session, *discord.Webhook, rest.Error)
	RefreshSession(identifier string, session Session, opts ...rest.RequestOpt) (Session, rest.Error)

	GetUser(session Session, opts ...rest.RequestOpt) (*discord.OAuth2User, rest.Error)
	GetGuilds(session Session, opts ...rest.RequestOpt) ([]discord.PartialGuild, rest.Error)
	GetConnections(session Session, opts ...rest.RequestOpt) ([]discord.Connection, rest.Error)

	ID() discord.Snowflake
	Secret() string
}

func New(id discord.Snowflake, secret string, opts ...ConfigOpt) Client {
	config := &DefaultConfig
	config.Apply(opts)

	if config.Logger == nil {
		config.Logger = log.Default()
	}

	if config.RestClient == nil {
		config.RestClient = rest.NewClient(config.RestClientConfig)
	}
	if config.OAuth2Service == nil {
		config.OAuth2Service = rest.NewOAuth2Service(config.RestClient)
	}
	if config.SessionController == nil {
		config.SessionController = NewSessionController()
	}
	if config.StateController == nil {
		config.StateController = NewStateController()
	}

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
		"scope":         discord.JoinScopes(scopes),
		"state":         c.StateController().GenerateNewState(redirectURI),
	}
	compiledRoute, _ := route.Authorize.Compile(values)
	return compiledRoute.URL()
}

func (c *clientImpl) StartSession(code string, state string, identifier string, opts ...rest.RequestOpt) (Session, *discord.Webhook, rest.Error) {
	redirectURI := c.StateController().ConsumeState(state)
	if redirectURI == nil {
		return nil, nil, rest.NewError(nil, ErrStateNotFound)
	}
	exchange, err := c.OAuth2Service().GetAccessToken(c.ID(), c.Secret(), code, *redirectURI, opts...)
	if err != nil {
		return nil, nil, err
	}

	return c.SessionController().CreateSession(identifier, exchange.AccessToken, exchange.RefreshToken, discord.SplitScopes(exchange.Scope), exchange.TokenType, time.Now().Add(exchange.ExpiresIn*time.Second)), exchange.Webhook, nil
}

func (c *clientImpl) RefreshSession(identifier string, session Session, opts ...rest.RequestOpt) (Session, rest.Error) {
	exchange, err := c.OAuth2Service().RefreshAccessToken(c.ID(), c.Secret(), session.RefreshToken(), opts...)
	if err != nil {
		return nil, err
	}
	return c.SessionController().CreateSession(identifier, exchange.AccessToken, exchange.RefreshToken, discord.SplitScopes(exchange.Scope), exchange.TokenType, time.Now().Add(exchange.ExpiresIn*time.Second)), nil
}

func (c *clientImpl) GetUser(session Session, opts ...rest.RequestOpt) (*discord.OAuth2User, rest.Error) {
	if session.Expiration().Before(time.Now()) {
		return nil, rest.NewError(nil, ErrAccessTokenExpired)
	}
	if !discord.HasScope(discord.ApplicationScopeIdentify, session.Scopes()...) {
		return nil, rest.NewError(nil, ErrMissingOAuth2Scope(discord.ApplicationScopeIdentify))
	}

	return c.OAuth2Service().GetCurrentUser(session.AccessToken(), opts...)
}

func (c *clientImpl) GetGuilds(session Session, opts ...rest.RequestOpt) ([]discord.PartialGuild, rest.Error) {
	if session.Expiration().Before(time.Now()) {
		return nil, rest.NewError(nil, ErrAccessTokenExpired)
	}
	if !discord.HasScope(discord.ApplicationScopeGuilds, session.Scopes()...) {
		return nil, rest.NewError(nil, ErrMissingOAuth2Scope(discord.ApplicationScopeGuilds))
	}

	return c.OAuth2Service().GetCurrentUserGuilds(session.AccessToken(), opts...)
}

func (c *clientImpl) GetConnections(session Session, opts ...rest.RequestOpt) ([]discord.Connection, rest.Error) {
	if session.Expiration().Before(time.Now()) {
		return nil, rest.NewError(nil, ErrAccessTokenExpired)
	}
	if !discord.HasScope(discord.ApplicationScopeConnections, session.Scopes()...) {
		return nil, rest.NewError(nil, ErrMissingOAuth2Scope(discord.ApplicationScopeConnections))
	}

	return c.OAuth2Service().GetCurrentUserConnections(session.AccessToken(), opts...)
}

func (c *clientImpl) ID() discord.Snowflake {
	return c.id
}
func (c *clientImpl) Secret() string {
	return c.secret
}
