package rest

import (
	"net/url"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest/route"
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
	GetCurrentUserGuilds(bearerToken string, before discord.Snowflake, after discord.Snowflake, limit int, opts ...RequestOpt) ([]discord.OAuth2Guild, error)
	GetCurrentUser(bearerToken string, opts ...RequestOpt) (*discord.OAuth2User, error)
	GetCurrentUserConnections(bearerToken string, opts ...RequestOpt) ([]discord.Connection, error)

	GetAccessToken(clientID discord.Snowflake, clientSecret string, code string, redirectURI string, opts ...RequestOpt) (*discord.AccessTokenExchange, error)
	RefreshAccessToken(clientID discord.Snowflake, clientSecret string, refreshToken string, opts ...RequestOpt) (*discord.AccessTokenExchange, error)
}

type oAuth2ServiceImpl struct {
	restClient Client
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
	err = s.restClient.DoBot(compiledRoute, nil, &application, opts...)
	return
}

func (s *oAuth2ServiceImpl) GetCurrentAuthorizationInfo(bearerToken string, opts ...RequestOpt) (info *discord.AuthorizationInformation, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetAuthorizationInfo.Compile(nil)
	if err != nil {
		return
	}
	err = s.restClient.DoBearer(compiledRoute, nil, &info, bearerToken, opts...)
	return
}

func (s *oAuth2ServiceImpl) GetCurrentUserGuilds(bearerToken string, before discord.Snowflake, after discord.Snowflake, limit int, opts ...RequestOpt) (guilds []discord.OAuth2Guild, err error) {
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

	// this route also supports bot token. Use it when no bearer token is provided.
	if bearerToken == "" {
		opts = applyBotToken(s.RestClient(), opts)
	} else {
		opts = applyBearerToken(bearerToken, opts)
	}

	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetCurrentUserGuilds.Compile(queryParams)
	if err != nil {
		return
	}

	err = s.restClient.Do(compiledRoute, nil, &guilds, opts...)
	return
}

func (s *oAuth2ServiceImpl) GetCurrentUser(bearerToken string, opts ...RequestOpt) (user *discord.OAuth2User, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetCurrentUser.Compile(nil)
	if err != nil {
		return
	}

	err = s.restClient.DoBearer(compiledRoute, nil, &user, bearerToken, opts...)
	return
}

func (s *oAuth2ServiceImpl) GetCurrentUserConnections(bearerToken string, opts ...RequestOpt) (connections []discord.Connection, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetCurrentUserConnections.Compile(nil)
	if err != nil {
		return
	}

	err = s.restClient.DoBearer(compiledRoute, nil, &connections, bearerToken, opts...)
	return
}

func (s *oAuth2ServiceImpl) exchangeAccessToken(clientID discord.Snowflake, clientSecret string, grantType discord.GrantType, codeOrRefreshToken string, redirectURI string, opts ...RequestOpt) (exchange *discord.AccessTokenExchange, err error) {
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

func (s *oAuth2ServiceImpl) GetAccessToken(clientID discord.Snowflake, clientSecret string, code string, redirectURI string, opts ...RequestOpt) (exchange *discord.AccessTokenExchange, err error) {
	return s.exchangeAccessToken(clientID, clientSecret, discord.GrantTypeAuthorizationCode, code, redirectURI, opts...)
}

func (s *oAuth2ServiceImpl) RefreshAccessToken(clientID discord.Snowflake, clientSecret string, refreshToken string, opts ...RequestOpt) (exchange *discord.AccessTokenExchange, err error) {
	return s.exchangeAccessToken(clientID, clientSecret, discord.GrantTypeRefreshToken, refreshToken, "", opts...)
}
