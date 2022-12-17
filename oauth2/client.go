package oauth2

import (
	"errors"
	"fmt"

	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
)

var (
	// ErrStateNotFound is returned when the state is not found in the SessionController.
	ErrStateNotFound = errors.New("state could not be found")

	// ErrAccessTokenExpired is returned when the access token has expired.
	ErrAccessTokenExpired = errors.New("access token expired. refresh the session")

	// ErrMissingOAuth2Scope is returned when a specific OAuth2 scope is missing.
	ErrMissingOAuth2Scope = func(scope discord.OAuth2Scope) error {
		return fmt.Errorf("missing '%s' scope", scope)
	}
)

// Client is a high level wrapper around Discord's OAuth2 API.
type Client interface {
	// ID returns the configured client ID
	ID() snowflake.ID
	// Secret returns the configured client secret
	Secret() string
	// Rest returns the underlying rest.OAuth2
	Rest() rest.OAuth2

	// SessionController returns the configured SessionController
	SessionController() SessionController
	// StateController returns the configured StateController
	StateController() StateController

	// GenerateAuthorizationURL generates an authorization URL with the given redirect URI, permissions, guildID, disableGuildSelect & scopes. State is automatically generated
	GenerateAuthorizationURL(redirectURI string, permissions discord.Permissions, guildID snowflake.ID, disableGuildSelect bool, scopes ...discord.OAuth2Scope) string
	// GenerateAuthorizationURLState generates an authorization URL with the given redirect URI, permissions, guildID, disableGuildSelect & scopes. State is automatically generated & returned
	GenerateAuthorizationURLState(redirectURI string, permissions discord.Permissions, guildID snowflake.ID, disableGuildSelect bool, scopes ...discord.OAuth2Scope) (string, string)

	// StartSession starts a new Session with the given authorization code & state
	StartSession(code string, state string, identifier string, opts ...rest.RequestOpt) (Session, error)
	// RefreshSession refreshes the given Session with the refresh token
	RefreshSession(identifier string, session Session, opts ...rest.RequestOpt) (Session, error)

	// GetUser returns the discord.OAuth2User associated with the given Session. Fields filled in the struct depend on the Session.Scopes
	GetUser(session Session, opts ...rest.RequestOpt) (*discord.OAuth2User, error)
	// GetMember returns the discord.Member associated with the given Session in a specific guild.
	GetMember(session Session, guildID snowflake.ID, opts ...rest.RequestOpt) (*discord.Member, error)
	// GetGuilds returns the discord.OAuth2Guild(s) the user is a member of. This requires the discord.OAuth2ScopeGuilds scope in the Session
	GetGuilds(session Session, opts ...rest.RequestOpt) ([]discord.OAuth2Guild, error)
	// GetConnections returns the discord.Connection(s) the user has connected. This requires the discord.OAuth2ScopeConnections scope in the Session
	GetConnections(session Session, opts ...rest.RequestOpt) ([]discord.Connection, error)
	// GetApplicationRoleConnection returns the discord.ApplicationRoleConnection for the given application. This requires the discord.OAuth2ScopeRoleConnectionsWrite scope in the Session
	GetApplicationRoleConnection(session Session, applicationID snowflake.ID, opts ...rest.RequestOpt) (*discord.ApplicationRoleConnection, error)
	// UpdateApplicationRoleConnection updates the discord.ApplicationRoleConnection for the given application. This requires the discord.OAuth2ScopeRoleConnectionsWrite scope in the Session
	UpdateApplicationRoleConnection(session Session, applicationID snowflake.ID, update discord.ApplicationRoleConnectionUpdate, opts ...rest.RequestOpt) (*discord.ApplicationRoleConnection, error)
}
