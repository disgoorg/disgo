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
	GetGateway(ctx context.Context) (*discord.Gateway, Error)
	GetGatewayBot(ctx context.Context) (*discord.GatewayBot, Error)
}

