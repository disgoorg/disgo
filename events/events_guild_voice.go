package events

import "github.com/DisgoOrg/disgo/core"

// GenericGuildVoiceEvent is called upon receiving GuildVoiceJoinEvent, GuildVoiceUpdateEvent, GuildVoiceLeaveEvent
type GenericGuildVoiceEvent struct {
	*GenericGuildMemberEvent
	VoiceState *core.VoiceState
}

// GuildVoiceJoinEvent indicates that a core.Member joined a Channel (requires discord.GatewayIntentGuildVoiceStates)
type GuildVoiceJoinEvent struct {
	*GenericGuildVoiceEvent
}

// GuildVoiceUpdateEvent indicates that a core.Member updated their VoiceState (requires discord.GatewayIntentGuildVoiceStates)
type GuildVoiceUpdateEvent struct {
	*GenericGuildVoiceEvent
	OldVoiceState *core.VoiceState
}

// GuildVoiceLeaveEvent indicates that a core.Member left a Channel (requires discord.GatewayIntentGuildVoiceStates)
type GuildVoiceLeaveEvent struct {
	*GenericGuildVoiceEvent
	OldVoiceState *core.VoiceState
}
