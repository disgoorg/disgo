package internal

import (
	"github.com/DiscoOrg/disgo"
	"github.com/DiscoOrg/disgo/constants"
	"github.com/DiscoOrg/disgo/handlers"
)

func GetHandlers() *map[string]disgo.GatewayEventProvider {
	return &map[string]disgo.GatewayEventProvider{
		constants.GuildCreateEvent:   handlers.GuildCreateHandler{},
		constants.GuildDeleteEvent:   handlers.GuildDeleteHandler{},
		constants.GuildUpdateEvent:   handlers.GuildUpdateHandler{},
		constants.MessageCreateEvent: handlers.MessageCreateHandler{},
	}
}
