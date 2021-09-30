package gatewayhandlers

import (
	"github.com/DisgoOrg/disgo/core"
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
func (h *gatewayHandlerVoiceStateUpdate) New() interface{} {
	return &discord.VoiceState{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerVoiceStateUpdate) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	voiceState := *v.(*discord.VoiceState)

	oldVoiceState := bot.Caches.VoiceStateCache().GetCopy(voiceState.GuildID, voiceState.UserID)

	var member *core.Member
	if voiceState.Member != nil {
		member = bot.EntityBuilder.CreateMember(voiceState.Member.GuildID, *voiceState.Member, core.CacheStrategyYes)
	}
	coreVoiceState := bot.EntityBuilder.CreateVoiceState(voiceState, core.CacheStrategyYes)

	if oldVoiceState != nil && oldVoiceState.ChannelID != nil {
		if channel := bot.Caches.ChannelCache().Get(*oldVoiceState.ChannelID); channel != nil {
			delete(channel.ConnectedMemberIDs, coreVoiceState.UserID)
		}
	}

	if coreVoiceState.ChannelID != nil {
		if channel := bot.Caches.ChannelCache().Get(*coreVoiceState.ChannelID); channel != nil {
			channel.ConnectedMemberIDs[coreVoiceState.UserID] = struct{}{}
		}
	}

	// voice state update for ourselves received
	// execute voice VoiceDispatchInterceptor.OnVoiceStateUpdate
	if bot.ClientID == voiceState.UserID {
		if interceptor := bot.EventManager.Config().VoiceDispatchInterceptor; interceptor != nil {
			interceptor.OnVoiceStateUpdate(&core.VoiceStateUpdateEvent{VoiceState: coreVoiceState})
		}
	}

	guild := coreVoiceState.Guild()
	if guild == nil {
		bot.Logger.Warnf("received guild voice state update for unknown guild: %s", voiceState.GuildID)
		return
	}

	genericGuildVoiceEvent := &events.GenericGuildVoiceEvent{
		GenericGuildMemberEvent: &events.GenericGuildMemberEvent{
			GenericGuildEvent: &events.GenericGuildEvent{
				GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
				Guild:        guild,
			},
			Member: member,
		},
		VoiceState: coreVoiceState,
	}

	if oldVoiceState != nil && oldVoiceState.ChannelID != nil && voiceState.ChannelID != nil {
		bot.EventManager.Dispatch(&events.GuildVoiceUpdateEvent{GenericGuildVoiceEvent: genericGuildVoiceEvent, OldVoiceState: oldVoiceState})
	} else if (oldVoiceState == nil || oldVoiceState.ChannelID == nil) && voiceState.ChannelID != nil {
		bot.EventManager.Dispatch(&events.GuildVoiceJoinEvent{GenericGuildVoiceEvent: genericGuildVoiceEvent})
	} else if voiceState.ChannelID == nil {
		bot.EventManager.Dispatch(&events.GuildVoiceLeaveEvent{GenericGuildVoiceEvent: genericGuildVoiceEvent, OldVoiceState: oldVoiceState})
	} else {
		bot.Logger.Warnf("could not decide which GuildVoiceEvent to fire")
	}
}
