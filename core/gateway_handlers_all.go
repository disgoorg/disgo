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
	&gatewayHandlerReady{},
	&gatewayHandlerResumed{},
	&gatewayHandlerInvalidSession{},

	&gatewayHandlerChannelCreat{},
	&gatewayHandlerChannelUpdate{},
	&gatewayHandlerChannelDelete{},
	&gatewayHandlerChannelPinsUpdate{},

	// &gatewayHandlerThreadCreate{},
	// &gatewayHandlerThreadUpdate{},
	// &gatewayHandlerThreadDelete{},
	// &gatewayHandlerThreadListSync{},
	// &gatewayHandlerThreadMemberUpdate{},
	// &gatewayHandlerThreadMembersUpdate{},

	&gatewayHandlerGuildCreate{},
	&gatewayHandlerGuildUpdate{},
	&gatewayHandlerGuildDelete{},

	&gatewayHandlerGuildBanAdd{},
	&gatewayHandlerGuildBanRemove{},

	&gatewayHandlerGuildEmojisUpdate{},
	&gatewayHandlerGuildStickersUpdate{},
	&gatewayHandlerGuildIntegrationsUpdate{},

	&gatewayHandlerGuildMemberAdd{},
	&gatewayHandlerGuildMemberRemove{},
	&gatewayHandlerGuildMemberUpdate{},
	// &gatewayHandlerGuildMemberChunk{},

	&gatewayHandlerGuildRoleCreate{},
	&gatewayHandlerGuildRoleUpdate{},
	&gatewayHandlerGuildRoleDelete{},

	// &gatewayHandlerIntegrationCreate{},
	// &gatewayHandlerIntegrationUpdate{},
	// &gatewayHandlerIntegrationDelete{},

	&gatewayHandlerInteractionCreate{},

	&gatewayHandlerInviteCreate{},
	&gatewayHandlerInviteDelete{},

	&gatewayHandlerMessageCreate{},
	&gatewayHandlerMessageUpdate{},
	&gatewayHandlerMessageDelete{},
	// &gatewayHandlerMessageDeleteBulk{},

	&gatewayHandlerMessageReactionAdd{},
	&gatewayHandlerMessageReactionRemove{},
	&gatewayHandlerMessageReactionRemoveAll{},
	&gatewayHandlerMessageReactionRemoveEmoji{},

	// &gatewayHandlerPresenceUpdate{},

	&gatewayHandlerStageInstanceCreate{},
	&gatewayHandlerStageInstanceDelete{},
	&gatewayHandlerStageInstanceUpdate{},

	&gatewayHandlerTypingStart{},
	// &gatewayHandlerUserUpdate{},

	&gatewayHandlerVoiceStateUpdate{},
	&gatewayHandlerVoiceServerUpdate{},

	&gatewayHandlerWebhooksUpdate{},
}
