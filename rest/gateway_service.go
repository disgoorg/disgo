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
	GetGateway(opts ...RequestOpt) (*discord.Gateway, Error)
	GetGatewayBot(opts ...RequestOpt) (*discord.GatewayBot, Error)
}

type gatewayServiceImpl struct {
	restClient Client
}

func (s *gatewayServiceImpl) RestClient() Client {
	return s.restClient
}

func (s *gatewayServiceImpl) GetGateway(opts ...RequestOpt) (gateway *discord.Gateway, rErr Error) {
	compiledRoute, err := route.GetGateway.Compile(nil)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &gateway, opts...)
	return
}

func (s *gatewayServiceImpl) GetGatewayBot(opts ...RequestOpt) (gatewayBot *discord.GatewayBot, rErr Error) {
	compiledRoute, err := route.GetGatewayBot.Compile(nil)
	if err != nil {
		return nil, NewError(nil, err)
	}
	rErr = s.restClient.Do(compiledRoute, nil, &gatewayBot, opts...)
	return
}
