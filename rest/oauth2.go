package rest

import (
	"net/url"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest/route"
	"github.com/disgoorg/snowflake/v2"
)

var _ OAuth2 = (*oAuth2Impl)(nil)

func NewOAuth2(client Client) OAuth2 {
	return &oAuth2Impl{client: client}
}

type OAuth2 interface {
	GetBotApplicationInfo(opts ...RequestOpt) (*discord.Application, error)

	GetCurrentAuthorizationInfo(bearerToken string, opts ...RequestOpt) (*discord.AuthorizationInformation, error)
	GetCurrentUser(bearerToken string, opts ...RequestOpt) (*discord.OAuth2User, error)
	GetCurrentMember(bearerToken string, guildID snowflake.ID, opts ...RequestOpt) (*discord.Member, error)
	GetCurrentUserGuilds(bearerToken string, before snowflake.ID, after snowflake.ID, limit int, opts ...RequestOpt) ([]discord.OAuth2Guild, error)
	GetCurrentUserConnections(bearerToken string, opts ...RequestOpt) ([]discord.Connection, error)

	SetGuildCommandPermissions(bearerToken string, applicationID snowflake.ID, guildID snowflake.ID, commandID snowflake.ID, commandPermissions []discord.ApplicationCommandPermission, opts ...RequestOpt) (*discord.ApplicationCommandPermissions, error)

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
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetBotApplicationInfo.Compile(nil)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, nil, &application, opts...)
	return
}

func (s *oAuth2Impl) GetCurrentAuthorizationInfo(bearerToken string, opts ...RequestOpt) (info *discord.AuthorizationInformation, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetAuthorizationInfo.Compile(nil)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, nil, &info, withBearerToken(bearerToken, opts)...)
	return
}

func (s *oAuth2Impl) GetCurrentUser(bearerToken string, opts ...RequestOpt) (user *discord.OAuth2User, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetCurrentUser.Compile(nil)
	if err != nil {
		return
	}

	err = s.client.Do(compiledRoute, nil, &user, withBearerToken(bearerToken, opts)...)
	return
}

func (s *oAuth2Impl) GetCurrentMember(bearerToken string, guildID snowflake.ID, opts ...RequestOpt) (member *discord.Member, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetCurrentMember.Compile(nil, guildID)
	if err != nil {
		return
	}

	err = s.client.Do(compiledRoute, nil, &member, withBearerToken(bearerToken, opts)...)
	return
}

func (s *oAuth2Impl) GetCurrentUserGuilds(bearerToken string, before snowflake.ID, after snowflake.ID, limit int, opts ...RequestOpt) (guilds []discord.OAuth2Guild, err error) {
	queryParams := route.QueryValues{}
	if before != 0 {
		queryParams["before"] = before
	}
	if after != 0 {
		queryParams["after"] = after
	}
	if limit != 0 {
		queryParams["limit"] = limit
	}

	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetCurrentUserGuilds.Compile(queryParams)
	if err != nil {
		return
	}

	err = s.client.Do(compiledRoute, nil, &guilds, withBearerToken(bearerToken, opts)...)
	return
}

func (s *oAuth2Impl) GetCurrentUserConnections(bearerToken string, opts ...RequestOpt) (connections []discord.Connection, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetCurrentUserConnections.Compile(nil)
	if err != nil {
		return
	}

	err = s.client.Do(compiledRoute, nil, &connections, withBearerToken(bearerToken, opts)...)
	return
}

func (s *oAuth2Impl) SetGuildCommandPermissions(bearerToken string, applicationID snowflake.ID, guildID snowflake.ID, commandID snowflake.ID, commandPermissions []discord.ApplicationCommandPermission, opts ...RequestOpt) (commandPerms *discord.ApplicationCommandPermissions, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.SetGuildCommandPermissions.Compile(nil, applicationID, guildID, commandID)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, discord.ApplicationCommandPermissionsSet{Permissions: commandPermissions}, &commandPerms, withBearerToken(bearerToken, opts)...)
	return
}

func (s *oAuth2Impl) exchangeAccessToken(clientID snowflake.ID, clientSecret string, grantType discord.GrantType, codeOrRefreshToken string, redirectURI string, opts ...RequestOpt) (exchange *discord.AccessTokenResponse, err error) {
	values := url.Values{
		"client_id":  []string{clientID.String()},
		"grant_type": []string{grantType.String()},
	}
	if clientSecret != "" {
		values["client_secret"] = []string{clientSecret}
	}
	switch grantType {
	case discord.GrantTypeAuthorizationCode:
		values["code"] = []string{codeOrRefreshToken}
		values["redirect_uri"] = []string{redirectURI}

	case discord.GrantTypeRefreshToken:
		values["refresh_token"] = []string{codeOrRefreshToken}
	}
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.Token.Compile(nil)
	if err != nil {
		return
	}

	err = s.client.Do(compiledRoute, values, &exchange, opts...)
	return
}

func (s *oAuth2Impl) GetAccessToken(clientID snowflake.ID, clientSecret string, code string, redirectURI string, opts ...RequestOpt) (exchange *discord.AccessTokenResponse, err error) {
	return s.exchangeAccessToken(clientID, clientSecret, discord.GrantTypeAuthorizationCode, code, redirectURI, opts...)
}

func (s *oAuth2Impl) RefreshAccessToken(clientID snowflake.ID, clientSecret string, refreshToken string, opts ...RequestOpt) (exchange *discord.AccessTokenResponse, err error) {
	return s.exchangeAccessToken(clientID, clientSecret, discord.GrantTypeRefreshToken, refreshToken, "", opts...)
}
