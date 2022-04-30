package rest

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest/route"
)

var _ Gateway = (*gatewayImpl)(nil)

func NewGateway(client Client) Gateway {
	return &gatewayImpl{client: client}
}

type Gateway interface {
	GetGateway(opts ...RequestOpt) (*discord.Gateway, error)
	GetGatewayBot(opts ...RequestOpt) (*discord.GatewayBot, error)
}

type gatewayImpl struct {
	client Client
}

func (s *gatewayImpl) GetGateway(opts ...RequestOpt) (gateway *discord.Gateway, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetGateway.Compile(nil)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, nil, &gateway, opts...)
	return
}

func (s *gatewayImpl) GetGatewayBot(opts ...RequestOpt) (gatewayBot *discord.GatewayBot, err error) {
	var compiledRoute *route.CompiledAPIRoute
	compiledRoute, err = route.GetGatewayBot.Compile(nil)
	if err != nil {
		return
	}
	err = s.client.Do(compiledRoute, nil, &gatewayBot, opts...)
	return
}
