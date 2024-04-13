package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerMessagePollVoteAdd(client bot.Client, sequenceNumber int, shardID int, event gateway.EventMessagePollVoteAdd) {
	genericEvent := events.NewGenericEvent(client, sequenceNumber, shardID)

	client.EventManager().DispatchEvent(&events.MessagePollVoteAdd{
		GenericMessagePollVote: &events.GenericMessagePollVote{
			GenericEvent: genericEvent,
			UserID:       event.UserID,
			ChannelID:    event.ChannelID,
			MessageID:    event.MessageID,
			GuildID:      event.GuildID,
			AnswerID:     event.AnswerID,
		},
	})

	if event.GuildID == nil {
		client.EventManager().DispatchEvent(&events.DMMessagePollVoteAdd{
			GenericDMMessagePollVote: &events.GenericDMMessagePollVote{
				GenericEvent: genericEvent,
				UserID:       event.UserID,
				ChannelID:    event.ChannelID,
				MessageID:    event.MessageID,
				AnswerID:     event.AnswerID,
			},
		})
	} else {
		client.EventManager().DispatchEvent(&events.GuildMessagePollVoteAdd{
			GenericGuildMessagePollVote: &events.GenericGuildMessagePollVote{
				GenericEvent: genericEvent,
				UserID:       event.UserID,
				ChannelID:    event.ChannelID,
				MessageID:    event.MessageID,
				GuildID:      *event.GuildID,
				AnswerID:     event.AnswerID,
			},
		})
	}
}

func gatewayHandlerMessagePollVoteRemove(client bot.Client, sequenceNumber int, shardID int, event gateway.EventMessagePollVoteRemove) {
	genericEvent := events.NewGenericEvent(client, sequenceNumber, shardID)

	client.EventManager().DispatchEvent(&events.MessagePollVoteRemove{
		GenericMessagePollVote: &events.GenericMessagePollVote{
			GenericEvent: genericEvent,
			UserID:       event.UserID,
			ChannelID:    event.ChannelID,
			MessageID:    event.MessageID,
			GuildID:      event.GuildID,
			AnswerID:     event.AnswerID,
		},
	})

	if event.GuildID == nil {
		client.EventManager().DispatchEvent(&events.DMMessagePollVoteRemove{
			GenericDMMessagePollVote: &events.GenericDMMessagePollVote{
				GenericEvent: genericEvent,
				UserID:       event.UserID,
				ChannelID:    event.ChannelID,
				MessageID:    event.MessageID,
				AnswerID:     event.AnswerID,
			},
		})
	} else {
		client.EventManager().DispatchEvent(&events.GuildMessagePollVoteRemove{
			GenericGuildMessagePollVote: &events.GenericGuildMessagePollVote{
				GenericEvent: genericEvent,
				UserID:       event.UserID,
				ChannelID:    event.ChannelID,
				MessageID:    event.MessageID,
				GuildID:      *event.GuildID,
				AnswerID:     event.AnswerID,
			},
		})
	}
}
