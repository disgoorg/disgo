package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/httpserver"
)

// DefaultHTTPServerEventHandlerFunc is the default handler for the httpserver.Server and sends payloads to the bot.EventManager.
func DefaultHTTPServerEventHandlerFunc(client bot.Client) httpserver.EventHandlerFunc {
	return client.EventManager().HandleHTTPEvent
}

// GetHTTPServerHandler returns the default httpserver.Server event handler for processing the raw payload which gets passed into the bot.EventManager
func GetHTTPServerHandler() bot.HTTPServerEventHandler {
	return &httpserverHandlerInteractionCreate{}
}

// DefaultGatewayEventHandlerFunc is the default handler for the gateway.Gateway and sends payloads to the bot.EventManager.
func DefaultGatewayEventHandlerFunc(client bot.Client) gateway.EventHandlerFunc {
	return client.EventManager().HandleGatewayEvent
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
	bot.NewGatewayEventHandler(gateway.EventTypeRaw, gatewayHandlerRaw),
	bot.NewGatewayEventHandler(gateway.EventTypeHeartbeatAck, gatewayHandlerHeartbeatAck),
	bot.NewGatewayEventHandler(gateway.EventTypeReady, gatewayHandlerReady),
	bot.NewGatewayEventHandler(gateway.EventTypeResumed, gatewayHandlerResumed),

	bot.NewGatewayEventHandler(gateway.EventTypeApplicationCommandPermissionsUpdate, gatewayHandlerApplicationCommandPermissionsUpdate),

	bot.NewGatewayEventHandler(gateway.EventTypeAutoModerationRuleCreate, gatewayHandlerAutoModerationRuleCreate),
	bot.NewGatewayEventHandler(gateway.EventTypeAutoModerationRuleUpdate, gatewayHandlerAutoModerationRuleUpdate),
	bot.NewGatewayEventHandler(gateway.EventTypeAutoModerationRuleDelete, gatewayHandlerAutoModerationRuleDelete),
	bot.NewGatewayEventHandler(gateway.EventTypeAutoModerationActionExecution, gatewayHandlerAutoModerationActionExecution),

	bot.NewGatewayEventHandler(gateway.EventTypeChannelCreate, gatewayHandlerChannelCreate),
	bot.NewGatewayEventHandler(gateway.EventTypeChannelUpdate, gatewayHandlerChannelUpdate),
	bot.NewGatewayEventHandler(gateway.EventTypeChannelDelete, gatewayHandlerChannelDelete),
	bot.NewGatewayEventHandler(gateway.EventTypeChannelPinsUpdate, gatewayHandlerChannelPinsUpdate),

	bot.NewGatewayEventHandler(gateway.EventTypeEntitlementCreate, gatewayHandlerEntitlementCreate),
	bot.NewGatewayEventHandler(gateway.EventTypeEntitlementUpdate, gatewayHandlerEntitlementUpdate),
	bot.NewGatewayEventHandler(gateway.EventTypeEntitlementDelete, gatewayHandlerEntitlementDelete),

	bot.NewGatewayEventHandler(gateway.EventTypeThreadCreate, gatewayHandlerThreadCreate),
	bot.NewGatewayEventHandler(gateway.EventTypeThreadUpdate, gatewayHandlerThreadUpdate),
	bot.NewGatewayEventHandler(gateway.EventTypeThreadDelete, gatewayHandlerThreadDelete),
	bot.NewGatewayEventHandler(gateway.EventTypeThreadListSync, gatewayHandlerThreadListSync),
	bot.NewGatewayEventHandler(gateway.EventTypeThreadMemberUpdate, gatewayHandlerThreadMemberUpdate),
	bot.NewGatewayEventHandler(gateway.EventTypeThreadMembersUpdate, gatewayHandlerThreadMembersUpdate),

	bot.NewGatewayEventHandler(gateway.EventTypeGuildCreate, gatewayHandlerGuildCreate),
	bot.NewGatewayEventHandler(gateway.EventTypeGuildUpdate, gatewayHandlerGuildUpdate),
	bot.NewGatewayEventHandler(gateway.EventTypeGuildDelete, gatewayHandlerGuildDelete),

	bot.NewGatewayEventHandler(gateway.EventTypeGuildAuditLogEntryCreate, gatewayHandlerGuildAuditLogEntryCreate),

	bot.NewGatewayEventHandler(gateway.EventTypeGuildBanAdd, gatewayHandlerGuildBanAdd),
	bot.NewGatewayEventHandler(gateway.EventTypeGuildBanRemove, gatewayHandlerGuildBanRemove),

	bot.NewGatewayEventHandler(gateway.EventTypeGuildEmojisUpdate, gatewayHandlerGuildEmojisUpdate),
	bot.NewGatewayEventHandler(gateway.EventTypeGuildStickersUpdate, gatewayHandlerGuildStickersUpdate),
	bot.NewGatewayEventHandler(gateway.EventTypeGuildIntegrationsUpdate, gatewayHandlerGuildIntegrationsUpdate),

	bot.NewGatewayEventHandler(gateway.EventTypeGuildMemberAdd, gatewayHandlerGuildMemberAdd),
	bot.NewGatewayEventHandler(gateway.EventTypeGuildMemberRemove, gatewayHandlerGuildMemberRemove),
	bot.NewGatewayEventHandler(gateway.EventTypeGuildMemberUpdate, gatewayHandlerGuildMemberUpdate),
	bot.NewGatewayEventHandler(gateway.EventTypeGuildMembersChunk, gatewayHandlerGuildMembersChunk),

	bot.NewGatewayEventHandler(gateway.EventTypeGuildRoleCreate, gatewayHandlerGuildRoleCreate),
	bot.NewGatewayEventHandler(gateway.EventTypeGuildRoleUpdate, gatewayHandlerGuildRoleUpdate),
	bot.NewGatewayEventHandler(gateway.EventTypeGuildRoleDelete, gatewayHandlerGuildRoleDelete),

	bot.NewGatewayEventHandler(gateway.EventTypeGuildScheduledEventCreate, gatewayHandlerGuildScheduledEventCreate),
	bot.NewGatewayEventHandler(gateway.EventTypeGuildScheduledEventUpdate, gatewayHandlerGuildScheduledEventUpdate),
	bot.NewGatewayEventHandler(gateway.EventTypeGuildScheduledEventDelete, gatewayHandlerGuildScheduledEventDelete),
	bot.NewGatewayEventHandler(gateway.EventTypeGuildScheduledEventUserAdd, gatewayHandlerGuildScheduledEventUserAdd),
	bot.NewGatewayEventHandler(gateway.EventTypeGuildScheduledEventUserRemove, gatewayHandlerGuildScheduledEventUserRemove),

	bot.NewGatewayEventHandler(gateway.EventTypeIntegrationCreate, gatewayHandlerIntegrationCreate),
	bot.NewGatewayEventHandler(gateway.EventTypeIntegrationUpdate, gatewayHandlerIntegrationUpdate),
	bot.NewGatewayEventHandler(gateway.EventTypeIntegrationDelete, gatewayHandlerIntegrationDelete),

	bot.NewGatewayEventHandler(gateway.EventTypeInteractionCreate, gatewayHandlerInteractionCreate),

	bot.NewGatewayEventHandler(gateway.EventTypeInviteCreate, gatewayHandlerInviteCreate),
	bot.NewGatewayEventHandler(gateway.EventTypeInviteDelete, gatewayHandlerInviteDelete),

	bot.NewGatewayEventHandler(gateway.EventTypeMessageCreate, gatewayHandlerMessageCreate),
	bot.NewGatewayEventHandler(gateway.EventTypeMessageUpdate, gatewayHandlerMessageUpdate),
	bot.NewGatewayEventHandler(gateway.EventTypeMessageDelete, gatewayHandlerMessageDelete),
	bot.NewGatewayEventHandler(gateway.EventTypeMessageDeleteBulk, gatewayHandlerMessageDeleteBulk),

	bot.NewGatewayEventHandler(gateway.EventTypeMessagePollVoteAdd, gatewayHandlerMessagePollVoteAdd),
	bot.NewGatewayEventHandler(gateway.EventTypeMessagePollVoteRemove, gatewayHandlerMessagePollVoteRemove),

	bot.NewGatewayEventHandler(gateway.EventTypeMessageReactionAdd, gatewayHandlerMessageReactionAdd),
	bot.NewGatewayEventHandler(gateway.EventTypeMessageReactionRemove, gatewayHandlerMessageReactionRemove),
	bot.NewGatewayEventHandler(gateway.EventTypeMessageReactionRemoveAll, gatewayHandlerMessageReactionRemoveAll),
	bot.NewGatewayEventHandler(gateway.EventTypeMessageReactionRemoveEmoji, gatewayHandlerMessageReactionRemoveEmoji),

	bot.NewGatewayEventHandler(gateway.EventTypePresenceUpdate, gatewayHandlerPresenceUpdate),

	bot.NewGatewayEventHandler(gateway.EventTypeStageInstanceCreate, gatewayHandlerStageInstanceCreate),
	bot.NewGatewayEventHandler(gateway.EventTypeStageInstanceUpdate, gatewayHandlerStageInstanceUpdate),
	bot.NewGatewayEventHandler(gateway.EventTypeStageInstanceDelete, gatewayHandlerStageInstanceDelete),

	bot.NewGatewayEventHandler(gateway.EventTypeTypingStart, gatewayHandlerTypingStart),
	bot.NewGatewayEventHandler(gateway.EventTypeUserUpdate, gatewayHandlerUserUpdate),

	bot.NewGatewayEventHandler(gateway.EventTypeVoiceStateUpdate, gatewayHandlerVoiceStateUpdate),
	bot.NewGatewayEventHandler(gateway.EventTypeVoiceServerUpdate, gatewayHandlerVoiceServerUpdate),

	bot.NewGatewayEventHandler(gateway.EventTypeWebhooksUpdate, gatewayHandlerWebhooksUpdate),
}
