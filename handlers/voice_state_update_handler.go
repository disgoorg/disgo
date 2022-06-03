package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

type gatewayHandlerVoiceStateUpdate struct{}

func (h *gatewayHandlerVoiceStateUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeVoiceStateUpdate
}

func (h *gatewayHandlerVoiceStateUpdate) New() any {
	return &discord.VoiceStateUpdate{}
}

func (h *gatewayHandlerVoiceStateUpdate) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	voiceState := *v.(*discord.VoiceStateUpdate)
	member := voiceState.Member

	oldVoiceState, oldOk := client.Caches().VoiceStates().Get(voiceState.GuildID, voiceState.UserID)
	if voiceState.ChannelID == nil {
		client.Caches().VoiceStates().Remove(voiceState.GuildID, voiceState.UserID)
	} else {
		client.Caches().VoiceStates().Put(voiceState.GuildID, voiceState.UserID, voiceState.VoiceState)
	}
	client.Caches().Members().Put(voiceState.GuildID, voiceState.UserID, member)

	if voiceState.UserID == client.ID() && client.VoiceManager() != nil {
		client.VoiceManager().HandleVoiceStateUpdate(voiceState.VoiceState)
	}

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
