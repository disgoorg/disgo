package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

// gatewayHandlerVoiceStateUpdate handles core.GatewayEventVoiceStateUpdate
type gatewayHandlerVoiceStateUpdate struct{}

// EventType returns the gateway.EventType
func (h *gatewayHandlerVoiceStateUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeVoiceStateUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerVoiceStateUpdate) New() interface{} {
	return &discord.VoiceState{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerVoiceStateUpdate) HandleGatewayEvent(bot core.Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*discord.VoiceState)

	oldVoiceState := bot.Caches().VoiceStates().GetCopy(payload.GuildID, payload.UserID)

	voiceState := bot.EntityBuilder().CreateVoiceState(payload, core.CacheStrategyYes)

	if oldVoiceState != nil && oldVoiceState.ChannelID != nil {
		if channel := bot.Caches().Channels().Get(*oldVoiceState.ChannelID); channel != nil {
			if ch, ok := channel.(*core.GuildVoiceChannel); ok {
				delete(ch.ConnectedMemberIDs, voiceState.UserID)
			} else if ch, ok := channel.(*core.GuildStageVoiceChannel); ok {
				delete(ch.ConnectedMemberIDs, voiceState.UserID)
			}
		}
	}

	if voiceState.ChannelID != nil {
		if channel := bot.Caches().Channels().Get(*voiceState.ChannelID); channel != nil {
			if ch, ok := channel.(*core.GuildVoiceChannel); ok {
				ch.ConnectedMemberIDs[voiceState.UserID] = struct{}{}
			} else if ch, ok := channel.(*core.GuildStageVoiceChannel); ok {
				ch.ConnectedMemberIDs[voiceState.UserID] = struct{}{}
			}

		}
	}

	genericGuildVoiceEvent := &events.GenericGuildVoiceEvent{
		GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
		VoiceState:   voiceState,
	}

	bot.EventManager().Dispatch(&events.GuildVoiceStateUpdateEvent{
		GenericGuildVoiceEvent: genericGuildVoiceEvent,
		OldVoiceState:          oldVoiceState,
	})

	if oldVoiceState != nil && oldVoiceState.ChannelID != nil && payload.ChannelID != nil {
		bot.EventManager().Dispatch(&events.GuildVoiceMoveEvent{
			GenericGuildVoiceEvent: genericGuildVoiceEvent,
			OldVoiceState:          oldVoiceState,
		})
	} else if (oldVoiceState == nil || oldVoiceState.ChannelID == nil) && payload.ChannelID != nil {
		bot.EventManager().Dispatch(&events.GuildVoiceJoinEvent{
			GenericGuildVoiceEvent: genericGuildVoiceEvent,
		})
	} else if payload.ChannelID == nil {
		bot.EventManager().Dispatch(&events.GuildVoiceLeaveEvent{
			GenericGuildVoiceEvent: genericGuildVoiceEvent,
			OldVoiceState:          oldVoiceState,
		})
	} else {
		bot.Logger.Warnf("could not decide which GuildVoiceEvent to fire")
	}
}
