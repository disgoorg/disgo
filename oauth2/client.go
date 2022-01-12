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

// errors returned by the OAuth2 client
var (
	ErrStateNotFound      = errors.New("state could not be found")
	ErrAccessTokenExpired = errors.New("access token expired. refresh the session")
	ErrMissingOAuth2Scope = func(scope discord.ApplicationScope) error {
		return fmt.Errorf("missing '%s' scope", scope)
	}
)

type Client interface {
	// ID returns the configured client ID
	ID() discord.Snowflake
	// Secret returns the configured client secret
	Secret() string
	// Config returns the configured Config
	Config() Config

	// SessionController returns the configured SessionController
	SessionController() SessionController
	// StateController returns the configured StateController
	StateController() StateController

	// GenerateAuthorizationURL generates an authorization URL with the given redirect URI, permissions, guildID, disableGuildSelect & scopes. State is automatically generated
	GenerateAuthorizationURL(redirectURI string, permissions discord.Permissions, guildID discord.Snowflake, disableGuildSelect bool, scopes ...discord.ApplicationScope) string
	// GenerateAuthorizationURLState generates an authorization URL with the given redirect URI, permissions, guildID, disableGuildSelect & scopes. State is automatically generated & returned
	GenerateAuthorizationURLState(redirectURI string, permissions discord.Permissions, guildID discord.Snowflake, disableGuildSelect bool, scopes ...discord.ApplicationScope) (string, string)

	// StartSession starts a new Session with the given authorization code & state
	StartSession(code string, state string, identifier string, opts ...rest.RequestOpt) (Session, error)
	// RefreshSession refreshes the given Session with the refresh token
	RefreshSession(identifier string, session Session, opts ...rest.RequestOpt) (Session, error)

	// GetUser returns the discord.OAuth2User associated with the given Session. Fields filled in the struct depend on the Session.Scopes
	GetUser(session Session, opts ...rest.RequestOpt) (*discord.OAuth2User, error)
	// GetGuilds returns the discord.OAuth2Guild(s) the user is a member of. This requires the discord.ApplicationScopeGuilds scope in the Session
	GetGuilds(session Session, opts ...rest.RequestOpt) ([]discord.OAuth2Guild, error)
	// GetConnections returns the discord.Connection(s) the user has connected. This requires the discord.ApplicationScopeConnections scope in the Session
	GetConnections(session Session, opts ...rest.RequestOpt) ([]discord.Connection, error)
}

// New returns a new OAuth2 client
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
		config.StateController = NewStateController(config.StateControllerConfig)
	}

	return &clientImpl{id: id, secret: secret, config: *config}
}

// Client is an OAuth2 client
type clientImpl struct {
	id     discord.Snowflake
	secret string
	config Config
}

// ID returns the configured client ID
func (c *clientImpl) ID() discord.Snowflake {
	return c.id
}

// Secret returns the configured client secret
func (c *clientImpl) Secret() string {
	return c.secret
}

// Config returns the configured Config
func (c *clientImpl) Config() Config {
	return c.config
}

// SessionController returns the configured SessionController
func (c *clientImpl) SessionController() SessionController {
	return c.config.SessionController
}

// StateController returns the configured StateController
func (c *clientImpl) StateController() StateController {
	return c.config.StateController
}

// GenerateAuthorizationURL generates an authorization URL with the given redirect URI, permissions, guildID, disableGuildSelect & scopes. State is automatically generated
func (c *clientImpl) GenerateAuthorizationURL(redirectURI string, permissions discord.Permissions, guildID discord.Snowflake, disableGuildSelect bool, scopes ...discord.ApplicationScope) string {
	url, _ := c.GenerateAuthorizationURLState(redirectURI, permissions, guildID, disableGuildSelect, scopes...)
	return url
}

