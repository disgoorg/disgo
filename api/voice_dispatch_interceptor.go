package api

// VoiceServerUpdateEvent sent when a guilds voice server is updated
type VoiceServerUpdateEvent struct {
	VoiceServerUpdate *VoiceServerUpdate
	Disgo             Disgo
}

// Guild returns the Guild for this VoiceServerUpdate from the Cache
func (u *VoiceServerUpdateEvent) Guild() *Guild {
	return u.Disgo.Cache().Guild(u.VoiceServerUpdate.GuildID)
}

// VoiceStateUpdateEvent sent when someone joins/leaves/moves voice channels
type VoiceStateUpdateEvent struct {
	VoiceState *VoiceState
	Member     *Member `json:"member"`
}

// Guild returns the Guild for this VoiceStateUpdate from the Cache
func (u *VoiceStateUpdateEvent) Guild() *Guild {
	return u.VoiceState.Disgo.Cache().Guild(u.VoiceState.GuildID)
}

// VoiceChannel returns the VoiceChannel for this VoiceStateUpdate from the Cache
func (u *VoiceStateUpdateEvent) VoiceChannel() VoiceChannel {
	if u.VoiceState.ChannelID == nil {
		return nil
	}
	return u.VoiceState.Disgo.Cache().VoiceChannel(*u.VoiceState.ChannelID)
}

// VoiceDispatchInterceptor lets you listen to VoiceServerUpdate & VoiceStateUpdate
type VoiceDispatchInterceptor interface {
	OnVoiceServerUpdate(voiceServerUpdateEvent *VoiceServerUpdateEvent)
	OnVoiceStateUpdate(voiceStateUpdateEvent *VoiceStateUpdateEvent)
}
