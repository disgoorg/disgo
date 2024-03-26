package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerMessageReactionAdd(client bot.Client, sequenceNumber int, shardID int, event gateway.EventMessageReactionAdd) {
	genericEvent := events.NewGenericEvent(client, sequenceNumber, shardID)

	client.EventManager().DispatchEvent(&events.MessageReactionAdd{
		GenericReaction: &events.GenericReaction{
			GenericEvent: genericEvent,
			MessageID:    event.MessageID,
			ChannelID:    event.ChannelID,
			GuildID:      event.GuildID,
			UserID:       event.UserID,
			Emoji:        event.Emoji,
			BurstColors:  event.BurstColors,
			Burst:        event.Burst,
		},
		Member: event.Member,
	})

	if event.GuildID == nil {
		client.EventManager().DispatchEvent(&events.DMMessageReactionAdd{
			GenericDMMessageReaction: &events.GenericDMMessageReaction{
				GenericEvent: genericEvent,
				MessageID:    event.MessageID,
				ChannelID:    event.ChannelID,
				UserID:       event.UserID,
				Emoji:        event.Emoji,
				BurstColors:  event.BurstColors,
				Burst:        event.Burst,
			},
			MessageAuthorID: event.MessageAuthorID,
		})
	} else {
		var member discord.Member
		// sometimes the member is nil for some reason
		if event.Member != nil {
			member = *event.Member
		}
		client.EventManager().DispatchEvent(&events.GuildMessageReactionAdd{
			GenericGuildMessageReaction: &events.GenericGuildMessageReaction{
				GenericEvent: genericEvent,
				MessageID:    event.MessageID,
				ChannelID:    event.ChannelID,
				GuildID:      *event.GuildID,
				UserID:       event.UserID,
				Emoji:        event.Emoji,
				BurstColors:  event.BurstColors,
				Burst:        event.Burst,
			},
			Member:          member,
			MessageAuthorID: event.MessageAuthorID,
		})
	}
}

func gatewayHandlerMessageReactionRemove(client bot.Client, sequenceNumber int, shardID int, event gateway.EventMessageReactionRemove) {
	genericEvent := events.NewGenericEvent(client, sequenceNumber, shardID)

	client.EventManager().DispatchEvent(&events.MessageReactionRemove{
		GenericReaction: &events.GenericReaction{
			GenericEvent: genericEvent,
			MessageID:    event.MessageID,
			ChannelID:    event.ChannelID,
			GuildID:      event.GuildID,
			UserID:       event.UserID,
			Emoji:        event.Emoji,
			BurstColors:  event.BurstColors,
			Burst:        event.Burst,
		},
	})

	if event.GuildID == nil {
		client.EventManager().DispatchEvent(&events.DMMessageReactionRemove{
			GenericDMMessageReaction: &events.GenericDMMessageReaction{
				GenericEvent: genericEvent,
				MessageID:    event.MessageID,
				ChannelID:    event.ChannelID,
				UserID:       event.UserID,
				Emoji:        event.Emoji,
				BurstColors:  event.BurstColors,
				Burst:        event.Burst,
			},
		})
	} else {
		client.EventManager().DispatchEvent(&events.GuildMessageReactionRemove{
			GenericGuildMessageReaction: &events.GenericGuildMessageReaction{
				GenericEvent: genericEvent,
				MessageID:    event.MessageID,
				ChannelID:    event.ChannelID,
				GuildID:      *event.GuildID,
				UserID:       event.UserID,
				Emoji:        event.Emoji,
				BurstColors:  event.BurstColors,
				Burst:        event.Burst,
			},
		})
	}
}

func gatewayHandlerMessageReactionRemoveAll(client bot.Client, sequenceNumber int, shardID int, event gateway.EventMessageReactionRemoveAll) {
	genericEvent := events.NewGenericEvent(client, sequenceNumber, shardID)

	client.EventManager().DispatchEvent(&events.MessageReactionRemoveAll{
		GenericEvent: genericEvent,
		MessageID:    event.MessageID,
		ChannelID:    event.ChannelID,
		GuildID:      event.GuildID,
	})

	if event.GuildID == nil {
		client.EventManager().DispatchEvent(&events.DMMessageReactionRemoveAll{
			GenericEvent: genericEvent,
			MessageID:    event.MessageID,
			ChannelID:    event.ChannelID,
		})
	} else {
		client.EventManager().DispatchEvent(&events.GuildMessageReactionRemoveAll{
			GenericEvent: genericEvent,
			MessageID:    event.MessageID,
			ChannelID:    event.ChannelID,
			GuildID:      *event.GuildID,
		})
	}
}

func gatewayHandlerMessageReactionRemoveEmoji(client bot.Client, sequenceNumber int, shardID int, event gateway.EventMessageReactionRemoveEmoji) {
	genericEvent := events.NewGenericEvent(client, sequenceNumber, shardID)

	client.EventManager().DispatchEvent(&events.MessageReactionRemoveEmoji{
		GenericEvent: genericEvent,
		MessageID:    event.MessageID,
		ChannelID:    event.ChannelID,
		GuildID:      event.GuildID,
		Emoji:        event.Emoji,
	})

	if event.GuildID == nil {
		client.EventManager().DispatchEvent(&events.DMMessageReactionRemoveEmoji{
			GenericEvent: genericEvent,
			MessageID:    event.MessageID,
			ChannelID:    event.ChannelID,
			Emoji:        event.Emoji,
		})
	} else {
		client.EventManager().DispatchEvent(&events.GuildMessageReactionRemoveEmoji{
			GenericEvent: genericEvent,
			MessageID:    event.MessageID,
			ChannelID:    event.ChannelID,
			GuildID:      *event.GuildID,
			Emoji:        event.Emoji,
		})
	}
}
