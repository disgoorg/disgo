package rest

import (
	"github.com/DisgoOrg/disgo/discord"
)

func NewGatewayService(client Client) GatewayService {
	return nil
}

type GatewayService interface {
	Service
	GetGateway(opts ...RequestOpt) (*discord.Gateway, Error)
	GetGatewayBot(opts ...RequestOpt) (*discord.GatewayBot, Error)
}
