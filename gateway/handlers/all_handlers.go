package handlers

import "github.com/DisgoOrg/disgo/core"

func init() {
	println("registering handlers...")
	for _, handler := range EventHandlers {
		core.GatewayEventHandlers[handler.EventType()] = handler
	}
}

var EventHandlers = []core.GatewayEventHandler{
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

	&InteractionCreateHandler{},

	&MessageCreateHandler{},
	&MessageUpdateHandler{},
	&MessageDeleteHandler{},

	&ReadyHandler{},

	&VoiceServerUpdateHandler{},
	&VoiceStateUpdateHandler{},
}
