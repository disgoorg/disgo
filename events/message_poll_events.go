package events

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

// GenericMessagePollVote is a generic poll vote event (requires gateway.IntentGuildMessagePolls and/or gateway.IntentDirectMessagePolls)
type GenericMessagePollVote struct {
	*GenericEvent
	UserID    snowflake.ID
	ChannelID snowflake.ID
	MessageID snowflake.ID
	GuildID   *snowflake.ID
	AnswerID  int
}

// Guild returns the discord.Guild where the GenericMessagePollVote happened or empty if it happened in DMs
func (e *GenericMessagePollVote) Guild() (discord.Guild, bool) {
	if e.GuildID == nil {
		return discord.Guild{}, false
	}
	return e.Client().Caches().Guild(*e.GuildID)
}

// MessagePollVoteAdd indicates that a discord.User voted on a discord.Poll (requires gateway.IntentGuildMessagePolls and/or gateway.IntentDirectMessagePolls)
type MessagePollVoteAdd struct {
	*GenericMessagePollVote
}

// MessagePollVoteRemove indicates that a discord.User removed their vote on a discord.Poll (requires gateway.IntentGuildMessagePolls and/or gateway.IntentDirectMessagePolls)
type MessagePollVoteRemove struct {
	*GenericMessagePollVote
}
