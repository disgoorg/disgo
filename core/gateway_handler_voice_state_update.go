package core

import (
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
	return &discord.VoiceState{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *VoiceStateUpdateHandler) HandleGatewayEvent(bot *Bot, sequenceNumber int, v interface{}) {
	discordVoiceState := *v.(*discord.VoiceState)

	oldVoiceState := bot.Caches.VoiceStateCache().GetCopy(discordVoiceState.GuildID, discordVoiceState.UserID)

	var member *Member
	if discordVoiceState.Member != nil {
		member = bot.EntityBuilder.CreateMember(discordVoiceState.Member.GuildID, *discordVoiceState.Member, CacheStrategyYes)
	}
	voiceState := bot.EntityBuilder.CreateVoiceState(discordVoiceState.GuildID, discordVoiceState, CacheStrategyYes)

	// voice state update for ourselves received
	// execute voice VoiceDispatchInterceptor.OnVoiceStateUpdate
	if bot.ClientID == discordVoiceState.UserID {
		if interceptor := bot.VoiceDispatchInterceptor; interceptor != nil {
			interceptor.OnVoiceStateUpdate(&VoiceStateUpdateEvent{VoiceState: voiceState})
		}
	}

	guild := voiceState.Guild()
	if guild == nil {
		bot.Logger.Warnf("received guild voice state update for unknown guild: %s", discordVoiceState.GuildID)
		return
	}

	genericGuildVoiceEvent := &GenericGuildVoiceEvent{
		GenericGuildMemberEvent: &GenericGuildMemberEvent{
			GenericGuildEvent: &GenericGuildEvent{
				GenericEvent: NewGenericEvent(bot, sequenceNumber),
				Guild:        guild,
			},
			Member: member,
		},
		VoiceState: voiceState,
	}

	if (oldVoiceState == nil || oldVoiceState.ChannelID == nil) && discordVoiceState.ChannelID != nil {
		bot.EventManager.Dispatch(&GuildVoiceJoinEvent{GenericGuildVoiceEvent: genericGuildVoiceEvent})
	} else if oldVoiceState != nil && oldVoiceState.ChannelID != nil && discordVoiceState.ChannelID == nil {
		bot.EventManager.Dispatch(&GuildVoiceLeaveEvent{GenericGuildVoiceEvent: genericGuildVoiceEvent})
	} else if oldVoiceState != nil && oldVoiceState.ChannelID != nil && discordVoiceState.ChannelID != nil {
		bot.EventManager.Dispatch(&GuildVoiceUpdateEvent{GenericGuildVoiceEvent: genericGuildVoiceEvent, OldVoiceState: oldVoiceState})
	} else {
		bot.Logger.Warnf("could not decide which GuildVoiceEvent to fire")
	}
}
