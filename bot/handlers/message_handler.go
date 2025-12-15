package handlers

import (
	"errors"
	"log/slog"

	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerMessageCreate(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventMessageCreate) {
	if event.Flags.Has(discord.MessageFlagEphemeral) {
		// Ignore ephemeral messages as they miss guild_id & member
		return
	}

	if event.Member != nil {
		event.Member.User = event.Author
	}

	client.Caches.AddMessage(event.Message)

	channel, err := client.Caches.GuildMessageChannel(event.ChannelID)
	if err != nil && !errors.Is(err, cache.ErrNotFound) {
		client.Logger.Error("failed to get channel from cache", slog.Any("err", err), slog.String("channel_id", event.ChannelID.String()))
	}
	if err == nil {
		client.Caches.AddChannel(discord.ApplyLastMessageIDToChannel(channel, event.ID))
	}

	thread, err := client.Caches.GuildThread(event.ChannelID)
	if err != nil && !errors.Is(err, cache.ErrNotFound) {
		client.Logger.Error("failed to get thread from cache", slog.Any("err", err), slog.String("channel_id", event.ChannelID.String()))
	}
	if err == nil {
		thread.TotalMessageSent++
		thread.MessageCount++
		client.Caches.AddChannel(thread)
	}

	genericEvent := events.NewGenericEvent(client, sequenceNumber, shardID)
	client.EventManager.DispatchEvent(&events.MessageCreate{
		GenericMessage: &events.GenericMessage{
			GenericEvent: genericEvent,
			MessageID:    event.ID,
			Message:      event.Message,
			ChannelID:    event.ChannelID,
			GuildID:      event.GuildID,
		},
	})

	if event.GuildID == nil {
		client.EventManager.DispatchEvent(&events.DMMessageCreate{
			GenericDMMessage: &events.GenericDMMessage{
				GenericEvent: genericEvent,
				MessageID:    event.ID,
				Message:      event.Message,
				ChannelID:    event.ChannelID,
			},
		})
	} else {
		client.EventManager.DispatchEvent(&events.GuildMessageCreate{
			GenericGuildMessage: &events.GenericGuildMessage{
				GenericEvent: genericEvent,
				MessageID:    event.ID,
				Message:      event.Message,
				ChannelID:    event.ChannelID,
				GuildID:      *event.GuildID,
			},
		})
	}
}

func gatewayHandlerMessageUpdate(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventMessageUpdate) {
	oldMessage, err := client.Caches.Message(event.ChannelID, event.ID)
	if err != nil && !errors.Is(err, cache.ErrNotFound) {
		client.Logger.Error("failed to get message from cache", slog.Any("err", err), slog.String("channel_id", event.ChannelID.String()), slog.String("message_id", event.ID.String()))
	}
	client.Caches.AddMessage(event.Message)

	genericEvent := events.NewGenericEvent(client, sequenceNumber, shardID)
	client.EventManager.DispatchEvent(&events.MessageUpdate{
		GenericMessage: &events.GenericMessage{
			GenericEvent: genericEvent,
			MessageID:    event.ID,
			Message:      event.Message,
			ChannelID:    event.ChannelID,
			GuildID:      event.GuildID,
		},
		OldMessage: oldMessage,
	})

	if event.GuildID == nil {
		client.EventManager.DispatchEvent(&events.DMMessageUpdate{
			GenericDMMessage: &events.GenericDMMessage{
				GenericEvent: genericEvent,
				MessageID:    event.ID,
				Message:      event.Message,
				ChannelID:    event.ChannelID,
			},
			OldMessage: oldMessage,
		})
	} else {
		client.EventManager.DispatchEvent(&events.GuildMessageUpdate{
			GenericGuildMessage: &events.GenericGuildMessage{
				GenericEvent: genericEvent,
				MessageID:    event.ID,
				Message:      event.Message,
				ChannelID:    event.ChannelID,
				GuildID:      *event.GuildID,
			},
			OldMessage: oldMessage,
		})
	}
}

func gatewayHandlerMessageDelete(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventMessageDelete) {
	handleMessageDelete(client, sequenceNumber, shardID, event.ID, event.ChannelID, event.GuildID)
}

func gatewayHandlerMessageDeleteBulk(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventMessageDeleteBulk) {
	for _, messageID := range event.IDs {
		handleMessageDelete(client, sequenceNumber, shardID, messageID, event.ChannelID, event.GuildID)
	}
}

func handleMessageDelete(client *bot.Client, sequenceNumber int, shardID int, messageID snowflake.ID, channelID snowflake.ID, guildID *snowflake.ID) {
	genericEvent := events.NewGenericEvent(client, sequenceNumber, shardID)

	message, err := client.Caches.RemoveMessage(channelID, messageID)
	if err != nil && !errors.Is(err, cache.ErrNotFound) {
		client.Logger.Error("failed to remove message from cache", slog.Any("err", err), slog.String("channel_id", channelID.String()), slog.String("message_id", messageID.String()))
	}

	thread, err := client.Caches.GuildThread(channelID)
	if err != nil && !errors.Is(err, cache.ErrNotFound) {
		client.Logger.Error("failed to get thread from cache", slog.Any("err", err), slog.String("channel_id", channelID.String()))
	}
	if err == nil {
		if thread.MessageCount > 0 {
			thread.MessageCount--
		}
		client.Caches.AddChannel(thread)
	}

	client.EventManager.DispatchEvent(&events.MessageDelete{
		GenericMessage: &events.GenericMessage{
			GenericEvent: genericEvent,
			MessageID:    messageID,
			Message:      message,
			ChannelID:    channelID,
			GuildID:      guildID,
		},
	})

	if guildID == nil {
		client.EventManager.DispatchEvent(&events.DMMessageDelete{
			GenericDMMessage: &events.GenericDMMessage{
				GenericEvent: genericEvent,
				MessageID:    messageID,
				Message:      message,
				ChannelID:    channelID,
			},
		})
	} else {
		client.EventManager.DispatchEvent(&events.GuildMessageDelete{
			GenericGuildMessage: &events.GenericGuildMessage{
				GenericEvent: genericEvent,
				MessageID:    messageID,
				Message:      message,
				ChannelID:    channelID,
				GuildID:      *guildID,
			},
		})
	}
}
