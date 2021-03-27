package handlers

import (
	"github.com/DiscoOrg/disgo/api"
)

func GetAllHandlers() []api.GatewayEventHandler {
	return []api.GatewayEventHandler{
		ReadyHandler{},

		GuildCreateHandler{},
		GuildUpdateHandler{},
		GuildDeleteHandler{},

		GuildMemberAddHandler{},
		GuildMemberUpdateHandler{},
		GuildMemberRemoveHandler{},

		GuildRoleCreateHandler{},
		GuildRoleUpdateHandler{},
		GuildRoleDeleteHandler{},

		MessageCreateHandler{},

		InteractionCreateHandler{},
	}
}
