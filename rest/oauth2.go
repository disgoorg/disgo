package rest

import (
	"net/url"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest/route"
	"github.com/disgoorg/snowflake"
)

var _ OAuth2 = (*oAuth2Impl)(nil)

func NewOAuth2(client Client) OAuth2 {
	return &oAuth2Impl{client: client}
}

type OAuth2 interface {
	GetBotApplicationInfo(opts ...RequestOpt) (*discord.Application, error)

	GetCurrentAuthorizationInfo(bearerToken string, opts ...RequestOpt) (*discord.AuthorizationInformation, error)
	GetCurrentUser(bearerToken string, opts ...RequestOpt) (*discord.OAuth2User, error)
	GetCurrentUserGuilds(bearerToken string, before snowflake.Snowflake, after snowflake.Snowflake, limit int, opts ...RequestOpt) ([]discord.OAuth2Guild, error)
	GetCurrentUserConnections(bearerToken string, opts ...RequestOpt) ([]discord.Connection, error)

	GetAccessToken(clientID snowflake.Snowflake, clientSecret string, code string, redirectURI string, opts ...RequestOpt) (*discord.AccessTokenExchange, error)
	RefreshAccessToken(clientID snowflake.Snowflake, clientSecret string, refreshToken string, opts ...RequestOpt) (*discord.AccessTokenExchange, error)
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

func (s *oAuth2Impl) GetCurrentUserGuilds(bearerToken string, before snowflake.Snowflake, after snowflake.Snowflake, limit int, opts ...RequestOpt) (guilds []discord.OAuth2Guild, err error) {
	queryParams := route.QueryValues{}
	if before != "" {
		queryParams["before"] = before
	}
	if after != "" {
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

func (s *oAuth2Impl) exchangeAccessToken(clientID snowflake.Snowflake, clientSecret string, grantType discord.GrantType, codeOrRefreshToken string, redirectURI string, opts ...RequestOpt) (exchange *discord.AccessTokenExchange, err error) {
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
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.Token.Compile(nil)
	if err != nil {
		return
	}

	err = s.client.Do(compiledRoute, values, &exchange, opts...)
	return
}

func (s *oAuth2Impl) GetAccessToken(clientID snowflake.Snowflake, clientSecret string, code string, redirectURI string, opts ...RequestOpt) (exchange *discord.AccessTokenExchange, err error) {
	return s.exchangeAccessToken(clientID, clientSecret, discord.GrantTypeAuthorizationCode, code, redirectURI, opts...)
}

func (s *oAuth2Impl) RefreshAccessToken(clientID snowflake.Snowflake, clientSecret string, refreshToken string, opts ...RequestOpt) (exchange *discord.AccessTokenExchange, err error) {
	return s.exchangeAccessToken(clientID, clientSecret, discord.GrantTypeRefreshToken, refreshToken, "", opts...)
}
