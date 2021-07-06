package handlers

import (
	"github.com/DisgoOrg/disgo/api"
)

// GetAllHandlers returns all api.GatewayEventHandler(s)
func GetAllHandlers() []api.EventHandler {
	return []api.EventHandler{
		&CommandCreateHandler{},
		&CommandDeleteHandler{},
		&CommandUpdateHandler{},

		&ChannelCreateHandler{},
		&ChannelDeleteHandler{},
		&ChannelUpdateHandler{},

		&GuildCreateHandler{},
		&GuildDeleteHandler{},
		&GuildUpdateHandler{},

		&GuildMemberAddHandler{},
		&GuildMemberRemoveHandler{},
		&GuildMemberUpdateHandler{},

		&GuildRoleCreateHandler{},
		&GuildRoleDeleteHandler{},
		&GuildRoleUpdateHandler{},

		&WebhooksUpdateHandler{},

		&InteractionCreateHandler{},
		&InteractionCreateWebhookHandler{},

		&MessageCreateHandler{},
		&MessageUpdateHandler{},
		&MessageDeleteHandler{},

		&ReadyHandler{},

		&VoiceServerUpdateHandler{},
		&VoiceStateUpdateHandler{},
	}
}
