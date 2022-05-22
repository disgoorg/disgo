package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

// gatewayHandlerVoiceStateUpdate handles discord.GatewayEventTypeVoiceStateUpdate
type gatewayHandlerVoiceStateUpdate struct{}

// EventType returns the discord.GatewayEventTypeVoiceStateUpdate
func (h *gatewayHandlerVoiceStateUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeVoiceStateUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerVoiceStateUpdate) New() any {
	return &discord.FullVoiceState{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerVoiceStateUpdate) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	voiceState := *v.(*discord.FullVoiceState)
	member := voiceState.Member

	oldVoiceState, oldOk := client.Caches().VoiceStates().Get(voiceState.GuildID, voiceState.UserID)
	if voiceState.ChannelID == nil {
		client.Caches().VoiceStates().Remove(voiceState.GuildID, voiceState.UserID)
	} else {
		client.Caches().VoiceStates().Put(voiceState.GuildID, voiceState.UserID, voiceState.VoiceState)
	}
	client.Caches().Members().Put(voiceState.GuildID, voiceState.UserID, member)

	genericGuildVoiceEvent := &events.GenericGuildVoiceState{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		VoiceState:   voiceState.VoiceState,
		Member:       member,
	}

	client.EventManager().DispatchEvent(&events.GuildVoiceStateUpdate{
		GenericGuildVoiceState: genericGuildVoiceEvent,
		OldVoiceState:          oldVoiceState,
	})

	if oldOk && oldVoiceState.ChannelID != nil && voiceState.ChannelID != nil {
		client.EventManager().DispatchEvent(&events.GuildVoiceMove{
			GenericGuildVoiceState: genericGuildVoiceEvent,
			OldVoiceState:          oldVoiceState,
		})
	} else if (oldOk || oldVoiceState.ChannelID == nil) && voiceState.ChannelID != nil {
		client.EventManager().DispatchEvent(&events.GuildVoiceJoin{
			GenericGuildVoiceState: genericGuildVoiceEvent,
		})
	} else if voiceState.ChannelID == nil {
		client.EventManager().DispatchEvent(&events.GuildVoiceLeave{
			GenericGuildVoiceState: genericGuildVoiceEvent,
			OldVoiceState:          oldVoiceState,
		})
	} else {
		client.Logger().Warnf("could not decide which GuildVoice to fire")
	}
}
