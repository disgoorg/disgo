package events

import (
	"github.com/disgoorg/snowflake/v2"
)

// GenericDMMessagePollVote is called upon receiving DMMessagePollVoteAdd or DMMessagePollVoteRemove (requires gateway.IntentDirectMessagePolls)
type GenericDMMessagePollVote struct {
	*GenericEvent
	UserID    snowflake.ID
	ChannelID snowflake.ID
	MessageID snowflake.ID
	AnswerID  int
}

// DMMessagePollVoteAdd  indicates that a discord.User voted on a discord.Poll in a DM (requires gateway.IntentDirectMessagePolls)
type DMMessagePollVoteAdd struct {
	*GenericDMMessagePollVote
}

// DMMessagePollVoteRemove indicates that a discord.User removed their vote on a discord.Poll in a DM (requires gateway.IntentDirectMessagePolls)
type DMMessagePollVoteRemove struct {
	*GenericDMMessagePollVote
}
