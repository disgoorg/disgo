package handlers

import (
	"github.com/DisgoOrg/disgo/api"
)

// GetAllHandlers returns all api.GatewayEventHandler(s)
func GetAllHandlers() []api.EventHandler {
	return []api.EventHandler{
		ReadyHandler{},

		VoiceServerUpdateHandler{},
		VoiceStateUpdateHandler{},

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
		InteractionCreateWebhookHandler{},
	}
}
