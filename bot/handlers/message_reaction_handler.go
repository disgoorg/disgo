package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerMessageReactionAdd(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventMessageReactionAdd) {
	client.EventManager.DispatchEvent(&events.MessageReactionAdd{
		GenericReaction: &events.GenericReaction{
			Event:        events.NewEvent(client),
			GatewayEvent: events.NewGatewayEvent(sequenceNumber, shardID),
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
		client.EventManager.DispatchEvent(&events.DMMessageReactionAdd{
			GenericDMMessageReaction: &events.GenericDMMessageReaction{
				Event:        events.NewEvent(client),
				GatewayEvent: events.NewGatewayEvent(sequenceNumber, shardID),
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
		client.EventManager.DispatchEvent(&events.GuildMessageReactionAdd{
			GenericGuildMessageReaction: &events.GenericGuildMessageReaction{
				Event:        events.NewEvent(client),
				GatewayEvent: events.NewGatewayEvent(sequenceNumber, shardID),
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

func gatewayHandlerMessageReactionRemove(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventMessageReactionRemove) {
	client.EventManager.DispatchEvent(&events.MessageReactionRemove{
		GenericReaction: &events.GenericReaction{
			Event:        events.NewEvent(client),
			GatewayEvent: events.NewGatewayEvent(sequenceNumber, shardID),
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
		client.EventManager.DispatchEvent(&events.DMMessageReactionRemove{
			GenericDMMessageReaction: &events.GenericDMMessageReaction{
				Event:        events.NewEvent(client),
				GatewayEvent: events.NewGatewayEvent(sequenceNumber, shardID),
				MessageID:    event.MessageID,
				ChannelID:    event.ChannelID,
				UserID:       event.UserID,
				Emoji:        event.Emoji,
				BurstColors:  event.BurstColors,
				Burst:        event.Burst,
			},
		})
	} else {
		client.EventManager.DispatchEvent(&events.GuildMessageReactionRemove{
			GenericGuildMessageReaction: &events.GenericGuildMessageReaction{
				Event:        events.NewEvent(client),
				GatewayEvent: events.NewGatewayEvent(sequenceNumber, shardID),
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

func gatewayHandlerMessageReactionRemoveAll(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventMessageReactionRemoveAll) {
	client.EventManager.DispatchEvent(&events.MessageReactionRemoveAll{
		Event:        events.NewEvent(client),
		GatewayEvent: events.NewGatewayEvent(sequenceNumber, shardID),
		MessageID:    event.MessageID,
		ChannelID:    event.ChannelID,
		GuildID:      event.GuildID,
	})

	if event.GuildID == nil {
		client.EventManager.DispatchEvent(&events.DMMessageReactionRemoveAll{
			Event:        events.NewEvent(client),
			GatewayEvent: events.NewGatewayEvent(sequenceNumber, shardID),
			MessageID:    event.MessageID,
			ChannelID:    event.ChannelID,
		})
	} else {
		client.EventManager.DispatchEvent(&events.GuildMessageReactionRemoveAll{
			Event:        events.NewEvent(client),
			GatewayEvent: events.NewGatewayEvent(sequenceNumber, shardID),
			MessageID:    event.MessageID,
			ChannelID:    event.ChannelID,
			GuildID:      *event.GuildID,
		})
	}
}

func gatewayHandlerMessageReactionRemoveEmoji(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventMessageReactionRemoveEmoji) {
	client.EventManager.DispatchEvent(&events.MessageReactionRemoveEmoji{
		Event:        events.NewEvent(client),
		GatewayEvent: events.NewGatewayEvent(sequenceNumber, shardID),
		MessageID:    event.MessageID,
		ChannelID:    event.ChannelID,
		GuildID:      event.GuildID,
		Emoji:        event.Emoji,
	})

	if event.GuildID == nil {
		client.EventManager.DispatchEvent(&events.DMMessageReactionRemoveEmoji{
			Event:        events.NewEvent(client),
			GatewayEvent: events.NewGatewayEvent(sequenceNumber, shardID),
			MessageID:    event.MessageID,
			ChannelID:    event.ChannelID,
			Emoji:        event.Emoji,
		})
	} else {
		client.EventManager.DispatchEvent(&events.GuildMessageReactionRemoveEmoji{
			Event:        events.NewEvent(client),
			GatewayEvent: events.NewGatewayEvent(sequenceNumber, shardID),
			MessageID:    event.MessageID,
			ChannelID:    event.ChannelID,
			GuildID:      *event.GuildID,
			Emoji:        event.Emoji,
		})
	}
}
