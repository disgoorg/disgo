package rest

import (
	"errors"
	"net/url"

	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
)

// ErrMissingBearerToken is returned when a bearer token is missing for a request which requires it.
var ErrMissingBearerToken = errors.New("missing bearer token")

var _ OAuth2 = (*oAuth2Impl)(nil)

func NewOAuth2(client Client) OAuth2 {
	return &oAuth2Impl{client: client}
}

type OAuth2 interface {
	GetBotApplicationInfo(opts ...RequestOpt) (*discord.Application, error)

	GetCurrentAuthorizationInfo(bearerToken string, opts ...RequestOpt) (*discord.AuthorizationInformation, error)
	// GetCurrentUser returns the current user
	// Leave bearerToken empty to use the bot token.
	GetCurrentUser(bearerToken string, opts ...RequestOpt) (*discord.OAuth2User, error)
	GetCurrentMember(bearerToken string, guildID snowflake.ID, opts ...RequestOpt) (*discord.Member, error)
	// GetCurrentUserGuilds returns a list of guilds the current user is a member of. Requires the discord.OAuth2ScopeGuilds scope.
	// Leave bearerToken empty to use the bot token.
	GetCurrentUserGuilds(bearerToken string, before snowflake.ID, after snowflake.ID, limit int, withCounts bool, opts ...RequestOpt) ([]discord.OAuth2Guild, error)
	// GetCurrentUserGuildsPage returns a Page of guilds the current user is a member of. Requires the discord.OAuth2ScopeGuilds scope.
	// Leave bearerToken empty to use the bot token.
	GetCurrentUserGuildsPage(bearerToken string, startID snowflake.ID, limit int, withCounts bool, opts ...RequestOpt) Page[discord.OAuth2Guild]
	GetCurrentUserConnections(bearerToken string, opts ...RequestOpt) ([]discord.Connection, error)

	SetGuildCommandPermissions(bearerToken string, applicationID snowflake.ID, guildID snowflake.ID, commandID snowflake.ID, commandPermissions []discord.ApplicationCommandPermission, opts ...RequestOpt) (*discord.ApplicationCommandPermissions, error)

	GetCurrentUserApplicationRoleConnection(bearerToken string, applicationID snowflake.ID, opts ...RequestOpt) (*discord.ApplicationRoleConnection, error)
	UpdateCurrentUserApplicationRoleConnection(bearerToken string, applicationID snowflake.ID, connectionUpdate discord.ApplicationRoleConnectionUpdate, opts ...RequestOpt) (*discord.ApplicationRoleConnection, error)

	GetAccessToken(clientID snowflake.ID, clientSecret string, code string, redirectURI string, opts ...RequestOpt) (*discord.AccessTokenResponse, error)
	RefreshAccessToken(clientID snowflake.ID, clientSecret string, refreshToken string, opts ...RequestOpt) (*discord.AccessTokenResponse, error)
}

type oAuth2Impl struct {
	client Client
}

func withBearerToken(bearerToken string, opts []RequestOpt) []RequestOpt {
	if bearerToken != "" {
		return append([]RequestOpt{WithToken(discord.TokenTypeBearer, bearerToken)}, opts...)
	}
	return opts
}

func (s *oAuth2Impl) GetBotApplicationInfo(opts ...RequestOpt) (application *discord.Application, err error) {
	err = s.client.Do(GetBotApplicationInfo.Compile(nil), nil, &application, opts...)
	return
}

func (s *oAuth2Impl) GetCurrentAuthorizationInfo(bearerToken string, opts ...RequestOpt) (info *discord.AuthorizationInformation, err error) {
	if bearerToken == "" {
		return nil, ErrMissingBearerToken
	}
	err = s.client.Do(GetAuthorizationInfo.Compile(nil), nil, &info, withBearerToken(bearerToken, opts)...)
	return
}

func (s *oAuth2Impl) GetCurrentUser(bearerToken string, opts ...RequestOpt) (user *discord.OAuth2User, err error) {
	err = s.client.Do(GetCurrentUser.Compile(nil), nil, &user, withBearerToken(bearerToken, opts)...)
	return
}

func (s *oAuth2Impl) GetCurrentMember(bearerToken string, guildID snowflake.ID, opts ...RequestOpt) (member *discord.Member, err error) {
	if bearerToken == "" {
		return nil, ErrMissingBearerToken
	}
	err = s.client.Do(GetCurrentMember.Compile(nil, guildID), nil, &member, withBearerToken(bearerToken, opts)...)
	return
}

