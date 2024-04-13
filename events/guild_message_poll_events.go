package events

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

// GenericGuildMessagePollVote is called upon receiving GuildMessagePollVoteAdd or GuildMessagePollVoteRemove (requires the gateway.IntentGuildMessagePolls)
type GenericGuildMessagePollVote struct {
	*GenericEvent
	UserID    snowflake.ID
	ChannelID snowflake.ID
	MessageID snowflake.ID
	GuildID   snowflake.ID
	AnswerID  int
}

// Guild returns the discord.Guild where the GenericGuildMessagePollVote happened
func (e *GenericGuildMessagePollVote) Guild() (discord.Guild, bool) {
	return e.Client().Caches().Guild(e.GuildID)
}

// GuildMessagePollVoteAdd indicates that a discord.User voted on a discord.Poll in a discord.Guild (requires gateway.IntentGuildMessagePolls)
type GuildMessagePollVoteAdd struct {
	*GenericGuildMessagePollVote
}

// GuildMessagePollVoteRemove indicates that a discord.User removed their vote on a discord.Poll in a discord.Guild (requires gateway.IntentGuildMessagePolls)
type GuildMessagePollVoteRemove struct {
	*GenericGuildMessagePollVote
}
