package core

// GenericGuildVoiceEvent is called upon receiving GuildVoiceJoinEvent, GuildVoiceUpdateEvent, GuildVoiceLeaveEvent
type GenericGuildVoiceEvent struct {
	*GenericGuildMemberEvent
	VoiceState *VoiceState
}

// GuildVoiceJoinEvent indicates that a core.Member joined a Channel (requires discord.GatewayIntentGuildVoiceStates)
type GuildVoiceJoinEvent struct {
	*GenericGuildVoiceEvent
}

// GuildVoiceUpdateEvent indicates that a core.Member updated their VoiceState (requires discord.GatewayIntentGuildVoiceStates)
type GuildVoiceUpdateEvent struct {
	*GenericGuildVoiceEvent
	OldVoiceState *VoiceState
}

// GuildVoiceLeaveEvent indicates that a core.Member left a Channel (requires discord.GatewayIntentGuildVoiceStates)
type GuildVoiceLeaveEvent struct {
	*GenericGuildVoiceEvent
	OldVoiceState *VoiceState
}
