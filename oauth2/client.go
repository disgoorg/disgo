package oauth2

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

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

func New(id discord.Snowflake, secret string, opts ...ConfigOpt) *Client {
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

	return &Client{ID: id, Secret: secret, Config: *config}
}

type Client struct {
	ID     discord.Snowflake
	Secret string
	Config
}

func (c *Client) GenerateAuthorizationURL(redirectURI string, scopes ...discord.ApplicationScope) string {
	values := route.QueryValues{
		"client_id":     c.ID,
		"redirect_uri":  redirectURI,
		"response_type": "code",
		"scope":         discord.JoinScopes(scopes),
		"state":         c.StateController.GenerateNewState(redirectURI),
	}
	compiledRoute, _ := route.Authorize.Compile(values)
	return compiledRoute.URL()
}

func (c *Client) StartSession(code string, state string, identifier string, opts ...rest.RequestOpt) (Session, discord.Webhook, error) {
	redirectURI := c.StateController.ConsumeState(state)
	if redirectURI == nil {
		return nil, nil, ErrStateNotFound
	}
	exchange, err := c.OAuth2Service.GetAccessToken(c.ID, c.Secret, code, *redirectURI, opts...)
	if err != nil {
		return nil, nil, err
	}

	return c.SessionController.CreateSession(identifier, exchange.AccessToken, exchange.RefreshToken, discord.SplitScopes(exchange.Scope), exchange.TokenType, time.Now().Add(exchange.ExpiresIn*time.Second)), exchange.Webhook, nil
}

func (c *Client) RefreshSession(identifier string, session Session, opts ...rest.RequestOpt) (Session, error) {
	exchange, err := c.OAuth2Service.RefreshAccessToken(c.ID, c.Secret, session.RefreshToken(), opts...)
	if err != nil {
		return nil, err
	}
	return c.SessionController.CreateSession(identifier, exchange.AccessToken, exchange.RefreshToken, discord.SplitScopes(exchange.Scope), exchange.TokenType, time.Now().Add(exchange.ExpiresIn*time.Second)), nil
}

func (c *Client) GetUser(session Session, opts ...rest.RequestOpt) (*discord.OAuth2User, error) {
	if session.Expiration().Before(time.Now()) {
		return nil, ErrAccessTokenExpired
	}
	if !discord.HasScope(discord.ApplicationScopeIdentify, session.Scopes()...) {
		return nil, ErrMissingOAuth2Scope(discord.ApplicationScopeIdentify)
	}

	return c.OAuth2Service.GetCurrentUser(session.AccessToken(), opts...)
}

func (c *Client) GetGuilds(session Session, opts ...rest.RequestOpt) ([]discord.OAuth2Guild, error) {
	if session.Expiration().Before(time.Now()) {
		return nil, ErrAccessTokenExpired
	}
	if !discord.HasScope(discord.ApplicationScopeGuilds, session.Scopes()...) {
		return nil, ErrMissingOAuth2Scope(discord.ApplicationScopeGuilds)
	}

	return c.OAuth2Service.GetCurrentUserGuilds(session.AccessToken(), opts...)
}

func (c *Client) GetConnections(session Session, opts ...rest.RequestOpt) ([]discord.Connection, error) {
	if session.Expiration().Before(time.Now()) {
		return nil, ErrAccessTokenExpired
	}
	if !discord.HasScope(discord.ApplicationScopeConnections, session.Scopes()...) {
		return nil, ErrMissingOAuth2Scope(discord.ApplicationScopeConnections)
	}

	return c.OAuth2Service.GetCurrentUserConnections(session.AccessToken(), opts...)
}
