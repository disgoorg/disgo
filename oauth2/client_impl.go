package oauth2

import (
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/disgo/rest/route"
	"github.com/disgoorg/snowflake"
)

// New returns a new OAuth2 client
func New(id snowflake.Snowflake, secret string, opts ...ConfigOpt) Client {
	config := DefaultConfig()
	config.Apply(opts)

	return &ClientImpl{id: id, secret: secret, config: *config}
}

// ClientImpl is an OAuth2 client
type ClientImpl struct {
	id     snowflake.Snowflake
	secret string
	config Config
}

// ID returns the configured client ID
func (c *ClientImpl) ID() snowflake.Snowflake {
	return c.id
}

// Secret returns the configured client secret
func (c *ClientImpl) Secret() string {
	return c.secret
}

// Rest returns the underlying rest.OAuth2
func (c *ClientImpl) Rest() rest.OAuth2 {
	return c.config.OAuth2
}

// SessionController returns the configured SessionController
func (c *ClientImpl) SessionController() SessionController {
	return c.config.SessionController
}

// StateController returns the configured StateController
func (c *ClientImpl) StateController() StateController {
	return c.config.StateController
}

// GenerateAuthorizationURL generates an authorization URL with the given redirect URI, permissions, guildID, disableGuildSelect & scopes, state is automatically generated
func (c *ClientImpl) GenerateAuthorizationURL(redirectURI string, permissions discord.Permissions, guildID snowflake.Snowflake, disableGuildSelect bool, scopes ...discord.ApplicationScope) string {
	url, _ := c.GenerateAuthorizationURLState(redirectURI, permissions, guildID, disableGuildSelect, scopes...)
	return url
}

// GenerateAuthorizationURLState generates an authorization URL with the given redirect URI, permissions, guildID, disableGuildSelect & scopes, state is automatically generated & returned
func (c *ClientImpl) GenerateAuthorizationURLState(redirectURI string, permissions discord.Permissions, guildID snowflake.Snowflake, disableGuildSelect bool, scopes ...discord.ApplicationScope) (string, string) {
	state := c.StateController().GenerateNewState(redirectURI)
	values := route.QueryValues{
		"client_id":     c.ID(),
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

// StartSession starts a new session with the given authorization code & state
func (c *ClientImpl) StartSession(code string, state string, identifier string, opts ...rest.RequestOpt) (Session, error) {
	redirectURI := c.StateController().ConsumeState(state)
	if redirectURI == "" {
		return nil, ErrStateNotFound
	}
	exchange, err := c.Rest().GetAccessToken(c.id, c.secret, code, redirectURI, opts...)
	if err != nil {
		return nil, err
	}

	return c.SessionController().CreateSessionFromExchange(identifier, *exchange), nil
}

// RefreshSession refreshes the given session with the refresh token
func (c *ClientImpl) RefreshSession(identifier string, session Session, opts ...rest.RequestOpt) (Session, error) {
	exchange, err := c.Rest().RefreshAccessToken(c.id, c.secret, session.RefreshToken(), opts...)
	if err != nil {
		return nil, err
	}
	return c.SessionController().CreateSessionFromExchange(identifier, *exchange), nil
}

// GetUser returns the discord.OAuth2User associated with the given session. Fields filled in the struct depend on the Session.Scopes
func (c *ClientImpl) GetUser(session Session, opts ...rest.RequestOpt) (*discord.OAuth2User, error) {
	if session.Expiration().Before(time.Now()) {
		return nil, ErrAccessTokenExpired
	}
	if !discord.HasScope(discord.ApplicationScopeIdentify, session.Scopes()...) {
		return nil, ErrMissingOAuth2Scope(discord.ApplicationScopeIdentify)
	}

	return c.Rest().GetCurrentUser(session.AccessToken(), opts...)
}

// GetGuilds returns the discord.OAuth2Guild(s) the user is a member of. This requires the discord.ApplicationScopeGuilds scope in the session
func (c *ClientImpl) GetGuilds(session Session, opts ...rest.RequestOpt) ([]discord.OAuth2Guild, error) {
	if session.Expiration().Before(time.Now()) {
		return nil, ErrAccessTokenExpired
	}
	if !discord.HasScope(discord.ApplicationScopeGuilds, session.Scopes()...) {
		return nil, ErrMissingOAuth2Scope(discord.ApplicationScopeGuilds)
	}

	return c.Rest().GetCurrentUserGuilds(session.AccessToken(), "", "", 0, opts...)
}

// GetConnections returns the discord.Connection(s) the user has connected. This requires the discord.ApplicationScopeConnections scope in the session
func (c *ClientImpl) GetConnections(session Session, opts ...rest.RequestOpt) ([]discord.Connection, error) {
	if session.Expiration().Before(time.Now()) {
		return nil, ErrAccessTokenExpired
	}
	if !discord.HasScope(discord.ApplicationScopeConnections, session.Scopes()...) {
		return nil, ErrMissingOAuth2Scope(discord.ApplicationScopeConnections)
	}

	return c.Rest().GetCurrentUserConnections(session.AccessToken(), opts...)
}
