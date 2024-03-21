package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerMessagePollVoteAdd(client bot.Client, sequenceNumber int, shardID int, event gateway.EventMessagePollVoteAdd) {
	client.EventManager().DispatchEvent(&events.MessagePollVoteAdd{
		GenericMessagePoll: &events.GenericMessagePoll{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			UserID:       event.UserID,
			ChannelID:    event.ChannelID,
			MessageID:    event.MessageID,
			GuildID:      event.GuildID,
			AnswerID:     event.AnswerID,
		},
	})
}

func gatewayHandlerMessagePollVoteRemove(client bot.Client, sequenceNumber int, shardID int, event gateway.EventMessagePollVoteRemove) {
	client.EventManager().DispatchEvent(&events.MessagePollVoteRemove{
		GenericMessagePoll: &events.GenericMessagePoll{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			UserID:       event.UserID,
			ChannelID:    event.ChannelID,
			MessageID:    event.MessageID,
			GuildID:      event.GuildID,
			AnswerID:     event.AnswerID,
		},
	})
}