// GenerateAuthorizationURLState generates an authorization URL with the given redirect URI, permissions, guildID, disableGuildSelect & scopes. State is automatically generated & returned
func (c *clientImpl) GenerateAuthorizationURLState(redirectURI string, permissions discord.Permissions, guildID discord.Snowflake, disableGuildSelect bool, scopes ...discord.ApplicationScope) (string, string) {
	state := c.config.StateController.GenerateNewState(redirectURI)
	values := route.QueryValues{
		"client_id":     c.id,
		"redirect_uri":  redirectURI,
		"response_type": "code",
		"scope":         discord.JoinScopes(scopes),
		"state":         state,
	}
	if permissions != discord.PermissionsNone {
		values["permissions"] = permissions
	}
	if guildID != "" {
		values["guild_id"] = guildID
	}
	if disableGuildSelect {
		values["disable_guild_select"] = true
	}
	compiledRoute, _ := route.Authorize.Compile(values)
	return compiledRoute.URL(), state
}

// StartSession starts a new Session with the given authorization code & state
func (c *clientImpl) StartSession(code string, state string, identifier string, opts ...rest.RequestOpt) (Session, error) {
	redirectURI := c.config.StateController.ConsumeState(state)
	if redirectURI == "" {
		return nil, ErrStateNotFound
	}
	exchange, err := c.config.OAuth2Service.GetAccessToken(c.id, c.secret, code, redirectURI, opts...)
	if err != nil {
		return nil, err
	}

	return c.config.SessionController.CreateSessionFromExchange(identifier, *exchange), nil
}

// RefreshSession refreshes the given Session with the refresh token
func (c *clientImpl) RefreshSession(identifier string, session Session, opts ...rest.RequestOpt) (Session, error) {
	exchange, err := c.config.OAuth2Service.RefreshAccessToken(c.id, c.secret, session.RefreshToken(), opts...)
	if err != nil {
		return nil, err
	}
	return c.config.SessionController.CreateSessionFromExchange(identifier, *exchange), nil
}

// GetUser returns the discord.OAuth2User associated with the given Session. Fields filled in the struct depend on the Session.Scopes
func (c *clientImpl) GetUser(session Session, opts ...rest.RequestOpt) (*discord.OAuth2User, error) {
	if session.Expiration().Before(time.Now()) {
		return nil, ErrAccessTokenExpired
	}
	if !discord.HasScope(discord.ApplicationScopeIdentify, session.Scopes()...) {
		return nil, ErrMissingOAuth2Scope(discord.ApplicationScopeIdentify)
	}

	return c.config.OAuth2Service.GetCurrentUser(session.AccessToken(), opts...)
}

// GetGuilds returns the discord.OAuth2Guild(s) the user is a member of. This requires the discord.ApplicationScopeGuilds scope in the Session
func (c *clientImpl) GetGuilds(session Session, opts ...rest.RequestOpt) ([]discord.OAuth2Guild, error) {
	if session.Expiration().Before(time.Now()) {
		return nil, ErrAccessTokenExpired
	}
	if !discord.HasScope(discord.ApplicationScopeGuilds, session.Scopes()...) {
		return nil, ErrMissingOAuth2Scope(discord.ApplicationScopeGuilds)
	}

	return c.config.OAuth2Service.GetCurrentUserGuilds(session.AccessToken(), "", "", 0, opts...)
}

// GetConnections returns the discord.Connection(s) the user has connected. This requires the discord.ApplicationScopeConnections scope in the Session
func (c *clientImpl) GetConnections(session Session, opts ...rest.RequestOpt) ([]discord.Connection, error) {
	if session.Expiration().Before(time.Now()) {
		return nil, ErrAccessTokenExpired
	}
	if !discord.HasScope(discord.ApplicationScopeConnections, session.Scopes()...) {
		return nil, ErrMissingOAuth2Scope(discord.ApplicationScopeConnections)
	}

	return c.config.OAuth2Service.GetCurrentUserConnections(session.AccessToken(), opts...)
}
