package handlers

import (
	"io"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/httpserver"
)

func DefaultHTTPServerEventHandler(client bot.Client) httpserver.EventHandlerFunc {
	return func(responseFunc httpserver.RespondFunc, reader io.Reader) {
		client.EventManager().HandleHTTPEvent(responseFunc, events.HandleRawEvent(client, discord.GatewayEventTypeInteractionCreate, -1, responseFunc, reader))
	}
}

func GetHTTPServerHandler() bot.HTTPServerEventHandler {
	return &httpserverHandlerInteractionCreate{}
}

func DefaultGatewayEventHandler(client bot.Client) gateway.EventHandlerFunc {
	return func(gatewayEventType discord.GatewayEventType, sequenceNumber int, reader io.Reader) {
		client.EventManager().HandleGatewayEvent(gatewayEventType, sequenceNumber, events.HandleRawEvent(client, gatewayEventType, sequenceNumber, nil, reader))
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
