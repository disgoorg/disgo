package rest

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest/route"
)

var _ GatewayService = (*gatewayServiceImpl)(nil)

func NewGatewayService(restClient Client) GatewayService {
	return &gatewayServiceImpl{restClient: restClient}
}

type GatewayService interface {
	Service
	GetGateway(opts ...RequestOpt) (*discord.Gateway, error)
	GetGatewayBot(opts ...RequestOpt) (*discord.GatewayBot, error)
}

type gatewayServiceImpl struct {
	restClient Client
}

func (s *gatewayServiceImpl) RestClient() Client {
	return s.restClient
}

func (s *gatewayServiceImpl) GetGateway(opts ...RequestOpt) (gateway *discord.Gateway, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetGateway.Compile(nil)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &gateway, opts...)
	return
}

func (s *gatewayServiceImpl) GetGatewayBot(opts ...RequestOpt) (gatewayBot *discord.GatewayBot, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetGatewayBot.Compile(nil)
	if err != nil {
		return
	}
	err = s.restClient.Do(compiledRoute, nil, &gatewayBot, opts...)
	return
}
