package internal

import (
	"github.com/DiscoOrg/disgo/api"
	"github.com/DiscoOrg/disgo/api/constants"
	"github.com/DiscoOrg/disgo/internal/handlers"
)

func GetHandlers() *map[string]api.GatewayEventProvider {
	return &map[string]api.GatewayEventProvider{
		constants.GuildCreateEvent:   handlers.GuildCreateHandler{},
		constants.GuildDeleteEvent:   handlers.GuildDeleteHandler{},
		constants.GuildUpdateEvent:   handlers.GuildUpdateHandler{},
		constants.MessageCreateEvent: handlers.MessageCreateHandler{},
	}
}
