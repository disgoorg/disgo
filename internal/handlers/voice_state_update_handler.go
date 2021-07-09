package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

// VoiceStateUpdateHandler handles api.GatewayEventVoiceStateUpdate
type VoiceStateUpdateHandler struct{}

// Event returns the raw gateway event Event
func (h *VoiceStateUpdateHandler) Event() api.GatewayEventType {
	return api.GatewayEventVoiceStateUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *VoiceStateUpdateHandler) New() interface{} {
	return &api.VoiceStateUpdateEvent{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *VoiceStateUpdateHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
	voiceStateUpdate, ok := i.(*api.VoiceStateUpdateEvent)
	if !ok {
		return
	}

	oldVoiceState := disgo.Cache().VoiceState(voiceStateUpdate.VoiceState.GuildID, voiceStateUpdate.VoiceState.UserID)
	if oldVoiceState != nil {
		oldVoiceState = &*oldVoiceState
	}
	member := disgo.EntityBuilder().CreateMember(voiceStateUpdate.Member.GuildID, voiceStateUpdate.Member, api.CacheStrategyYes)
	voiceState := disgo.EntityBuilder().CreateVoiceState(voiceStateUpdate.VoiceState.GuildID, voiceStateUpdate.VoiceState, api.CacheStrategyYes)

	// voice state update for ourself received
	// execute voice VoiceDispatchInterceptor.OnVoiceStateUpdate
	if disgo.ClientID() == voiceStateUpdate.VoiceState.UserID {
		if interceptor := disgo.VoiceDispatchInterceptor(); interceptor != nil {
			interceptor.OnVoiceStateUpdate(voiceStateUpdate)
		}
	}

	guild := voiceStateUpdate.Guild()
	if guild == nil {
		disgo.Logger().Warnf("received guild voice state update for unknown guild: %s", voiceStateUpdate.VoiceState.GuildID)
		return
	}

	genericGuildVoiceEvent := &events.GenericGuildVoiceEvent{
		GenericGuildMemberEvent: &events.GenericGuildMemberEvent{
			GenericGuildEvent: &events.GenericGuildEvent{
				GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
				Guild:        guild,
			},
			Member: member,
		},
		VoiceState: voiceState,
	}

	if (oldVoiceState == nil || oldVoiceState.ChannelID == nil) && voiceStateUpdate.VoiceState.ChannelID != nil {
		disgo.EventManager().Dispatch(&events.GuildVoiceJoinEvent{GenericGuildVoiceEvent: genericGuildVoiceEvent})
	} else if oldVoiceState != nil && oldVoiceState.ChannelID != nil && voiceStateUpdate.VoiceState.ChannelID == nil {
		disgo.EventManager().Dispatch(&events.GuildVoiceLeaveEvent{GenericGuildVoiceEvent: genericGuildVoiceEvent})
	} else if oldVoiceState != nil && oldVoiceState.ChannelID != nil && voiceStateUpdate.VoiceState.ChannelID != nil {
		disgo.EventManager().Dispatch(&events.GuildVoiceUpdateEvent{GenericGuildVoiceEvent: genericGuildVoiceEvent, OldVoiceState: oldVoiceState})
	} else {
		disgo.Logger().Warnf("could not decide which GuildVoiceEvent to fire")
	}
}
