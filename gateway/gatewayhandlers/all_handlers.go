package gatewayhandlers

import (
	"io"

	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
	"github.com/DisgoOrg/disgo/gateway"
)

func DefaultGatewayEventHandler(bot *core.Bot) gateway.EventHandlerFunc {
	return func(gatewayEventType discord.GatewayEventType, sequenceNumber int, reader io.Reader) {
		reader = events.HandleRawEvent(bot, gatewayEventType, sequenceNumber, reader)

		bot.EventManager.HandleGateway(gatewayEventType, sequenceNumber, reader)
	}
}

func GetGatewayHandlers() map[discord.GatewayEventType]core.GatewayEventHandler {
	handlers := make(map[discord.GatewayEventType]core.GatewayEventHandler, len(AllEventHandlers))
	for _, handler := range AllEventHandlers {
		handlers[handler.EventType()] = handler
	}
	return handlers
}

var AllEventHandlers = []core.GatewayEventHandler{
	&gatewayHandlerReady{},
	&gatewayHandlerResumed{},
	&gatewayHandlerInvalidSession{},

	&gatewayHandlerChannelCreate{},
	&gatewayHandlerChannelUpdate{},
	&gatewayHandlerChannelDelete{},
	&gatewayHandlerChannelPinsUpdate{},

	&gatewayHandlerThreadCreate{},
	&gatewayHandlerThreadUpdate{},
	&gatewayHandlerThreadDelete{},
	&gatewayHandlerThreadListSync{},
	&gatewayHandlerThreadMemberUpdate{},
	&gatewayHandlerThreadMembersUpdate{},

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
