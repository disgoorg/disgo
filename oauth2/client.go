package oauth2

import (
	"errors"
	"fmt"
	"time"

	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
)

var (
	// ErrStateNotFound is returned when the state is not found in the SessionController.
	ErrStateNotFound = errors.New("state could not be found")

	// ErrSessionExpired is returned when the Session has expired.
	ErrSessionExpired = errors.New("access token expired. refresh the session")

	// ErrMissingOAuth2Scope is returned when a specific OAuth2 scope is missing.
	ErrMissingOAuth2Scope = func(scope discord.OAuth2Scope) error {
		return fmt.Errorf("missing '%s' scope", scope)
	}
)

// Session represents a discord access token response (https://discord.com/developers/docs/topics/oauth2#authorization-code-grant-access-token-response)
type Session struct {
	// AccessToken allows requesting user information
	AccessToken string `json:"access_token"`

	// RefreshToken allows refreshing the AccessToken
	RefreshToken string `json:"refresh_token"`

	// Scopes returns the discord.OAuth2Scope(s) of the Session
	Scopes []discord.OAuth2Scope `json:"scopes"`

	// TokenType returns the discord.TokenType of the AccessToken
	TokenType discord.TokenType `json:"token_type"`

	// Expiration returns the time.Time when the AccessToken expires and needs to be refreshed
	Expiration time.Time `json:"expiration"`
}

func (s Session) Expired() bool {
	return s.Expiration.Before(time.Now())
}

type AuthorizationURLParams struct {
	RedirectURI        string
	Permissions        discord.Permissions
	GuildID            snowflake.ID
	DisableGuildSelect bool
	IntegrationType    discord.ApplicationIntegrationType
	Scopes             []discord.OAuth2Scope
}

// New returns a new OAuth2 client with the given ID, Secret and ConfigOpt(s).
func New(id snowflake.ID, secret string, opts ...ConfigOpt) *Client {
	cfg := defaultConfig()
	cfg.apply(opts)

	return &Client{
		ID:              id,
		Secret:          secret,
		Rest:            cfg.OAuth2,
		StateController: cfg.StateController,
	}
}

type Client struct {
	ID              snowflake.ID
	Secret          string
	Rest            rest.OAuth2
	StateController StateController
}

func (c *Client) GenerateAuthorizationURL(params AuthorizationURLParams) string {
	authURL, _ := c.GenerateAuthorizationURLState(params)
	return authURL
}

func (c *Client) GenerateAuthorizationURLState(params AuthorizationURLParams) (string, string) {
	state := c.StateController.NewState(params.RedirectURI)
	values := discord.QueryValues{
		"client_id":     c.ID,
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

func (c *Client) StartSession(code string, state string, opts ...rest.RequestOpt) (Session, *discord.IncomingWebhook, error) {
	redirectURI := c.StateController.UseState(state)
	if redirectURI == "" {
		return Session{}, nil, ErrStateNotFound
	}
	accessToken, err := c.Rest.GetAccessToken(c.ID, c.Secret, code, redirectURI, opts...)
	if err != nil {
		return Session{}, nil, err
	}

	return newSession(*accessToken), accessToken.Webhook, nil
}

func (c *Client) RefreshSession(session Session, opts ...rest.RequestOpt) (Session, error) {
	accessToken, err := c.Rest.RefreshAccessToken(c.ID, c.Secret, session.RefreshToken, opts...)
	if err != nil {
		return Session{}, err
	}
	return newSession(*accessToken), nil
}

func (c *Client) VerifySession(session Session, opts ...rest.RequestOpt) (Session, error) {
	if session.Expired() {
		return c.RefreshSession(session, opts...)
	}
	return session, nil
}

func (c *Client) GetUser(session Session, opts ...rest.RequestOpt) (*discord.OAuth2User, error) {
	if err := checkSession(session, discord.OAuth2ScopeIdentify); err != nil {
		return nil, err
	}
	return c.Rest.GetCurrentUser(session.AccessToken, opts...)
}

func (c *Client) GetMember(session Session, guildID snowflake.ID, opts ...rest.RequestOpt) (*discord.Member, error) {
	if err := checkSession(session, discord.OAuth2ScopeGuildsMembersRead); err != nil {
		return nil, err
	}
	return c.Rest.GetCurrentMember(session.AccessToken, guildID, opts...)
}

func (c *Client) GetGuilds(session Session, opts ...rest.RequestOpt) ([]discord.OAuth2Guild, error) {
	if err := checkSession(session, discord.OAuth2ScopeGuilds); err != nil {
		return nil, err
	}
	return c.Rest.GetCurrentUserGuilds(session.AccessToken, 0, 0, 0, false, opts...)
}

func (c *Client) GetConnections(session Session, opts ...rest.RequestOpt) ([]discord.Connection, error) {
	if err := checkSession(session, discord.OAuth2ScopeConnections); err != nil {
		return nil, err
	}
	return c.Rest.GetCurrentUserConnections(session.AccessToken, opts...)
}

func (c *Client) GetApplicationRoleConnection(session Session, applicationID snowflake.ID, opts ...rest.RequestOpt) (*discord.ApplicationRoleConnection, error) {
	if err := checkSession(session, discord.OAuth2ScopeRoleConnectionsWrite); err != nil {
		return nil, err
	}
	return c.Rest.GetCurrentUserApplicationRoleConnection(session.AccessToken, applicationID, opts...)
}

func (c *Client) UpdateApplicationRoleConnection(session Session, applicationID snowflake.ID, update discord.ApplicationRoleConnectionUpdate, opts ...rest.RequestOpt) (*discord.ApplicationRoleConnection, error) {
	if err := checkSession(session, discord.OAuth2ScopeRoleConnectionsWrite); err != nil {
		return nil, err
	}
	return c.Rest.UpdateCurrentUserApplicationRoleConnection(session.AccessToken, applicationID, update, opts...)
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
