package handlers

import (
	"io"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/httpserver"
)

// DefaultHTTPServerEventHandler is the default handler for the httpserver.Server and sends payloads to the bot.EventManager.
func DefaultHTTPServerEventHandler(client bot.Client) httpserver.EventHandlerFunc {
	return func(responseFunc httpserver.RespondFunc, reader io.Reader) {
		client.EventManager().HandleHTTPEvent(responseFunc, events.HandleRawEvent(client, gateway.EventTypeInteractionCreate, -1, -1, responseFunc, reader))
	}
}

// GetHTTPServerHandler returns the default httpserver.Server event handler for processing the raw payload which gets passed into the bot.EventManager
func GetHTTPServerHandler() bot.HTTPServerEventHandler {
	return &httpserverHandlerInteractionCreate{}
}

// DefaultGatewayEventHandler is the default handler for the gateway.Gateway and sends payloads to the bot.EventManager.
func DefaultGatewayEventHandler(client bot.Client) gateway.EventHandlerFunc {
	return func(gatewayEventType gateway.EventType, sequenceNumber int, shardID int, reader io.Reader) {
		client.EventManager().HandleGatewayEvent(gatewayEventType, sequenceNumber, shardID, events.HandleRawEvent(client, gatewayEventType, sequenceNumber, shardID, nil, reader))
	}
}

// GetGatewayHandlers returns the default gateway.Gateway event handlers for processing the raw payload which gets passed into the bot.EventManager
func GetGatewayHandlers() map[gateway.EventType]bot.GatewayEventHandler {
	handlers := make(map[gateway.EventType]bot.GatewayEventHandler, len(allEventHandlers))
	for _, handler := range allEventHandlers {
		handlers[handler.EventType()] = handler
	}
	return handlers
}

var allEventHandlers = []bot.GatewayEventHandler{
	&gatewayHandlerReady{},
	&gatewayHandlerResumed{},

	&gatewayHandlerApplicationCommandPermissionsUpdate{},

	&gatewayHandlerAutoModerationRuleCreate{},
	&gatewayHandlerAutoModerationRuleUpdate{},
	&gatewayHandlerAutoModerationRuleDelete{},
	&gatewayHandlerAutoModerationActionExecution{},

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
