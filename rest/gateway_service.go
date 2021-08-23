package rest

import (
	"context"

	"github.com/DisgoOrg/disgo/discord"
)

func NewGatewayService(client Client) GatewayService {
	return nil
}

type GatewayService interface {
	Service
	GetGateway(opts ...rest.RequestOpt) (*discord.Gateway, Error)
	GetGatewayBot(opts ...rest.RequestOpt) (*discord.GatewayBot, Error)
}
