package handlers

import (
	"io"

	"github.com/DisgoOrg/disgo/bot"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/disgo/httpserver"
)

func DefaultHTTPServerEventHandler(client bot.Client) httpserver.EventHandlerFunc {
	return func(responseChannel chan<- discord.InteractionResponse, reader io.Reader) {
		client.EventManager().HandleHTTP(responseChannel, events.HandleRawEvent(client, discord.GatewayEventTypeInteractionCreate, -1, reader))
	}
}

func GetHTTPServerHandler() bot.HTTPServerEventHandler {
	return &httpserverHandlerInteractionCreate{}
}

func DefaultGatewayEventHandler(client bot.Client) gateway.EventHandlerFunc {
	return func(gatewayEventType discord.GatewayEventType, sequenceNumber discord.GatewaySequence, reader io.Reader) {
		reader = events.HandleRawEvent(client, gatewayEventType, sequenceNumber, reader)

		client.EventManager().HandleGateway(gatewayEventType, sequenceNumber, reader)
	}
}

func GetGatewayHandlers() map[discord.GatewayEventType]bot.GatewayEventHandler {
	handlers := make(map[discord.GatewayEventType]bot.GatewayEventHandler, len(AllEventHandlers))
	for _, handler := range AllEventHandlers {
		handlers[handler.EventType()] = handler
	}
	return handlers
}

var AllEventHandlers = []bot.GatewayEventHandler{
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
