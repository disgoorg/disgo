package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

func GetGatewayHandlers() map[discord.GatewayEventType]GatewayEventHandler {
	handlers := make(map[discord.GatewayEventType]GatewayEventHandler, len(EventHandlers))
	for _, handler := range EventHandlers {
		handlers[handler.EventType()] = handler
	}
	return handlers
}

var EventHandlers = []GatewayEventHandler{
	&ChannelCreateHandler{},
	&ChannelDeleteHandler{},
	&ChannelUpdateHandler{},

	&GuildCreateHandler{},
	&GuildDeleteHandler{},
	&GuildUpdateHandler{},

	&GuildMemberAddHandler{},
	&GuildMemberRemoveHandler{},
	&GuildMemberUpdateHandler{},

	&GuildBanAddHandler{},
	&GuildBanRemoveHandler{},

	&GuildRoleCreateHandler{},
	&GuildRoleDeleteHandler{},
	&GuildRoleUpdateHandler{},

	&WebhooksUpdateHandler{},

	&InteractionCreateGatewayHandler{},

	&MessageCreateHandler{},
	&MessageUpdateHandler{},
	&MessageDeleteHandler{},

	&ReadyHandler{},

	&VoiceServerUpdateHandler{},
	&VoiceStateUpdateHandler{},
}
