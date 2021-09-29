package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

func GetGatewayHandlers() map[discord.GatewayEventType]GatewayEventHandler {
	handlers := make(map[discord.GatewayEventType]GatewayEventHandler, len(AllEventHandlers))
	for _, handler := range AllEventHandlers {
		handlers[handler.EventType()] = handler
	}
	return handlers
}

var AllEventHandlers = []GatewayEventHandler{
	&gatewayHandlerReady{},
	&gatewayHandlerResumed{},
	&gatewayHandlerInvalidSession{},

	&gatewayHandlerChannelCreate{},
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
	&gatewayHandlerGuildMembersChunk{},

	&gatewayHandlerGuildRoleCreate{},
	&gatewayHandlerGuildRoleUpdate{},
	&gatewayHandlerGuildRoleDelete{},

	&gatewayHandlerIntegrationCreate{},
	&gatewayHandlerIntegrationUpdate{},
	&gatewayHandlerIntegrationDelete{},

	&gatewayHandlerInteractionCreate{},

	&gatewayHandlerInviteCreate{},
	&gatewayHandlerInviteDelete{},

	&gatewayHandlerMessageCreate{},
	&gatewayHandlerMessageUpdate{},
	&gatewayHandlerMessageDelete{},
	&gatewayHandlerMessageDeleteBulk{},

	&gatewayHandlerMessageReactionAdd{},
	&gatewayHandlerMessageReactionRemove{},
	&gatewayHandlerMessageReactionRemoveAll{},
	&gatewayHandlerMessageReactionRemoveEmoji{},

	&gatewayHandlerPresenceUpdate{},

	&gatewayHandlerStageInstanceCreate{},
	&gatewayHandlerStageInstanceDelete{},
	&gatewayHandlerStageInstanceUpdate{},

	&gatewayHandlerTypingStart{},
	&gatewayHandlerUserUpdate{},

	&gatewayHandlerVoiceStateUpdate{},
	&gatewayHandlerVoiceServerUpdate{},

	&gatewayHandlerWebhooksUpdate{},
}
