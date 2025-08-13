package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/gateway"
)

type gatewayEventHandler interface {
	bot.GatewayEventHandler
	EventType() gateway.EventType
}

// GetGatewayHandler returns the default gateway.Gateway event handlers for processing the raw payload which gets passed into the bot.EventManager
func GetGatewayHandler() bot.GatewayEventHandler {
	handlers := make(map[gateway.EventType]gatewayEventHandler, len(allEventHandlers))
	for _, handler := range allEventHandlers {
		handlers[handler.EventType()] = handler
	}
	return &gatewayHandler{
		handlers: handlers,
	}
}

type gatewayHandler struct {
	handlers map[gateway.EventType]gatewayEventHandler
}

func (g *gatewayHandler) HandleGatewayEvent(client *bot.Client, message gateway.Message, shardID int) {
	if handler, ok := g.handlers[message.T]; ok {
		handler.HandleGatewayEvent(client, message, shardID)
	} else {
		client.Logger.Warn("No handler found for gateway event", "event_type", message.T, "shard_id", shardID)
	}
}

func newGatewayEventHandler[T gateway.EventData](eventType gateway.EventType, handleFunc func(client *bot.Client, sequenceNumber int, shardID int, event T)) *genericGatewayEventHandler[T] {
	return &genericGatewayEventHandler[T]{
		eventType:  eventType,
		handleFunc: handleFunc,
	}
}

type genericGatewayEventHandler[T gateway.EventData] struct {
	eventType  gateway.EventType
	handleFunc func(client *bot.Client, sequenceNumber int, shardID int, event T)
}

func (h *genericGatewayEventHandler[T]) EventType() gateway.EventType {
	return h.eventType
}

func (h *genericGatewayEventHandler[T]) HandleGatewayEvent(client *bot.Client, message gateway.Message, shardID int) {
	if data, ok := message.D.(T); ok {
		h.handleFunc(client, message.S, shardID, data)
	}
}

