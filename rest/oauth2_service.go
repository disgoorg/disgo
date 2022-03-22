package rest

import (
	"net/url"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest/route"
	"github.com/DisgoOrg/snowflake"
)

var (
	_ Service       = (*oAuth2ServiceImpl)(nil)
	_ OAuth2Service = (*oAuth2ServiceImpl)(nil)
)

func NewOAuth2Service(restClient Client) OAuth2Service {
	return &oAuth2ServiceImpl{restClient: restClient}
}

type OAuth2Service interface {
	Service
	GetBotApplicationInfo(opts ...RequestOpt) (*discord.Application, error)

	GetCurrentAuthorizationInfo(bearerToken string, opts ...RequestOpt) (*discord.AuthorizationInformation, error)
	GetCurrentUser(bearerToken string, opts ...RequestOpt) (*discord.OAuth2User, error)
	GetCurrentUserGuilds(bearerToken string, before snowflake.Snowflake, after snowflake.Snowflake, limit int, opts ...RequestOpt) ([]discord.OAuth2Guild, error)
	GetCurrentUserConnections(bearerToken string, opts ...RequestOpt) ([]discord.Connection, error)

	GetAccessToken(clientID snowflake.Snowflake, clientSecret string, code string, redirectURI string, opts ...RequestOpt) (*discord.AccessTokenExchange, error)
	RefreshAccessToken(clientID snowflake.Snowflake, clientSecret string, refreshToken string, opts ...RequestOpt) (*discord.AccessTokenExchange, error)
}

type oAuth2ServiceImpl struct {
	restClient Client
}

func withBearerToken(bearerToken string, opts []RequestOpt) []RequestOpt {
	if bearerToken != "" {
		return append(opts, WithToken(discord.TokenTypeBearer, bearerToken))
	}
	return opts
}

func (s *oAuth2ServiceImpl) RestClient() Client {
	return s.restClient
}

func (s *oAuth2ServiceImpl) GetBotApplicationInfo(opts ...RequestOpt) (application *discord.Application, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetBotApplicationInfo.Compile(nil)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &application, opts...)
	return
}

func (s *oAuth2ServiceImpl) GetCurrentAuthorizationInfo(bearerToken string, opts ...RequestOpt) (info *discord.AuthorizationInformation, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetAuthorizationInfo.Compile(nil)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &info, withBearerToken(bearerToken, opts)...)
	return
}

func (s *oAuth2ServiceImpl) GetCurrentUser(bearerToken string, opts ...RequestOpt) (user *discord.OAuth2User, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetCurrentUser.Compile(nil)
	if err != nil {
		return
	}

	err = s.restClient.Do(compiledRoute, nil, &user, withBearerToken(bearerToken, opts)...)
	return
}

func (s *oAuth2ServiceImpl) GetCurrentUserGuilds(bearerToken string, before snowflake.Snowflake, after snowflake.Snowflake, limit int, opts ...RequestOpt) (guilds []discord.OAuth2Guild, err error) {
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

	err = s.restClient.Do(compiledRoute, nil, &guilds, withBearerToken(bearerToken, opts)...)
	return
}

func (s *oAuth2ServiceImpl) GetCurrentUserConnections(bearerToken string, opts ...RequestOpt) (connections []discord.Connection, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetCurrentUserConnections.Compile(nil)
	if err != nil {
		return
	}

	err = s.restClient.Do(compiledRoute, nil, &connections, withBearerToken(bearerToken, opts)...)
	return
}

func (s *oAuth2ServiceImpl) exchangeAccessToken(clientID snowflake.Snowflake, clientSecret string, grantType discord.GrantType, codeOrRefreshToken string, redirectURI string, opts ...RequestOpt) (exchange *discord.AccessTokenExchange, err error) {
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

	err = s.restClient.Do(compiledRoute, values, &exchange, opts...)
	return
}

func (s *oAuth2ServiceImpl) GetAccessToken(clientID snowflake.Snowflake, clientSecret string, code string, redirectURI string, opts ...RequestOpt) (exchange *discord.AccessTokenExchange, err error) {
	return s.exchangeAccessToken(clientID, clientSecret, discord.GrantTypeAuthorizationCode, code, redirectURI, opts...)
}

func (s *oAuth2ServiceImpl) RefreshAccessToken(clientID snowflake.Snowflake, clientSecret string, refreshToken string, opts ...RequestOpt) (exchange *discord.AccessTokenExchange, err error) {
	return s.exchangeAccessToken(clientID, clientSecret, discord.GrantTypeRefreshToken, refreshToken, "", opts...)
}
