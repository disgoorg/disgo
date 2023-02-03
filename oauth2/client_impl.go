package oauth2

import (
	"time"

	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
)

// New returns a new OAuth2 client with the given ID, secret and ConfigOpt(s).
func New(id snowflake.ID, secret string, opts ...ConfigOpt) Client {
	config := DefaultConfig()
	config.Apply(opts)

	return &clientImpl{id: id, secret: secret, config: *config}
}

type clientImpl struct {
	id     snowflake.ID
	secret string
	config Config
}

func (c *clientImpl) ID() snowflake.ID {
	return c.id
}

func (c *clientImpl) Secret() string {
	return c.secret
}

func (c *clientImpl) Rest() rest.OAuth2 {
	return c.config.OAuth2
}

func (c *clientImpl) SessionController() SessionController {
	return c.config.SessionController
}

func (c *clientImpl) StateController() StateController {
	return c.config.StateController
}

func (c *clientImpl) GenerateAuthorizationURL(redirectURI string, permissions discord.Permissions, guildID snowflake.ID, disableGuildSelect bool, scopes ...discord.OAuth2Scope) string {
	authURL, _ := c.GenerateAuthorizationURLState(redirectURI, permissions, guildID, disableGuildSelect, scopes...)
	return authURL
}

func (c *clientImpl) GenerateAuthorizationURLState(redirectURI string, permissions discord.Permissions, guildID snowflake.ID, disableGuildSelect bool, scopes ...discord.OAuth2Scope) (string, string) {
	state := c.StateController().GenerateNewState(redirectURI)
	values := discord.QueryValues{
		"client_id":     c.id,
		"redirect_uri":  redirectURI,
		"response_type": "code",
		"scope":         discord.JoinScopes(scopes),
		"state":         state,
	}

	if permissions != discord.PermissionsNone {
		values["permissions"] = permissions
	}
	if guildID != 0 {
		values["guild_id"] = guildID
	}
	if disableGuildSelect {
		values["disable_guild_select"] = true
	}
	return discord.AuthorizeURL(values), state
}

func (c *clientImpl) StartSession(code string, state string, identifier string, opts ...rest.RequestOpt) (Session, error) {
	redirectURI := c.StateController().ConsumeState(state)
	if redirectURI == "" {
		return nil, ErrStateNotFound
	}
	exchange, err := c.Rest().GetAccessToken(c.id, c.secret, code, redirectURI, opts...)
	if err != nil {
		return nil, err
	}
	return c.SessionController().CreateSessionFromResponse(identifier, *exchange), nil
}

func (c *clientImpl) RefreshSession(identifier string, session Session, opts ...rest.RequestOpt) (Session, error) {
	exchange, err := c.Rest().RefreshAccessToken(c.id, c.secret, session.RefreshToken(), opts...)
	if err != nil {
		return nil, err
	}
	return c.SessionController().CreateSessionFromResponse(identifier, *exchange), nil
}

func (c *clientImpl) GetUser(session Session, opts ...rest.RequestOpt) (*discord.OAuth2User, error) {
	if err := checkSession(session, discord.OAuth2ScopeIdentify); err != nil {
		return nil, err
	}
	return c.Rest().GetCurrentUser(session.AccessToken(), opts...)
}

func (c *clientImpl) GetMember(session Session, guildID snowflake.ID, opts ...rest.RequestOpt) (*discord.Member, error) {
	if err := checkSession(session, discord.OAuth2ScopeGuildsMembersRead); err != nil {
		return nil, err
	}
	return c.Rest().GetCurrentMember(session.AccessToken(), guildID, opts...)
}

func (c *clientImpl) GetGuilds(session Session, opts ...rest.RequestOpt) ([]discord.OAuth2Guild, error) {
	if err := checkSession(session, discord.OAuth2ScopeGuilds); err != nil {
		return nil, err
	}
	return c.Rest().GetCurrentUserGuilds(session.AccessToken(), 0, 0, 0, opts...)
}

func (c *clientImpl) GetConnections(session Session, opts ...rest.RequestOpt) ([]discord.Connection, error) {
	if err := checkSession(session, discord.OAuth2ScopeConnections); err != nil {
		return nil, err
	}
	return c.Rest().GetCurrentUserConnections(session.AccessToken(), opts...)
}

func (c *clientImpl) GetApplicationRoleConnection(session Session, applicationID snowflake.ID, opts ...rest.RequestOpt) (*discord.ApplicationRoleConnection, error) {
	if err := checkSession(session, discord.OAuth2ScopeRoleConnectionsWrite); err != nil {
		return nil, err
	}
	return c.Rest().GetCurrentUserApplicationRoleConnection(session.AccessToken(), applicationID, opts...)
}

func (c *clientImpl) UpdateApplicationRoleConnection(session Session, applicationID snowflake.ID, update discord.ApplicationRoleConnectionUpdate, opts ...rest.RequestOpt) (*discord.ApplicationRoleConnection, error) {
	if err := checkSession(session, discord.OAuth2ScopeRoleConnectionsWrite); err != nil {
		return nil, err
	}
	return c.Rest().UpdateCurrentUserApplicationRoleConnection(session.AccessToken(), applicationID, update, opts...)
}

func checkSession(session Session, scope discord.OAuth2Scope) error {
	if session.Expiration().Before(time.Now()) {
		return ErrAccessTokenExpired
	}
	if !discord.HasScope(scope, session.Scopes()...) {
		return ErrMissingOAuth2Scope(scope)
	}
	return nil
}
