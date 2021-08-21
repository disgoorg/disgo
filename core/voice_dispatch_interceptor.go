package core

import "github.com/DisgoOrg/disgo/discord"

// VoiceServerUpdateEvent sent when a guilds voice server is updated
type VoiceServerUpdateEvent struct {
	discord.VoiceServerUpdate
	Disgo Disgo
}

// Guild returns the Guild for this VoiceServerUpdate from the Cache
func (u *VoiceServerUpdateEvent) Guild() *Guild {
	return u.Disgo.Cache().GuildCache().Get(u.GuildID)
}

// VoiceStateUpdateEvent sent when someone joins/leaves/moves voice channels
type VoiceStateUpdateEvent struct {
	*VoiceState
}

// Guild returns the Guild for this VoiceStateUpdate from the Cache
func (u *VoiceStateUpdateEvent) Guild() *Guild {
	if u.GuildID == nil {
		return nil
	}
	return u.Disgo.Cache().GuildCache().Get(*u.GuildID)
}

// VoiceChannel returns the VoiceChannel for this VoiceStateUpdate from the Cache
func (u *VoiceStateUpdateEvent) VoiceChannel() VoiceChannel {
	if u.ChannelID == nil {
		return nil
	}
	return u.Disgo.Cache().VoiceChannelCache().Get(*u.ChannelID)
}

// VoiceDispatchInterceptor lets you listen to VoiceServerUpdate & VoiceStateUpdate
type VoiceDispatchInterceptor interface {
	OnVoiceServerUpdate(voiceServerUpdateEvent *VoiceServerUpdateEvent)
	OnVoiceStateUpdate(voiceStateUpdateEvent *VoiceStateUpdateEvent)
}
