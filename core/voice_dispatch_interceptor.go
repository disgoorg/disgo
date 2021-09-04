package core

import "github.com/DisgoOrg/disgo/discord"

// VoiceServerUpdateEvent sent when a guilds voice server is updated
type VoiceServerUpdateEvent struct {
	discord.VoiceServerUpdate
	Disgo Disgo
}

// Guild returns the Guild for this VoiceServerUpdate from the Caches
func (u *VoiceServerUpdateEvent) Guild() *Guild {
	return u.Disgo.Caches().GuildCache().Get(u.GuildID)
}

// VoiceStateUpdateEvent sent when someone joins/leaves/moves voice channels
type VoiceStateUpdateEvent struct {
	*VoiceState
}

// Guild returns the Guild for this VoiceStateUpdate from the Caches
func (u *VoiceStateUpdateEvent) Guild() *Guild {
	return u.VoiceState.Guild()
}

// VoiceChannel returns the VoiceChannel for this VoiceStateUpdate from the Caches
func (u *VoiceStateUpdateEvent) VoiceChannel() VoiceChannel {
	return u.VoiceState.VoiceChannel()
}

// VoiceDispatchInterceptor lets you listen to VoiceServerUpdate & VoiceStateUpdate
type VoiceDispatchInterceptor interface {
	OnVoiceServerUpdate(voiceServerUpdateEvent *VoiceServerUpdateEvent)
	OnVoiceStateUpdate(voiceStateUpdateEvent *VoiceStateUpdateEvent)
}
