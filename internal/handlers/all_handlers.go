package handlers

import (
	"github.com/DisgoOrg/disgo/api"
)

// GetAllHandlers returns all api.GatewayEventHandler(s)
func GetAllHandlers() []api.EventHandler {
	return []api.EventHandler{
		ApplicationCommandCreateHandler{},
		ApplicationCommandDeleteHandler{},
		ApplicationCommandUpdateHandler{},

		ChannelCreateHandler{},
		ChannelDeleteHandler{},
		ChannelUpdateHandler{},

		ThreadCreateHandler{},
		ThreadDeleteHandler{},
		ThreadUpdateHandler{},
		ThreadListSyncHandler{},

		GuildCreateHandler{},
		GuildDeleteHandler{},
		GuildUpdateHandler{},

		GuildMemberAddHandler{},
		GuildMemberRemoveHandler{},
		GuildMemberUpdateHandler{},

		GuildRoleCreateHandler{},
		GuildRoleDeleteHandler{},
		GuildRoleUpdateHandler{},

		InteractionCreateHandler{},
		InteractionCreateWebhookHandler{},

		MessageCreateHandler{},

		ReadyHandler{},

		VoiceServerUpdateHandler{},
		VoiceStateUpdateHandler{},
	}
}
