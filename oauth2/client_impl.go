package oauth2

import (
	"log/slog"
	"time"

	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
)

// New returns a new OAuth2 client with the given ID, secret and ConfigOpt(s).
func New(id snowflake.ID, secret string, opts ...ConfigOpt) Client {
	config := DefaultConfig()
	config.Apply(opts)
	config.Logger = config.Logger.With(slog.String("name", "oauth2"))

	return &clientImpl{
		id:              id,
		secret:          secret,
		oAuth2:          config.OAuth2,
		stateController: config.StateController,
	}
}

type clientImpl struct {
	id              snowflake.ID
	secret          string
	oAuth2          rest.OAuth2
	stateController StateController
}

func (c *clientImpl) ID() snowflake.ID {
	return c.id
}

func (c *clientImpl) Secret() string {
	return c.secret
}

func (c *clientImpl) Rest() rest.OAuth2 {
	return c.oAuth2
}

func (c *clientImpl) StateController() StateController {
	return c.stateController
}

func (c *clientImpl) GenerateAuthorizationURL(params AuthorizationURLParams) string {
	authURL, _ := c.GenerateAuthorizationURLState(params)
	return authURL
}

func (c *clientImpl) GenerateAuthorizationURLState(params AuthorizationURLParams) (string, string) {
	state := c.StateController().NewState(params.RedirectURI)
	values := discord.QueryValues{
		"client_id":     c.id,
		"redirect_uri":  params.RedirectURI,
		"response_type": "code",
		"scope":         discord.JoinScopes(params.Scopes),
		"state":         state,
	}

	if params.Permissions != discord.PermissionsNone {
		values["permissions"] = params.Permissions
	}
	if params.GuildID != 0 {
		values["guild_id"] = params.GuildID
	}
	if params.DisableGuildSelect {
		values["disable_guild_select"] = true
	}
	if params.IntegrationType != 0 {
		values["integration_type"] = params.IntegrationType
	}
	return discord.AuthorizeURL(values), state
}

func (c *clientImpl) StartSession(code string, state string, opts ...rest.RequestOpt) (Session, *discord.IncomingWebhook, error) {
	redirectURI := c.StateController().UseState(state)
	if redirectURI == "" {
		return Session{}, nil, ErrStateNotFound
	}
	accessToken, err := c.Rest().GetAccessToken(c.id, c.secret, code, redirectURI, opts...)
	if err != nil {
		return Session{}, nil, err
	}

	return newSession(*accessToken), accessToken.Webhook, nil
}

func (c *clientImpl) RefreshSession(session Session, opts ...rest.RequestOpt) (Session, error) {
	accessToken, err := c.Rest().RefreshAccessToken(c.id, c.secret, session.RefreshToken, opts...)
	if err != nil {
		return Session{}, err
	}
	return newSession(*accessToken), nil
}

func (c *clientImpl) VerifySession(session Session, opts ...rest.RequestOpt) (Session, error) {
	if session.Expired() {
		return c.RefreshSession(session, opts...)
	}
	return session, nil
}

func (c *clientImpl) GetUser(session Session, opts ...rest.RequestOpt) (*discord.OAuth2User, error) {
	if err := checkSession(session, discord.OAuth2ScopeIdentify); err != nil {
		return nil, err
	}
	return c.Rest().GetCurrentUser(session.AccessToken, opts...)
}

func (c *clientImpl) GetMember(session Session, guildID snowflake.ID, opts ...rest.RequestOpt) (*discord.Member, error) {
	if err := checkSession(session, discord.OAuth2ScopeGuildsMembersRead); err != nil {
		return nil, err
	}
	return c.Rest().GetCurrentMember(session.AccessToken, guildID, opts...)
}

func (c *clientImpl) GetGuilds(session Session, opts ...rest.RequestOpt) ([]discord.OAuth2Guild, error) {
	if err := checkSession(session, discord.OAuth2ScopeGuilds); err != nil {
		return nil, err
	}
	return c.Rest().GetCurrentUserGuilds(session.AccessToken, 0, 0, 0, false, opts...)
}

func (c *clientImpl) GetConnections(session Session, opts ...rest.RequestOpt) ([]discord.Connection, error) {
	if err := checkSession(session, discord.OAuth2ScopeConnections); err != nil {
		return nil, err
	}
	return c.Rest().GetCurrentUserConnections(session.AccessToken, opts...)
}

func (c *clientImpl) GetApplicationRoleConnection(session Session, applicationID snowflake.ID, opts ...rest.RequestOpt) (*discord.ApplicationRoleConnection, error) {
	if err := checkSession(session, discord.OAuth2ScopeRoleConnectionsWrite); err != nil {
		return nil, err
	}
	return c.Rest().GetCurrentUserApplicationRoleConnection(session.AccessToken, applicationID, opts...)
}

func (c *clientImpl) UpdateApplicationRoleConnection(session Session, applicationID snowflake.ID, update discord.ApplicationRoleConnectionUpdate, opts ...rest.RequestOpt) (*discord.ApplicationRoleConnection, error) {
	if err := checkSession(session, discord.OAuth2ScopeRoleConnectionsWrite); err != nil {
		return nil, err
	}
	return c.Rest().UpdateCurrentUserApplicationRoleConnection(session.AccessToken, applicationID, update, opts...)
}

func checkSession(session Session, scope discord.OAuth2Scope) error {
	if session.Expired() {
		return ErrSessionExpired
	}
	if !discord.HasScope(scope, session.Scopes...) {
		return ErrMissingOAuth2Scope(scope)
	}
	return nil
}

func newSession(accessToken discord.AccessTokenResponse) Session {
	return Session{
		AccessToken:  accessToken.AccessToken,
		RefreshToken: accessToken.RefreshToken,
		Scopes:       accessToken.Scope,
		TokenType:    accessToken.TokenType,
		Expiration:   time.Now().Add(accessToken.ExpiresIn),
	}
}
