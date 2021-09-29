package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

// VoiceServerUpdateEvent is sent when a guilds voice server is updated
type VoiceServerUpdateEvent struct {
	discord.VoiceServerUpdate
	Bot *Bot
}

// Guild returns the Guild for this VoiceServerUpdateEvent.
// This will only check cached guilds!
func (u *VoiceServerUpdateEvent) Guild() *Guild {
	return u.Bot.Caches.GuildCache().Get(u.GuildID)
}

// VoiceStateUpdateEvent is sent when someone joins/leaves/moves voice channels
type VoiceStateUpdateEvent struct {
	*VoiceState
}

// VoiceDispatchInterceptor lets you listen to VoiceServerUpdate & VoiceStateUpdate
type VoiceDispatchInterceptor interface {
	OnVoiceServerUpdate(voiceServerUpdateEvent *VoiceServerUpdateEvent)
	OnVoiceStateUpdate(voiceStateUpdateEvent *VoiceStateUpdateEvent)
}
