package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

// VoiceStateUpdateHandler handles api.GatewayEventVoiceStateUpdate
type VoiceStateUpdateHandler struct{}

// EventType returns the gateway.EventType
func (h *VoiceStateUpdateHandler) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeVoiceStateUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *VoiceStateUpdateHandler) New() interface{} {
	return discord.VoiceState{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *VoiceStateUpdateHandler) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	discordVoiceState, ok := v.(discord.VoiceState)
	if !ok {
		return
	}

	oldVoiceState := bot.Caches.VoiceStateCache().GetCopy(discordVoiceState.GuildID, discordVoiceState.UserID)

	var member *core.Member
	if discordVoiceState.Member != nil {
		member = bot.EntityBuilder.CreateMember(discordVoiceState.Member.GuildID, *discordVoiceState.Member, core.CacheStrategyYes)
	}
	voiceState := bot.EntityBuilder.CreateVoiceState(discordVoiceState.GuildID, discordVoiceState, core.CacheStrategyYes)

	// voice state update for ourselves received
	// execute voice VoiceDispatchInterceptor.OnVoiceStateUpdate
	if bot.ClientID == discordVoiceState.UserID {
		if interceptor := bot.VoiceDispatchInterceptor; interceptor != nil {
			interceptor.OnVoiceStateUpdate(&core.VoiceStateUpdateEvent{VoiceState: voiceState})
		}
	}

	guild := voiceState.Guild()
	if guild == nil {
		bot.Logger.Warnf("received guild voice state update for unknown guild: %s", discordVoiceState.GuildID)
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
		VoiceState: voiceState,
	}

	if (oldVoiceState == nil || oldVoiceState.ChannelID == nil) && discordVoiceState.ChannelID != nil {
		bot.EventManager.Dispatch(&events.GuildVoiceJoinEvent{GenericGuildVoiceEvent: genericGuildVoiceEvent})
	} else if oldVoiceState != nil && oldVoiceState.ChannelID != nil && discordVoiceState.ChannelID == nil {
		bot.EventManager.Dispatch(&events.GuildVoiceLeaveEvent{GenericGuildVoiceEvent: genericGuildVoiceEvent})
	} else if oldVoiceState != nil && oldVoiceState.ChannelID != nil && discordVoiceState.ChannelID != nil {
		bot.EventManager.Dispatch(&events.GuildVoiceUpdateEvent{GenericGuildVoiceEvent: genericGuildVoiceEvent, OldVoiceState: oldVoiceState})
	} else {
		bot.Logger.Warnf("could not decide which GuildVoiceEvent to fire")
	}
}
