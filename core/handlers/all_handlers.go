package handlers

import (
	"io"

	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/disgo/httpserver"
)

func DefaultHTTPServerEventHandler(bot core.Bot) httpserver.EventHandlerFunc {
	return func(responseChannel chan<- discord.InteractionResponse, reader io.Reader) {
		bot.EventManager().HandleHTTP(responseChannel, events.HandleRawEvent(bot, discord.GatewayEventTypeInteractionCreate, -1, reader))
	}
}

func GetHTTPServerHandler() core.HTTPServerEventHandler {
	return &httpserverHandlerInteractionCreate{}
}

func DefaultGatewayEventHandler(bot core.Bot) gateway.EventHandlerFunc {
	return func(gatewayEventType discord.GatewayEventType, sequenceNumber discord.GatewaySequence, reader io.Reader) {
		reader = events.HandleRawEvent(bot, gatewayEventType, sequenceNumber, reader)

		bot.EventManager().HandleGateway(gatewayEventType, sequenceNumber, reader)
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

	&gatewayHandlerGuildScheduledEventCreate{},
	&gatewayHandlerGuildScheduledEventUpdate{},
	&gatewayHandlerGuildScheduledEventDelete{},
	&gatewayHandlerGuildScheduledEventUserAdd{},
	&gatewayHandlerGuildScheduledEventUserRemove{},

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
