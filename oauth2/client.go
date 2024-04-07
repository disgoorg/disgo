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

// Client is a high level wrapper around Discord's OAuth2 API.
type Client interface {
	// ID returns the configured client ID.
	ID() snowflake.ID
	// Secret returns the configured client secret.
	Secret() string
	// Rest returns the underlying rest.OAuth2.
	Rest() rest.OAuth2

	// StateController returns the configured StateController.
	StateController() StateController

	// GenerateAuthorizationURL generates an authorization URL with the given authorization params. State is automatically generated.
	GenerateAuthorizationURL(params AuthorizationURLParams) string
	// GenerateAuthorizationURLState generates an authorization URL with the given authorization params. State is automatically generated & returned.
	GenerateAuthorizationURLState(params AuthorizationURLParams) (string, string)

	// StartSession starts a new Session with the given authorization code & state.
	StartSession(code string, state string, opts ...rest.RequestOpt) (Session, *discord.IncomingWebhook, error)
	// RefreshSession refreshes the given Session with the refresh token.
	RefreshSession(session Session, opts ...rest.RequestOpt) (Session, error)
	// VerifySession verifies the given Session & refreshes it if needed.
	VerifySession(session Session, opts ...rest.RequestOpt) (Session, error)

	// GetUser returns the discord.OAuth2User associated with the given Session. Fields filled in the struct depend on the Session.Scopes.
	GetUser(session Session, opts ...rest.RequestOpt) (*discord.OAuth2User, error)
	// GetMember returns the discord.Member associated with the given Session in a specific guild.
	GetMember(session Session, guildID snowflake.ID, opts ...rest.RequestOpt) (*discord.Member, error)
	// GetGuilds returns the discord.OAuth2Guild(s) the user is a member of. This requires the discord.OAuth2ScopeGuilds scope in the Session.
	GetGuilds(session Session, opts ...rest.RequestOpt) ([]discord.OAuth2Guild, error)
	// GetConnections returns the discord.Connection(s) the user has connected. This requires the discord.OAuth2ScopeConnections scope in the Session.
	GetConnections(session Session, opts ...rest.RequestOpt) ([]discord.Connection, error)
	// GetApplicationRoleConnection returns the discord.ApplicationRoleConnection for the given application. This requires the discord.OAuth2ScopeRoleConnectionsWrite scope in the Session.
	GetApplicationRoleConnection(session Session, applicationID snowflake.ID, opts ...rest.RequestOpt) (*discord.ApplicationRoleConnection, error)
	// UpdateApplicationRoleConnection updates the discord.ApplicationRoleConnection for the given application. This requires the discord.OAuth2ScopeRoleConnectionsWrite scope in the Session.
	UpdateApplicationRoleConnection(session Session, applicationID snowflake.ID, update discord.ApplicationRoleConnectionUpdate, opts ...rest.RequestOpt) (*discord.ApplicationRoleConnection, error)
}
