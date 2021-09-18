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
	&gatewayHandlerChannelCreat{},
	&gatewayHandlerChannelDelete{},
	&gatewayHandlerChannelUpdate{},

	&gatewayHandlerGuildCreate{},
	&gatewayHandlerGuildDelete{},
	&gatewayHandlerGuildUpdate{},

	&gatewayHandlerGuildMemberAdd{},
	&gatewayHandlerGuildMemberRemove{},
	&gatewayHandlerGuildMemberUpdate{},

	&gatewayHandlerGuildBanAdd{},
	&gatewayHandlerGuildBanRemove{},

	&gatewayHandlerGuildRoleCreate{},
	&gatewayHandlerGuildRoleDelete{},
	&gatewayHandlerGuildRoleUpdate{},

	&gatewayHandlerGuildEmojisUpdate{},
	&gatewayHandlerGuildStickersUpdate{},

	&gatewayHandlerInviteCreate{},
	&gatewayHandlerInviteDelete{},

	&gatewayHandlerStageInstanceCreate{},
	&gatewayHandlerStageInstanceUpdate{},
	&gatewayHandlerStageInstanceDelete{},

	&gatewayHandlerWebhooksUpdate{},

	&gatewayHandlerInteractionCreate{},

	&gatewayHandlerTypingStart{},

	&gatewayHandlerMessageCreate{},
	&gatewayHandlerMessageUpdate{},
	&gatewayHandlerMessageDelete{},

	&gatewayHandlerMessageReactionAdd{},
	&gatewayHandlerMessageReactionRemove{},

	&gatewayHandlerReady{},

	&gatewayHandlerVoiceServerUpdate{},
	&gatewayHandlerVoiceStateUpdate{},
}
