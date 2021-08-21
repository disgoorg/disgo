package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
)

// VoiceStateUpdateHandler handles api.GatewayEventVoiceStateUpdate
type VoiceStateUpdateHandler struct{}

// EventType returns the gateway.EventType
func (h *VoiceStateUpdateHandler) EventType() gateway.EventType {
	return gateway.EventTypeVoiceStateUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *VoiceStateUpdateHandler) New() interface{} {
	return discord.VoiceState{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *VoiceStateUpdateHandler) HandleGatewayEvent(disgo core.Disgo, _ core.EventManager, sequenceNumber int, i interface{}) {
	discordVoiceState, ok := i.(discord.VoiceState)
	if !ok {
		return
	}

	oldVoiceState := disgo.Cache().VoiceStateCache().Get(*discordVoiceState.GuildID, discordVoiceState.UserID)
	if oldVoiceState != nil {
		oldVoiceState = &*oldVoiceState
	}
	var member *core.Member
	if discordVoiceState.Member != nil {
		member = disgo.EntityBuilder().CreateMember(discordVoiceState.Member.GuildID, *discordVoiceState.Member, core.CacheStrategyYes)
	}
	voiceState := disgo.EntityBuilder().CreateVoiceState(*discordVoiceState.GuildID, discordVoiceState, core.CacheStrategyYes)

	// voice state update for ourselves received
	// execute voice VoiceDispatchInterceptor.OnVoiceStateUpdate
	if disgo.ClientID() == discordVoiceState.UserID {
		if interceptor := disgo.VoiceDispatchInterceptor(); interceptor != nil {
			interceptor.OnVoiceStateUpdate(&core.VoiceStateUpdateEvent{VoiceState: voiceState})
		}
	}

	guild := voiceState.Guild()
	if guild == nil {
		disgo.Logger().Warnf("received guild voice state update for unknown guild: %s", discordVoiceState.GuildID)
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

	if (oldVoiceState == nil || oldVoiceState.ChannelID == nil) && discordVoiceState.ChannelID != nil {
		disgo.EventManager().Dispatch(&events.GuildVoiceJoinEvent{GenericGuildVoiceEvent: genericGuildVoiceEvent})
	} else if oldVoiceState != nil && oldVoiceState.ChannelID != nil && discordVoiceState.ChannelID == nil {
		disgo.EventManager().Dispatch(&events.GuildVoiceLeaveEvent{GenericGuildVoiceEvent: genericGuildVoiceEvent})
	} else if oldVoiceState != nil && oldVoiceState.ChannelID != nil && discordVoiceState.ChannelID != nil {
		disgo.EventManager().Dispatch(&events.GuildVoiceUpdateEvent{GenericGuildVoiceEvent: genericGuildVoiceEvent, OldVoiceState: oldVoiceState})
	} else {
		disgo.Logger().Warnf("could not decide which GuildVoiceEvent to fire")
	}
}