func (s *oAuth2Impl) GetCurrentUserGuilds(bearerToken string, before snowflake.ID, after snowflake.ID, limit int, withCounts bool, opts ...RequestOpt) (guilds []discord.OAuth2Guild, err error) {
	queryParams := discord.QueryValues{
		"with_counts": withCounts,
	}
	if before != 0 {
		queryParams["before"] = before
	}
	if after != 0 {
		queryParams["after"] = after
	}
	if limit != 0 {
		queryParams["limit"] = limit
	}
	err = s.client.Do(GetCurrentUserGuilds.Compile(queryParams), nil, &guilds, withBearerToken(bearerToken, opts)...)
	return
}

func (s *oAuth2Impl) GetCurrentUserGuildsPage(bearerToken string, startID snowflake.ID, limit int, withCounts bool, opts ...RequestOpt) Page[discord.OAuth2Guild] {
	return Page[discord.OAuth2Guild]{
		getItemsFunc: func(before snowflake.ID, after snowflake.ID) ([]discord.OAuth2Guild, error) {
			return s.GetCurrentUserGuilds(bearerToken, before, after, limit, withCounts, opts...)
		},
		getIDFunc: func(guild discord.OAuth2Guild) snowflake.ID {
			return guild.ID
		},
		ID: startID,
	}
}

func (s *oAuth2Impl) GetCurrentUserConnections(bearerToken string, opts ...RequestOpt) (connections []discord.Connection, err error) {
	if bearerToken == "" {
		return nil, ErrMissingBearerToken
	}
	err = s.client.Do(GetCurrentUserConnections.Compile(nil), nil, &connections, withBearerToken(bearerToken, opts)...)
	return
}

func (s *oAuth2Impl) SetGuildCommandPermissions(bearerToken string, applicationID snowflake.ID, guildID snowflake.ID, commandID snowflake.ID, commandPermissions []discord.ApplicationCommandPermission, opts ...RequestOpt) (commandPerms *discord.ApplicationCommandPermissions, err error) {
	if bearerToken == "" {
		return nil, ErrMissingBearerToken
	}
	err = s.client.Do(SetGuildCommandPermissions.Compile(nil, applicationID, guildID, commandID), discord.ApplicationCommandPermissionsSet{Permissions: commandPermissions}, &commandPerms, withBearerToken(bearerToken, opts)...)
	return
}

func (s *oAuth2Impl) GetCurrentUserApplicationRoleConnection(bearerToken string, applicationID snowflake.ID, opts ...RequestOpt) (connection *discord.ApplicationRoleConnection, err error) {
	if bearerToken == "" {
		return nil, ErrMissingBearerToken
	}
	err = s.client.Do(GetCurrentUserApplicationRoleConnection.Compile(nil, applicationID), nil, &connection, withBearerToken(bearerToken, opts)...)
	return
}

func (s *oAuth2Impl) UpdateCurrentUserApplicationRoleConnection(bearerToken string, applicationID snowflake.ID, connectionUpdate discord.ApplicationRoleConnectionUpdate, opts ...RequestOpt) (connection *discord.ApplicationRoleConnection, err error) {
	if bearerToken == "" {
		return nil, ErrMissingBearerToken
	}
	err = s.client.Do(UpdateCurrentUserApplicationRoleConnection.Compile(nil, applicationID), connectionUpdate, &connection, withBearerToken(bearerToken, opts)...)
	return
}

func (s *oAuth2Impl) exchangeAccessToken(clientID snowflake.ID, clientSecret string, grantType discord.GrantType, codeOrRefreshToken string, redirectURI string, opts ...RequestOpt) (exchange *discord.AccessTokenResponse, err error) {
	values := url.Values{
		"client_id":     []string{clientID.String()},
		"client_secret": []string{clientSecret},
		"grant_type":    []string{grantType.String()},
	}
	switch grantType {
	case discord.GrantTypeAuthorizationCode:
		values["code"] = []string{codeOrRefreshToken}
		values["redirect_uri"] = []string{redirectURI}

	case discord.GrantTypeRefreshToken:
		values["refresh_token"] = []string{codeOrRefreshToken}
	}
	err = s.client.Do(Token.Compile(nil), values, &exchange, opts...)
	return
}

func (s *oAuth2Impl) GetAccessToken(clientID snowflake.ID, clientSecret string, code string, redirectURI string, opts ...RequestOpt) (exchange *discord.AccessTokenResponse, err error) {
	return s.exchangeAccessToken(clientID, clientSecret, discord.GrantTypeAuthorizationCode, code, redirectURI, opts...)
}

func (s *oAuth2Impl) RefreshAccessToken(clientID snowflake.ID, clientSecret string, refreshToken string, opts ...RequestOpt) (exchange *discord.AccessTokenResponse, err error) {
	return s.exchangeAccessToken(clientID, clientSecret, discord.GrantTypeRefreshToken, refreshToken, "", opts...)
}
