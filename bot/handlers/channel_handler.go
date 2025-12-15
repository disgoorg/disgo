package handlers

import (
	"errors"
	"log/slog"
	"time"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerChannelCreate(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventChannelCreate) {
	client.Caches.AddChannel(event.GuildChannel)

	client.EventManager.DispatchEvent(&events.GuildChannelCreate{
		GenericGuildChannel: &events.GenericGuildChannel{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			ChannelID:    event.ID(),
			Channel:      event.GuildChannel,
			GuildID:      event.GuildChannel.GuildID(),
		},
	})
}

func gatewayHandlerChannelUpdate(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventChannelUpdate) {
	oldGuildChannel, err := client.Caches.Channel(event.ID())
	if err != nil && !errors.Is(err, cache.ErrNotFound) {
		client.Logger.Error("failed to get channel from cache", slog.Any("err", err), slog.String("channel_id", event.ID().String()))
	}
	client.Caches.AddChannel(event.GuildChannel)

	client.EventManager.DispatchEvent(&events.GuildChannelUpdate{
		GenericGuildChannel: &events.GenericGuildChannel{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			ChannelID:    event.ID(),
			Channel:      event.GuildChannel,
			GuildID:      event.GuildChannel.GuildID(),
		},
		OldChannel: oldGuildChannel,
	})

	if event.Type() == discord.ChannelTypeGuildText || event.Type() == discord.ChannelTypeGuildNews {
		member, err := client.Caches.Member(event.GuildChannel.GuildID(), client.ID())
		if err != nil && !errors.Is(err, cache.ErrNotFound) {
			client.Logger.Error("failed to get member from cache", slog.Any("err", err), slog.String("guild_id", event.GuildChannel.GuildID().String()), slog.String("user_id", client.ID().String()))
		}
		if err == nil {
			permissions, err := client.Caches.MemberPermissionsInChannel(event.GuildChannel, member)
			if err != nil && !errors.Is(err, cache.ErrNotFound) {
				client.Logger.Error("failed to get member permissions from cache", slog.Any("err", err))
			}
			if err == nil && permissions.Missing(discord.PermissionViewChannel) {
				guildThreads, err := client.Caches.GuildThreadsInChannel(event.ID())
				if err != nil && !errors.Is(err, cache.ErrNotFound) {
					client.Logger.Error("failed to get guild threads from cache", slog.Any("err", err))
				}
				if err == nil {
					for _, guildThread := range guildThreads {
						client.Caches.RemoveThreadMembersByThreadID(guildThread.ID())
						client.Caches.RemoveChannel(guildThread.ID())
						client.EventManager.DispatchEvent(&events.ThreadHide{
							GenericThread: &events.GenericThread{
								GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
								Thread:       guildThread,
								ThreadID:     guildThread.ID(),
								GuildID:      guildThread.GuildID(),
								ParentID:     *guildThread.ParentID(),
							},
						})
					}
				}
			}
		}

	}
}

func gatewayHandlerChannelDelete(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventChannelDelete) {
	client.Caches.RemoveChannel(event.ID())

	client.EventManager.DispatchEvent(&events.GuildChannelDelete{
		GenericGuildChannel: &events.GenericGuildChannel{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			ChannelID:    event.ID(),
			Channel:      event.GuildChannel,
			GuildID:      event.GuildChannel.GuildID(),
		},
	})
}

func gatewayHandlerChannelPinsUpdate(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventChannelPinsUpdate) {
	if event.GuildID == nil {
		client.EventManager.DispatchEvent(&events.DMChannelPinsUpdate{
			GenericEvent:        events.NewGenericEvent(client, sequenceNumber, shardID),
			ChannelID:           event.ChannelID,
			NewLastPinTimestamp: event.LastPinTimestamp,
		})
		return
	}

	var oldTime *time.Time
	channel, err := client.Caches.GuildMessageChannel(event.ChannelID)
	if err != nil && !errors.Is(err, cache.ErrNotFound) {
		client.Logger.Error("failed to get channel from cache", slog.Any("err", err), slog.String("channel_id", event.ChannelID.String()))
	}
	if err == nil {
		oldTime = channel.LastPinTimestamp()
		client.Caches.AddChannel(discord.ApplyLastPinTimestampToChannel(channel, event.LastPinTimestamp))
	}

	client.EventManager.DispatchEvent(&events.GuildChannelPinsUpdate{
		GenericEvent:        events.NewGenericEvent(client, sequenceNumber, shardID),
		GuildID:             *event.GuildID,
		ChannelID:           event.ChannelID,
		OldLastPinTimestamp: oldTime,
		NewLastPinTimestamp: event.LastPinTimestamp,
	})

}
