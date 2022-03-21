package handlers

import (
	"github.com/DisgoOrg/disgo/bot"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
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
func (h *gatewayHandlerVoiceStateUpdate) HandleGatewayEvent(client bot.Client, sequenceNumber discord.GatewaySequence, v any) {
	voiceState := *v.(*discord.VoiceState)

	oldVoiceState, oldOk := client.Caches().VoiceStates().Get(voiceState.GuildID, voiceState.UserID)
	if voiceState.ChannelID == nil {
		client.Caches().VoiceStates().Remove(voiceState.GuildID, voiceState.UserID)
	} else {
		client.Caches().VoiceStates().Put(voiceState.GuildID, voiceState.UserID, voiceState)
	}

	genericGuildVoiceEvent := &events.GenericGuildVoiceStateEvent{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber),
		VoiceState:   voiceState,
	}

	client.EventManager().Dispatch(&events.GuildVoiceStateUpdateEvent{
		GenericGuildVoiceStateEvent: genericGuildVoiceEvent,
		OldVoiceState:               oldVoiceState,
	})

	if oldOk && oldVoiceState.ChannelID != nil && voiceState.ChannelID != nil {
		client.EventManager().Dispatch(&events.GuildVoiceMoveEvent{
			GenericGuildVoiceStateEvent: genericGuildVoiceEvent,
			OldVoiceState:               oldVoiceState,
		})
	} else if (oldOk || oldVoiceState.ChannelID == nil) && voiceState.ChannelID != nil {
		client.EventManager().Dispatch(&events.GuildVoiceJoinEvent{
			GenericGuildVoiceStateEvent: genericGuildVoiceEvent,
		})
	} else if voiceState.ChannelID == nil {
		client.EventManager().Dispatch(&events.GuildVoiceLeaveEvent{
			GenericGuildVoiceStateEvent: genericGuildVoiceEvent,
			OldVoiceState:               oldVoiceState,
		})
	} else {
		client.Logger().Warnf("could not decide which GuildVoiceEvent to fire")
	}
}
