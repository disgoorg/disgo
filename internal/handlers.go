package internal

import (
	"github.com/DiscoOrg/disgo/api"
	"github.com/DiscoOrg/disgo/internal/handlers"
)

func getHandlers() *map[string]api.GatewayEventProvider {
	return &map[string]api.GatewayEventProvider{
		api.GuildCreateGatewayEvent:   handlers.GuildCreateHandler{},
		api.GuildDeleteGatewayEvent:   handlers.GuildDeleteHandler{},
		api.GuildUpdateGatewayEvent:   handlers.GuildUpdateHandler{},
		api.MessageCreateGatewayEvent: handlers.MessageCreateHandler{},
	}
}
