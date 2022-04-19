package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

// gatewayHandlerVoiceStateUpdate handles core.GatewayEventVoiceStateUpdate
type gatewayHandlerVoiceStateUpdate struct{}

// EventType returns the gateway.EventType
func (h *gatewayHandlerVoiceStateUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeVoiceStateUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerVoiceStateUpdate) New() any {
	return &discord.FullVoiceState{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerVoiceStateUpdate) HandleGatewayEvent(client bot.Client, sequenceNumber int, v any) {
	voiceState := *v.(*discord.FullVoiceState)
	member := voiceState.Member
	// populate unset fields
	member.GuildID = voiceState.GuildID

	oldVoiceState, oldOk := client.Caches().VoiceStates().Get(voiceState.GuildID, voiceState.UserID)
	if voiceState.ChannelID == nil {
		client.Caches().VoiceStates().Remove(voiceState.GuildID, voiceState.UserID)
	} else {
		client.Caches().VoiceStates().Put(voiceState.GuildID, voiceState.UserID, voiceState.VoiceState)
	}
	client.Caches().Members().Put(voiceState.GuildID, voiceState.UserID, member)

	genericGuildVoiceEvent := &events.GenericGuildVoiceStateEvent{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber),
		VoiceState:   voiceState.VoiceState,
		Member:       member,
	}

	client.EventManager().DispatchEvent(&events.GuildVoiceStateUpdateEvent{
		GenericGuildVoiceStateEvent: genericGuildVoiceEvent,
		OldVoiceState:               oldVoiceState,
	})

	if oldOk && oldVoiceState.ChannelID != nil && voiceState.ChannelID != nil {
		client.EventManager().DispatchEvent(&events.GuildVoiceMoveEvent{
			GenericGuildVoiceStateEvent: genericGuildVoiceEvent,
			OldVoiceState:               oldVoiceState,
		})
	} else if (oldOk || oldVoiceState.ChannelID == nil) && voiceState.ChannelID != nil {
		client.EventManager().DispatchEvent(&events.GuildVoiceJoinEvent{
			GenericGuildVoiceStateEvent: genericGuildVoiceEvent,
		})
	} else if voiceState.ChannelID == nil {
		client.EventManager().DispatchEvent(&events.GuildVoiceLeaveEvent{
			GenericGuildVoiceStateEvent: genericGuildVoiceEvent,
			OldVoiceState:               oldVoiceState,
		})
	} else {
		client.Logger().Warnf("could not decide which GuildVoiceEvent to fire")
	}
}
