package handlers

import (
	"time"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerChannelCreate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventChannelCreate) {
	client.Caches().Channels().Put(event.ID(), event.Channel)

	if guildChannel, ok := event.Channel.(discord.GuildChannel); ok {
		client.EventManager().DispatchEvent(&events.GuildChannelCreate{
			GenericGuildChannel: &events.GenericGuildChannel{
				GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
				ChannelID:    event.ID(),
				Channel:      guildChannel,
				GuildID:      guildChannel.GuildID(),
			},
		})
	} else if dmChannel, ok := event.Channel.(discord.DMChannel); ok {
		client.EventManager().DispatchEvent(&events.DMChannelCreate{
			GenericDMChannel: &events.GenericDMChannel{
				GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
				ChannelID:    event.ID(),
				Channel:      dmChannel,
			},
		})
	}
}

func gatewayHandlerChannelUpdate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventChannelUpdate) {
	if guildChannel, ok := event.Channel.(discord.GuildChannel); ok {
		oldGuildChannel, _ := client.Caches().Channels().GetGuildChannel(event.ID())
		client.Caches().Channels().Put(event.ID(), event.Channel)

		client.EventManager().DispatchEvent(&events.GuildChannelUpdate{
			GenericGuildChannel: &events.GenericGuildChannel{
				GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
				ChannelID:    event.ID(),
				Channel:      guildChannel,
				GuildID:      guildChannel.GuildID(),
			},
			OldChannel: oldGuildChannel,
		})

		if event.Type() == discord.ChannelTypeGuildText || event.Type() == discord.ChannelTypeGuildNews {
			if member, ok := client.Caches().Members().Get(guildChannel.GuildID(), client.ID()); ok &&
				client.Caches().GetMemberPermissionsInChannel(guildChannel, member).Missing(discord.PermissionViewChannel) {
				for _, guildThread := range client.Caches().Channels().GuildThreadsInChannel(event.ID()) {
					client.Caches().ThreadMembers().RemoveAll(guildThread.ID())
					client.Caches().Channels().Remove(guildThread.ID())
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
	} else if dmChannel, ok := event.Channel.(discord.DMChannel); ok {
		oldDMChannel, _ := client.Caches().Channels().GetDMChannel(event.ID())
		client.Caches().Channels().Put(event.ID(), event.Channel)

		client.EventManager().DispatchEvent(&events.DMChannelUpdate{
			GenericDMChannel: &events.GenericDMChannel{
				GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
				ChannelID:    event.ID(),
				Channel:      dmChannel,
			},
			OldChannel: oldDMChannel,
		})
	}
}

func gatewayHandlerChannelDelete(client bot.Client, sequenceNumber int, shardID int, event gateway.EventChannelDelete) {
	client.Caches().Channels().Remove(event.ID())

	if guildChannel, ok := event.Channel.(discord.GuildChannel); ok {
		client.EventManager().DispatchEvent(&events.GuildChannelDelete{
			GenericGuildChannel: &events.GenericGuildChannel{
				GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
				ChannelID:    event.ID(),
				Channel:      guildChannel,
				GuildID:      guildChannel.GuildID(),
			},
		})
	} else if dmChannel, ok := event.Channel.(discord.DMChannel); ok {
		client.EventManager().DispatchEvent(&events.DMChannelDelete{
			GenericDMChannel: &events.GenericDMChannel{
				GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
				ChannelID:    event.ID(),
				Channel:      dmChannel,
			},
		})
	}
}

func gatewayHandlerChannelPinsUpdate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventChannelPinsUpdate) {
	var oldTime *time.Time
	channel, ok := client.Caches().Channels().GetMessageChannel(event.ChannelID)
	if ok {
		// TODO: update channels last pinned timestamp
		oldTime = channel.LastPinTimestamp()
		client.Caches().Channels().Put(event.ChannelID, discord.ApplyLastPinTimestampToChannel(channel, event.LastPinTimestamp))
	}

	if event.GuildID == nil {
		client.EventManager().DispatchEvent(&events.DMChannelPinsUpdate{
			GenericEvent:        events.NewGenericEvent(client, sequenceNumber, shardID),
			ChannelID:           event.ChannelID,
			OldLastPinTimestamp: oldTime,
			NewLastPinTimestamp: event.LastPinTimestamp,
		})
	} else {
		client.EventManager().DispatchEvent(&events.GuildChannelPinsUpdate{
			GenericEvent:        events.NewGenericEvent(client, sequenceNumber, shardID),
			GuildID:             *event.GuildID,
			ChannelID:           event.ChannelID,
			OldLastPinTimestamp: oldTime,
			NewLastPinTimestamp: event.LastPinTimestamp,
		})
	}
}
