package rest

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest/route"
)

var _ GatewayService = (*GatewayServiceImpl)(nil)

func NewGatewayService(restClient Client) GatewayService {
	return &GatewayServiceImpl{restClient: restClient}
}

type GatewayService interface {
	Service
	GetGateway(opts ...RequestOpt) (*discord.Gateway, Error)
	GetGatewayBot(opts ...RequestOpt) (*discord.GatewayBot, Error)
}

type GatewayServiceImpl struct {
	restClient Client
}

func (s *GatewayServiceImpl) RestClient() Client {
	return s.restClient
}

func (s *GatewayServiceImpl) GetGateway(opts ...RequestOpt) (gateway *discord.Gateway, rErr Error) {
	compiledRoute, err := route.GetGateway.Compile(nil)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &gateway, opts...)
	return
}

func (s *GatewayServiceImpl) GetGatewayBot(opts ...RequestOpt) (gatewayBot *discord.GatewayBot, rErr Error) {
	compiledRoute, err := route.GetGatewayBot.Compile(nil)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &gatewayBot, opts...)
	return
}
