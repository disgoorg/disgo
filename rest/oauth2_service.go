package rest

import (
	"net/url"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest/route"
)

var _ OAuth2Service = (*oAuth2ServiceImpl)(nil)

func NewOAuth2Service(restClient Client) OAuth2Service {
	return &oAuth2ServiceImpl{restClient: restClient}
}

type OAuth2Service interface {
	Service
	GetBotApplicationInfo(opts ...RequestOpt) (*discord.Application, error)
	GetAuthorizationInfo(opts ...RequestOpt) (*discord.AuthorizationInformation, error)

	GetCurrentUserGuilds(token string, opts ...RequestOpt) ([]discord.OAuth2Guild, error)
	GetCurrentUser(token string, opts ...RequestOpt) (*discord.OAuth2User, error)
	GetCurrentUserConnections(token string, opts ...RequestOpt) ([]discord.Connection, error)

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
	err = s.restClient.Do(compiledRoute, nil, &application, opts...)
	return
}

func (s *oAuth2ServiceImpl) GetAuthorizationInfo(opts ...RequestOpt) (info *discord.AuthorizationInformation, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetAuthorizationInfo.Compile(nil)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &info, opts...)
	return
}

func (s *oAuth2ServiceImpl) GetCurrentUserGuilds(token string, opts ...RequestOpt) (guilds []discord.OAuth2Guild, err error) {
	queryParams := route.QueryValues{}

	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetCurrentUserGuilds.Compile(queryParams)
	if err != nil {
		return
	}

	err = s.restClient.Do(compiledRoute, nil, &guilds, append(opts, WithHeader("authorization", discord.TokenTypeBearer.Apply(token)))...)
	return
}

func (s *oAuth2ServiceImpl) GetCurrentUser(token string, opts ...RequestOpt) (user *discord.OAuth2User, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetCurrentUser.Compile(nil)
	if err != nil {
		return
	}

	err = s.restClient.Do(compiledRoute, nil, &user, append(opts, WithHeader("authorization", discord.TokenTypeBearer.Apply(token)))...)
	return
}

func (s *oAuth2ServiceImpl) GetCurrentUserConnections(token string, opts ...RequestOpt) (connections []discord.Connection, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetCurrentUserConnections.Compile(nil)
	if err != nil {
		return
	}

	err = s.restClient.Do(compiledRoute, nil, &connections, append(opts, WithHeader("authorization", discord.TokenTypeBearer.Apply(token)))...)
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
