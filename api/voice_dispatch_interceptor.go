package api

// VoiceServerUpdate sent when a guilds voice server is updated
type VoiceServerUpdate struct {
	Disgo    Disgo
	Token    string    `json:"token"`
	GuildID  Snowflake `json:"guild_id"`
	Endpoint *string   `json:"endpoint"`
}

// Guild returns the Guild for this VoiceServerUpdate from the Cache
func (u VoiceServerUpdate) Guild() *Guild {
	return u.Disgo.Cache().Guild(u.GuildID)
}

// VoiceStateUpdate sent when someone joins/leaves/moves voice channels
type VoiceStateUpdate struct {
	VoiceState
	Member *Member `json:"member"`
}

// Guild returns the Guild for this VoiceStateUpdate from the Cache
func (u VoiceStateUpdate) Guild() *Guild {
	return u.Disgo.Cache().Guild(u.GuildID)
}

// VoiceChannel returns the VoiceChannel for this VoiceStateUpdate from the Cache
func (u VoiceStateUpdate) VoiceChannel() *VoiceChannel {
	if u.ChannelID == nil {
		return nil
	}
	return u.Disgo.Cache().VoiceChannel(*u.ChannelID)
}

// VoiceDispatchInterceptor lets you listen to VoiceServerUpdate & VoiceStateUpdate
type VoiceDispatchInterceptor interface {
	OnVoiceServerUpdate(voiceServerUpdateEvent VoiceServerUpdate)
	OnVoiceStateUpdate(voiceStateUpdateEvent VoiceStateUpdate)
}
