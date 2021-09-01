package rest

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest/route"
)

func NewOAuth2Service(client Client) OAuth2Service {
	return &OAuth2ServiceImpl{restClient: client}
}

type OAuth2Service interface {
	Service
	GetBotApplicationInfo(opts ...RequestOpt) (*discord.Application, Error)
	GetAuthorizationInfo(opts ...RequestOpt) (*discord.AuthorizationInformation, Error)

	GetCurrentUserGuilds(token string, opts ...RequestOpt) ([]discord.PartialGuild, Error)
	GetCurrentUser(token string, opts ...RequestOpt) (*discord.OAuth2User, Error)

	GetAccessToken(clientID discord.Snowflake, clientSecret string, code string, redirectURI string, opts ...RequestOpt) (*discord.AccessTokenExchange, Error)
	RefreshAccessToken(clientID discord.Snowflake, clientSecret string, refreshToken string, opts ...RequestOpt) (*discord.AccessTokenExchange, Error)
}

type OAuth2ServiceImpl struct {
	restClient Client
}

func (s *OAuth2ServiceImpl) RestClient() Client {
	return s.restClient
}

func (s *OAuth2ServiceImpl) GetBotApplicationInfo(opts ...RequestOpt) (application *discord.Application, rErr Error) {
	compiledRoute, err := route.GetBotApplicationInfo.Compile(nil)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &application, opts...)
	return
}

func (s *OAuth2ServiceImpl) GetAuthorizationInfo(opts ...RequestOpt) (info *discord.AuthorizationInformation, rErr Error) {
	compiledRoute, err := route.GetAuthorizationInfo.Compile(nil)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &info, opts...)
	return
}

func (s *OAuth2ServiceImpl) GetCurrentUserGuilds(token string, opts ...RequestOpt) (guilds []discord.PartialGuild, rErr Error) {
	queryParams := route.QueryValues{}

	compiledRoute, err := route.GetCurrentUserGuilds.Compile(queryParams)
	if err != nil {
		return nil, NewError(nil, NewError(nil, err))
	}

	opts = append(opts, WithHeader("authorization", discord.TokenTypeBearer.Apply(token)))

	rErr = s.restClient.Do(compiledRoute, nil, &guilds, opts...)
	return
}

func (s *OAuth2ServiceImpl) GetCurrentUser(token string, opts ...RequestOpt) (user *discord.OAuth2User, rErr Error) {
	compiledRoute, err := route.GetCurrentUser.Compile(nil)
	if err != nil {
		return nil, NewError(nil, err)
	}

	opts = append(opts, WithHeader("authorization", discord.TokenTypeBearer.Apply(token)))

	rErr = s.restClient.Do(compiledRoute, nil, &user, opts...)
	return
}

func (s *OAuth2ServiceImpl) exchangeAccessToken(clientID discord.Snowflake, clientSecret string, grantType discord.GrantType, codeOrRefreshToken string, redirectURI string, opts ...RequestOpt) (exchange *discord.AccessTokenExchange, rErr Error) {
	values := map[string]interface{}{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"grant_type":    grantType,
	}
	switch grantType {
	case discord.GrantTypeAuthorizationCode:
		values["code"] = codeOrRefreshToken
		values["redirect_uri"] = redirectURI

	case discord.GrantTypeRefreshToken:
		values["refresh_token"] = codeOrRefreshToken
	}
	compiledRoute, err := route.Token.Compile(values)
	if err != nil {
		return nil, NewError(nil, err)
	}

	rErr = s.restClient.Do(compiledRoute, nil, &exchange, opts...)
	return
}

func (s *OAuth2ServiceImpl) GetAccessToken(clientID discord.Snowflake, clientSecret string, code string, redirectURI string, opts ...RequestOpt) (exchange *discord.AccessTokenExchange, rErr Error) {
	return s.exchangeAccessToken(clientID, clientSecret, discord.GrantTypeAuthorizationCode, code, redirectURI, opts...)
}

func (s *OAuth2ServiceImpl) RefreshAccessToken(clientID discord.Snowflake, clientSecret string, refreshToken string, opts ...RequestOpt) (exchange *discord.AccessTokenExchange, rErr Error) {
	return s.exchangeAccessToken(clientID, clientSecret, discord.GrantTypeRefreshToken, refreshToken, "", opts...)
}