var allEventHandlers = []gatewayEventHandler{
	newGatewayEventHandler(gateway.EventTypeRaw, gatewayHandlerRaw),
	newGatewayEventHandler(gateway.EventTypeHeartbeatAck, gatewayHandlerHeartbeatAck),
	newGatewayEventHandler(gateway.EventTypeReady, gatewayHandlerReady),
	newGatewayEventHandler(gateway.EventTypeResumed, gatewayHandlerResumed),

	newGatewayEventHandler(gateway.EventTypeApplicationCommandPermissionsUpdate, gatewayHandlerApplicationCommandPermissionsUpdate),

	newGatewayEventHandler(gateway.EventTypeAutoModerationRuleCreate, gatewayHandlerAutoModerationRuleCreate),
	newGatewayEventHandler(gateway.EventTypeAutoModerationRuleUpdate, gatewayHandlerAutoModerationRuleUpdate),
	newGatewayEventHandler(gateway.EventTypeAutoModerationRuleDelete, gatewayHandlerAutoModerationRuleDelete),
	newGatewayEventHandler(gateway.EventTypeAutoModerationActionExecution, gatewayHandlerAutoModerationActionExecution),

	newGatewayEventHandler(gateway.EventTypeChannelCreate, gatewayHandlerChannelCreate),
	newGatewayEventHandler(gateway.EventTypeChannelUpdate, gatewayHandlerChannelUpdate),
	newGatewayEventHandler(gateway.EventTypeChannelDelete, gatewayHandlerChannelDelete),
	newGatewayEventHandler(gateway.EventTypeChannelPinsUpdate, gatewayHandlerChannelPinsUpdate),

	newGatewayEventHandler(gateway.EventTypeEntitlementCreate, gatewayHandlerEntitlementCreate),
	newGatewayEventHandler(gateway.EventTypeEntitlementUpdate, gatewayHandlerEntitlementUpdate),
	newGatewayEventHandler(gateway.EventTypeEntitlementDelete, gatewayHandlerEntitlementDelete),

	newGatewayEventHandler(gateway.EventTypeThreadCreate, gatewayHandlerThreadCreate),
	newGatewayEventHandler(gateway.EventTypeThreadUpdate, gatewayHandlerThreadUpdate),
	newGatewayEventHandler(gateway.EventTypeThreadDelete, gatewayHandlerThreadDelete),
	newGatewayEventHandler(gateway.EventTypeThreadListSync, gatewayHandlerThreadListSync),
	newGatewayEventHandler(gateway.EventTypeThreadMemberUpdate, gatewayHandlerThreadMemberUpdate),
	newGatewayEventHandler(gateway.EventTypeThreadMembersUpdate, gatewayHandlerThreadMembersUpdate),

	newGatewayEventHandler(gateway.EventTypeGuildCreate, gatewayHandlerGuildCreate),
	newGatewayEventHandler(gateway.EventTypeGuildUpdate, gatewayHandlerGuildUpdate),
	newGatewayEventHandler(gateway.EventTypeGuildDelete, gatewayHandlerGuildDelete),

	newGatewayEventHandler(gateway.EventTypeGuildAuditLogEntryCreate, gatewayHandlerGuildAuditLogEntryCreate),

	newGatewayEventHandler(gateway.EventTypeGuildBanAdd, gatewayHandlerGuildBanAdd),
	newGatewayEventHandler(gateway.EventTypeGuildBanRemove, gatewayHandlerGuildBanRemove),

	newGatewayEventHandler(gateway.EventTypeGuildEmojisUpdate, gatewayHandlerGuildEmojisUpdate),
	newGatewayEventHandler(gateway.EventTypeGuildStickersUpdate, gatewayHandlerGuildStickersUpdate),
	newGatewayEventHandler(gateway.EventTypeGuildIntegrationsUpdate, gatewayHandlerGuildIntegrationsUpdate),

	newGatewayEventHandler(gateway.EventTypeGuildMemberAdd, gatewayHandlerGuildMemberAdd),
	newGatewayEventHandler(gateway.EventTypeGuildMemberRemove, gatewayHandlerGuildMemberRemove),
	newGatewayEventHandler(gateway.EventTypeGuildMemberUpdate, gatewayHandlerGuildMemberUpdate),
	newGatewayEventHandler(gateway.EventTypeGuildMembersChunk, gatewayHandlerGuildMembersChunk),

	newGatewayEventHandler(gateway.EventTypeGuildRoleCreate, gatewayHandlerGuildRoleCreate),
	newGatewayEventHandler(gateway.EventTypeGuildRoleUpdate, gatewayHandlerGuildRoleUpdate),
	newGatewayEventHandler(gateway.EventTypeGuildRoleDelete, gatewayHandlerGuildRoleDelete),

	newGatewayEventHandler(gateway.EventTypeGuildScheduledEventCreate, gatewayHandlerGuildScheduledEventCreate),
	newGatewayEventHandler(gateway.EventTypeGuildScheduledEventUpdate, gatewayHandlerGuildScheduledEventUpdate),
	newGatewayEventHandler(gateway.EventTypeGuildScheduledEventDelete, gatewayHandlerGuildScheduledEventDelete),
	newGatewayEventHandler(gateway.EventTypeGuildScheduledEventUserAdd, gatewayHandlerGuildScheduledEventUserAdd),
	newGatewayEventHandler(gateway.EventTypeGuildScheduledEventUserRemove, gatewayHandlerGuildScheduledEventUserRemove),

	newGatewayEventHandler(gateway.EventTypeGuildSoundboardSoundCreate, gatewayHandlerGuildSoundboardSoundCreate),
	newGatewayEventHandler(gateway.EventTypeGuildSoundboardSoundUpdate, gatewayHandlerGuildSoundboardSoundUpdate),
	newGatewayEventHandler(gateway.EventTypeGuildSoundboardSoundDelete, gatewayHandlerGuildSoundboardSoundDelete),
	newGatewayEventHandler(gateway.EventTypeGuildSoundboardSoundsUpdate, gatewayHandlerGuildSoundboardSoundsUpdate),
	newGatewayEventHandler(gateway.EventTypeSoundboardSounds, gatewayHandlerSoundboardSounds),

	newGatewayEventHandler(gateway.EventTypeIntegrationCreate, gatewayHandlerIntegrationCreate),
	newGatewayEventHandler(gateway.EventTypeIntegrationUpdate, gatewayHandlerIntegrationUpdate),
	newGatewayEventHandler(gateway.EventTypeIntegrationDelete, gatewayHandlerIntegrationDelete),

	newGatewayEventHandler(gateway.EventTypeInteractionCreate, gatewayHandlerInteractionCreate),

	newGatewayEventHandler(gateway.EventTypeInviteCreate, gatewayHandlerInviteCreate),
	newGatewayEventHandler(gateway.EventTypeInviteDelete, gatewayHandlerInviteDelete),

	newGatewayEventHandler(gateway.EventTypeMessageCreate, gatewayHandlerMessageCreate),
	newGatewayEventHandler(gateway.EventTypeMessageUpdate, gatewayHandlerMessageUpdate),
	newGatewayEventHandler(gateway.EventTypeMessageDelete, gatewayHandlerMessageDelete),
	newGatewayEventHandler(gateway.EventTypeMessageDeleteBulk, gatewayHandlerMessageDeleteBulk),

	newGatewayEventHandler(gateway.EventTypeMessagePollVoteAdd, gatewayHandlerMessagePollVoteAdd),
	newGatewayEventHandler(gateway.EventTypeMessagePollVoteRemove, gatewayHandlerMessagePollVoteRemove),

	newGatewayEventHandler(gateway.EventTypeMessageReactionAdd, gatewayHandlerMessageReactionAdd),
	newGatewayEventHandler(gateway.EventTypeMessageReactionRemove, gatewayHandlerMessageReactionRemove),
	newGatewayEventHandler(gateway.EventTypeMessageReactionRemoveAll, gatewayHandlerMessageReactionRemoveAll),
	newGatewayEventHandler(gateway.EventTypeMessageReactionRemoveEmoji, gatewayHandlerMessageReactionRemoveEmoji),

	newGatewayEventHandler(gateway.EventTypePresenceUpdate, gatewayHandlerPresenceUpdate),

	newGatewayEventHandler(gateway.EventTypeStageInstanceCreate, gatewayHandlerStageInstanceCreate),
	newGatewayEventHandler(gateway.EventTypeStageInstanceUpdate, gatewayHandlerStageInstanceUpdate),
	newGatewayEventHandler(gateway.EventTypeStageInstanceDelete, gatewayHandlerStageInstanceDelete),

	newGatewayEventHandler(gateway.EventTypeSubscriptionCreate, gatewayHandlerSubscriptionCreate),
	newGatewayEventHandler(gateway.EventTypeSubscriptionUpdate, gatewayHandlerSubscriptionUpdate),
	newGatewayEventHandler(gateway.EventTypeSubscriptionDelete, gatewayHandlerSubscriptionDelete),

	newGatewayEventHandler(gateway.EventTypeTypingStart, gatewayHandlerTypingStart),
	newGatewayEventHandler(gateway.EventTypeUserUpdate, gatewayHandlerUserUpdate),

	newGatewayEventHandler(gateway.EventTypeVoiceChannelEffectSend, gatewayHandlerVoiceChannelEffectSend),
	newGatewayEventHandler(gateway.EventTypeVoiceStateUpdate, gatewayHandlerVoiceStateUpdate),
	newGatewayEventHandler(gateway.EventTypeVoiceServerUpdate, gatewayHandlerVoiceServerUpdate),

	newGatewayEventHandler(gateway.EventTypeWebhooksUpdate, gatewayHandlerWebhooksUpdate),
}
