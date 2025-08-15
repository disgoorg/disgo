package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/gateway"
)

// GetGatewayHandler returns the default gateway.Gateway event handlers for processing the raw payload which gets passed into the bot.EventManager
func GetGatewayHandler() bot.GatewayEventHandler {
	return &gatewayHandler{}
}

type gatewayHandler struct{}

func (g *gatewayHandler) HandleGatewayEvent(client *bot.Client, shardID int, message gateway.Message) {
	switch event := message.D.(type) {
	case gateway.EventRaw:
		gatewayHandlerRaw(client, message.S, shardID, event)
	case gateway.EventHeartbeatAck:
		gatewayHandlerHeartbeatAck(client, message.S, shardID, event)
	case gateway.EventReady:
		gatewayHandlerReady(client, message.S, shardID, event)
	case gateway.EventResumed:
		gatewayHandlerResumed(client, message.S, shardID, event)
	case gateway.EventApplicationCommandPermissionsUpdate:
		gatewayHandlerApplicationCommandPermissionsUpdate(client, message.S, shardID, event)
	case gateway.EventAutoModerationRuleCreate:
		gatewayHandlerAutoModerationRuleCreate(client, message.S, shardID, event)
	case gateway.EventAutoModerationRuleUpdate:
		gatewayHandlerAutoModerationRuleUpdate(client, message.S, shardID, event)
	case gateway.EventAutoModerationRuleDelete:
		gatewayHandlerAutoModerationRuleDelete(client, message.S, shardID, event)
	case gateway.EventAutoModerationActionExecution:
		gatewayHandlerAutoModerationActionExecution(client, message.S, shardID, event)
	case gateway.EventChannelCreate:
		gatewayHandlerChannelCreate(client, message.S, shardID, event)
	case gateway.EventChannelUpdate:
		gatewayHandlerChannelUpdate(client, message.S, shardID, event)
	case gateway.EventChannelDelete:
		gatewayHandlerChannelDelete(client, message.S, shardID, event)
	case gateway.EventChannelPinsUpdate:
		gatewayHandlerChannelPinsUpdate(client, message.S, shardID, event)
	case gateway.EventEntitlementCreate:
		gatewayHandlerEntitlementCreate(client, message.S, shardID, event)
	case gateway.EventEntitlementUpdate:
		gatewayHandlerEntitlementUpdate(client, message.S, shardID, event)
	case gateway.EventEntitlementDelete:
		gatewayHandlerEntitlementDelete(client, message.S, shardID, event)
	case gateway.EventThreadCreate:
		gatewayHandlerThreadCreate(client, message.S, shardID, event)
	case gateway.EventThreadUpdate:
		gatewayHandlerThreadUpdate(client, message.S, shardID, event)
	case gateway.EventThreadDelete:
		gatewayHandlerThreadDelete(client, message.S, shardID, event)
	case gateway.EventThreadListSync:
		gatewayHandlerThreadListSync(client, message.S, shardID, event)
	case gateway.EventThreadMemberUpdate:
		gatewayHandlerThreadMemberUpdate(client, message.S, shardID, event)
	case gateway.EventThreadMembersUpdate:
		gatewayHandlerThreadMembersUpdate(client, message.S, shardID, event)
	case gateway.EventGuildCreate:
		gatewayHandlerGuildCreate(client, message.S, shardID, event)
	case gateway.EventGuildUpdate:
		gatewayHandlerGuildUpdate(client, message.S, shardID, event)
	case gateway.EventGuildDelete:
		gatewayHandlerGuildDelete(client, message.S, shardID, event)
	case gateway.EventGuildAuditLogEntryCreate:
		gatewayHandlerGuildAuditLogEntryCreate(client, message.S, shardID, event)
	case gateway.EventGuildBanAdd:
		gatewayHandlerGuildBanAdd(client, message.S, shardID, event)
	case gateway.EventGuildBanRemove:
		gatewayHandlerGuildBanRemove(client, message.S, shardID, event)
	case gateway.EventGuildEmojisUpdate:
		gatewayHandlerGuildEmojisUpdate(client, message.S, shardID, event)
	case gateway.EventGuildStickersUpdate:
		gatewayHandlerGuildStickersUpdate(client, message.S, shardID, event)
	case gateway.EventGuildIntegrationsUpdate:
		gatewayHandlerGuildIntegrationsUpdate(client, message.S, shardID, event)
	case gateway.EventGuildMemberAdd:
		gatewayHandlerGuildMemberAdd(client, message.S, shardID, event)
	case gateway.EventGuildMemberRemove:
		gatewayHandlerGuildMemberRemove(client, message.S, shardID, event)
	case gateway.EventGuildMemberUpdate:
		gatewayHandlerGuildMemberUpdate(client, message.S, shardID, event)
	case gateway.EventGuildMembersChunk:
		gatewayHandlerGuildMembersChunk(client, message.S, shardID, event)
	case gateway.EventGuildRoleCreate:
		gatewayHandlerGuildRoleCreate(client, message.S, shardID, event)
	case gateway.EventGuildRoleUpdate:
		gatewayHandlerGuildRoleUpdate(client, message.S, shardID, event)
	case gateway.EventGuildRoleDelete:
		gatewayHandlerGuildRoleDelete(client, message.S, shardID, event)
	case gateway.EventGuildScheduledEventCreate:
		gatewayHandlerGuildScheduledEventCreate(client, message.S, shardID, event)
	case gateway.EventGuildScheduledEventUpdate:
		gatewayHandlerGuildScheduledEventUpdate(client, message.S, shardID, event)
	case gateway.EventGuildScheduledEventDelete:
		gatewayHandlerGuildScheduledEventDelete(client, message.S, shardID, event)
	case gateway.EventGuildScheduledEventUserAdd:
		gatewayHandlerGuildScheduledEventUserAdd(client, message.S, shardID, event)
	case gateway.EventGuildScheduledEventUserRemove:
		gatewayHandlerGuildScheduledEventUserRemove(client, message.S, shardID, event)
	case gateway.EventGuildSoundboardSoundCreate:
		gatewayHandlerGuildSoundboardSoundCreate(client, message.S, shardID, event)
	case gateway.EventGuildSoundboardSoundUpdate:
		gatewayHandlerGuildSoundboardSoundUpdate(client, message.S, shardID, event)
	case gateway.EventGuildSoundboardSoundDelete:
		gatewayHandlerGuildSoundboardSoundDelete(client, message.S, shardID, event)
	case gateway.EventGuildSoundboardSoundsUpdate:
		gatewayHandlerGuildSoundboardSoundsUpdate(client, message.S, shardID, event)
	case gateway.EventSoundboardSounds:
		gatewayHandlerSoundboardSounds(client, message.S, shardID, event)
	case gateway.EventIntegrationCreate:
		gatewayHandlerIntegrationCreate(client, message.S, shardID, event)
	case gateway.EventIntegrationUpdate:
		gatewayHandlerIntegrationUpdate(client, message.S, shardID, event)
	case gateway.EventIntegrationDelete:
		gatewayHandlerIntegrationDelete(client, message.S, shardID, event)
	case gateway.EventInteractionCreate:
		gatewayHandlerInteractionCreate(client, message.S, shardID, event)
	case gateway.EventInviteCreate:
		gatewayHandlerInviteCreate(client, message.S, shardID, event)
	case gateway.EventInviteDelete:
		gatewayHandlerInviteDelete(client, message.S, shardID, event)
	case gateway.EventMessageCreate:
		gatewayHandlerMessageCreate(client, message.S, shardID, event)
	case gateway.EventMessageUpdate:
		gatewayHandlerMessageUpdate(client, message.S, shardID, event)
	case gateway.EventMessageDelete:
		gatewayHandlerMessageDelete(client, message.S, shardID, event)
	case gateway.EventMessageDeleteBulk:
		gatewayHandlerMessageDeleteBulk(client, message.S, shardID, event)
	case gateway.EventMessagePollVoteAdd:
		gatewayHandlerMessagePollVoteAdd(client, message.S, shardID, event)
	case gateway.EventMessagePollVoteRemove:
		gatewayHandlerMessagePollVoteRemove(client, message.S, shardID, event)
	case gateway.EventMessageReactionAdd:
		gatewayHandlerMessageReactionAdd(client, message.S, shardID, event)
	case gateway.EventMessageReactionRemove:
		gatewayHandlerMessageReactionRemove(client, message.S, shardID, event)
	case gateway.EventMessageReactionRemoveAll:
		gatewayHandlerMessageReactionRemoveAll(client, message.S, shardID, event)
	case gateway.EventMessageReactionRemoveEmoji:
		gatewayHandlerMessageReactionRemoveEmoji(client, message.S, shardID, event)
	case gateway.EventPresenceUpdate:
		gatewayHandlerPresenceUpdate(client, message.S, shardID, event)
	case gateway.EventStageInstanceCreate:
		gatewayHandlerStageInstanceCreate(client, message.S, shardID, event)
	case gateway.EventStageInstanceUpdate:
		gatewayHandlerStageInstanceUpdate(client, message.S, shardID, event)
	case gateway.EventStageInstanceDelete:
		gatewayHandlerStageInstanceDelete(client, message.S, shardID, event)
	case gateway.EventSubscriptionCreate:
		gatewayHandlerSubscriptionCreate(client, message.S, shardID, event)
	case gateway.EventSubscriptionUpdate:
		gatewayHandlerSubscriptionUpdate(client, message.S, shardID, event)
	case gateway.EventSubscriptionDelete:
		gatewayHandlerSubscriptionDelete(client, message.S, shardID, event)
	case gateway.EventTypingStart:
		gatewayHandlerTypingStart(client, message.S, shardID, event)
	case gateway.EventUserUpdate:
		gatewayHandlerUserUpdate(client, message.S, shardID, event)
	case gateway.EventVoiceChannelEffectSend:
		gatewayHandlerVoiceChannelEffectSend(client, message.S, shardID, event)
	case gateway.EventVoiceStateUpdate:
		gatewayHandlerVoiceStateUpdate(client, message.S, shardID, event)
	case gateway.EventVoiceServerUpdate:
		gatewayHandlerVoiceServerUpdate(client, message.S, shardID, event)
	case gateway.EventWebhooksUpdate:
		gatewayHandlerWebhooksUpdate(client, message.S, shardID, event)
	case gateway.EventUnknown:
		// nothing to do
	default:
		client.Logger.Warn("Unknown gateway event type", "event_type", message.T, "shard_id", shardID)
	}
}
