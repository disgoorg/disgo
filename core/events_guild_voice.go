package core

// GenericGuildVoiceEvent is called upon receiving GuildVoiceJoinEvent, GuildVoiceUpdateEvent, GuildVoiceLeaveEvent
type GenericGuildVoiceEvent struct {
	*GenericGuildMemberEvent
	VoiceState *VoiceState
}

// GuildVoiceJoinEvent indicates that an api.Member joined an api.VoiceChannel(requires api.GatewayIntentsGuildVoiceStates)
type GuildVoiceJoinEvent struct {
	*GenericGuildVoiceEvent
}

// GuildVoiceUpdateEvent indicates that an api.Member moved an api.VoiceChannel(requires api.GatewayIntentsGuildVoiceStates)
type GuildVoiceUpdateEvent struct {
	*GenericGuildVoiceEvent
	OldVoiceState *VoiceState
}

// GuildVoiceLeaveEvent indicates that an api.Member left an api.VoiceChannel(requires api.GatewayIntentsGuildVoiceStates)
type GuildVoiceLeaveEvent struct {
	*GenericGuildVoiceEvent
}
