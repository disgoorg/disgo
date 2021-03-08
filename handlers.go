package disgo

import (
	"github.com/DiscoOrg/disgo/constants"
)

func GetHandlers() map[string]GatewayEventProvider {
	return map[string]GatewayEventProvider{
		constants.GuildCreateEvent:   GuildCreateHandler{},
		constants.GuildDeleteEvent:   GuildDeleteHandler{},
		constants.GuildUpdateEvent:   GuildUpdateHandler{},
		constants.MessageCreateEvent: MessageCreateHandler{},
	}
}