package handlers

import (
	"time"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerChannelCreate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventChannelCreate) {
	client.Caches().AddChannel(event.GuildChannel)

	client.EventManager().DispatchEvent(&events.GuildChannelCreate{
		GenericGuildChannel: &events.GenericGuildChannel{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			ChannelID:    event.ID(),
			Channel:      event.GuildChannel,
			GuildID:      event.GuildChannel.GuildID(),
		},
	})
}

func gatewayHandlerChannelUpdate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventChannelUpdate) {
	oldGuildChannel, _ := client.Caches().Channel(event.ID())
	client.Caches().AddChannel(event.GuildChannel)

	client.EventManager().DispatchEvent(&events.GuildChannelUpdate{
		GenericGuildChannel: &events.GenericGuildChannel{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			ChannelID:    event.ID(),
			Channel:      event.GuildChannel,
			GuildID:      event.GuildChannel.GuildID(),
		},
		OldChannel: oldGuildChannel,
	})

	if event.Type() == discord.ChannelTypeGuildText || event.Type() == discord.ChannelTypeGuildNews {
		if member, ok := client.Caches().Member(event.GuildChannel.GuildID(), client.ID()); ok &&
			client.Caches().MemberPermissionsInChannel(event.GuildChannel, member).Missing(discord.PermissionViewChannel) {
			for _, guildThread := range client.Caches().GuildThreadsInChannel(event.ID()) {
				client.Caches().RemoveThreadMembersByThreadID(guildThread.ID())
				client.Caches().RemoveChannel(guildThread.ID())
				client.EventManager().DispatchEvent(&events.ThreadHide{
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

func gatewayHandlerChannelDelete(client bot.Client, sequenceNumber int, shardID int, event gateway.EventChannelDelete) {
	client.Caches().RemoveChannel(event.ID())

	client.EventManager().DispatchEvent(&events.GuildChannelDelete{
		GenericGuildChannel: &events.GenericGuildChannel{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			ChannelID:    event.ID(),
			Channel:      event.GuildChannel,
			GuildID:      event.GuildChannel.GuildID(),
		},
	})
}

func gatewayHandlerChannelPinsUpdate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventChannelPinsUpdate) {
	if event.GuildID == nil {
		client.EventManager().DispatchEvent(&events.DMChannelPinsUpdate{
			GenericEvent:        events.NewGenericEvent(client, sequenceNumber, shardID),
			ChannelID:           event.ChannelID,
			NewLastPinTimestamp: event.LastPinTimestamp,
		})
		return
	}

	var oldTime *time.Time
	channel, ok := client.Caches().GuildMessageChannel(event.ChannelID)
	if ok {
		oldTime = channel.LastPinTimestamp()
		client.Caches().AddChannel(discord.ApplyLastPinTimestampToChannel(channel, event.LastPinTimestamp))
	}

	client.EventManager().DispatchEvent(&events.GuildChannelPinsUpdate{
		GenericEvent:        events.NewGenericEvent(client, sequenceNumber, shardID),
		GuildID:             *event.GuildID,
		ChannelID:           event.ChannelID,
		OldLastPinTimestamp: oldTime,
		NewLastPinTimestamp: event.LastPinTimestamp,
	})

}
