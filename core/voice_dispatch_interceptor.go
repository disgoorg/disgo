package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

// VoiceServerUpdateEvent sent when a guilds voice server is updated
type VoiceServerUpdateEvent struct {
	discord.VoiceServerUpdate
	Bot *Bot
}

// Guild returns the Guild for this VoiceServerUpdate from the Caches
func (u *VoiceServerUpdateEvent) Guild() *Guild {
	return u.Bot.Caches.GuildCache().Get(u.GuildID)
}

// VoiceStateUpdateEvent sent when someone joins/leaves/moves voice channels
type VoiceStateUpdateEvent struct {
	*VoiceState
}

// VoiceDispatchInterceptor lets you listen to VoiceServerUpdate & VoiceStateUpdate
type VoiceDispatchInterceptor interface {
	OnVoiceServerUpdate(voiceServerUpdateEvent *VoiceServerUpdateEvent)
	OnVoiceStateUpdate(voiceStateUpdateEvent *VoiceStateUpdateEvent)
}
