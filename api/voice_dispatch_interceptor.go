package api

// VoiceServerUpdateEvent sent when a guilds voice server is updated
type VoiceServerUpdateEvent struct {
	*VoiceServerUpdate
	Disgo Disgo
}

// Guild returns the Guild for this VoiceServerUpdate from the Cache
func (u *VoiceServerUpdateEvent) Guild() *Guild {
	return u.Disgo.Cache().Guild(u.GuildID)
}

// VoiceStateUpdateEvent sent when someone joins/leaves/moves voice channels
type VoiceStateUpdateEvent struct {
	*VoiceState
	Member *Member `json:"member"`
}

// Guild returns the Guild for this VoiceStateUpdate from the Cache
func (u *VoiceStateUpdateEvent) Guild() *Guild {
	return u.Disgo.Cache().Guild(u.GuildID)
}

// VoiceChannel returns the VoiceChannel for this VoiceStateUpdate from the Cache
func (u *VoiceStateUpdateEvent) VoiceChannel() *VoiceChannel {
	if u.ChannelID == nil {
		return nil
	}
	return u.Disgo.Cache().VoiceChannel(*u.ChannelID)
}

// VoiceDispatchInterceptor lets you listen to VoiceServerUpdate & VoiceStateUpdate
type VoiceDispatchInterceptor interface {
	OnVoiceServerUpdate(voiceServerUpdateEvent *VoiceServerUpdateEvent)
	OnVoiceStateUpdate(voiceStateUpdateEvent *VoiceStateUpdateEvent)
}
