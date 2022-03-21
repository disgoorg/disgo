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
func (h *gatewayHandlerVoiceStateUpdate) New() any {
	return &discord.VoiceState{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerVoiceStateUpdate) HandleGatewayEvent(bot core.Bot, sequenceNumber discord.GatewaySequence, v any) {
	voiceState := *v.(*discord.VoiceState)

	oldVoiceState, oldOk := bot.Caches().VoiceStates().Get(voiceState.GuildID, voiceState.UserID)
	if voiceState.ChannelID == nil {
		bot.Caches().VoiceStates().Remove(voiceState.GuildID, voiceState.UserID)
	} else {
		bot.Caches().VoiceStates().Put(voiceState.GuildID, voiceState.UserID, voiceState)
	}

	genericGuildVoiceEvent := &events.GenericGuildVoiceStateEvent{
		GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
		VoiceState:   voiceState,
	}

	bot.EventManager().Dispatch(&events.GuildVoiceStateUpdateEvent{
		GenericGuildVoiceStateEvent: genericGuildVoiceEvent,
		OldVoiceState:               oldVoiceState,
	})

	if oldOk && oldVoiceState.ChannelID != nil && voiceState.ChannelID != nil {
		bot.EventManager().Dispatch(&events.GuildVoiceMoveEvent{
			GenericGuildVoiceStateEvent: genericGuildVoiceEvent,
			OldVoiceState:               oldVoiceState,
		})
	} else if (oldOk || oldVoiceState.ChannelID == nil) && voiceState.ChannelID != nil {
		bot.EventManager().Dispatch(&events.GuildVoiceJoinEvent{
			GenericGuildVoiceStateEvent: genericGuildVoiceEvent,
		})
	} else if voiceState.ChannelID == nil {
		bot.EventManager().Dispatch(&events.GuildVoiceLeaveEvent{
			GenericGuildVoiceStateEvent: genericGuildVoiceEvent,
			OldVoiceState:               oldVoiceState,
		})
	} else {
		bot.Logger().Warnf("could not decide which GuildVoiceEvent to fire")
	}
}
