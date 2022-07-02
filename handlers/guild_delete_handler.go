package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/snowflake/v2"
)

func gatewayHandlerGuildDelete(client bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildDelete) {
	guild, _ := client.Caches().Guilds().Remove(event.ID)
	client.Caches().VoiceStates().RemoveAll(event.ID)
	client.Caches().Presences().RemoveAll(event.ID)
	client.Caches().ThreadMembers().RemoveIf(func(_ snowflake.ID, threadMember discord.ThreadMember) bool {
		// TODO: figure out a better way to remove thread members from cache via guild id without requiring cached GuildThreads
		if thread, ok := client.Caches().Channels().GetGuildThread(threadMember.ThreadID); ok {
			return thread.GuildID() == event.ID
		}
		return false
	})
	client.Caches().Channels().RemoveIf(func(channel discord.Channel) bool {
		if guildChannel, ok := channel.(discord.GuildChannel); ok {
			return guildChannel.GuildID() == event.ID
		}
		return false
	})
	client.Caches().Emojis().RemoveAll(event.ID)
	client.Caches().Stickers().RemoveAll(event.ID)
	client.Caches().Roles().RemoveAll(event.ID)
	client.Caches().StageInstances().RemoveAll(event.ID)

	client.Caches().Messages().RemoveIf(func(channelID snowflake.ID, message discord.Message) bool {
		return message.GuildID != nil && *message.GuildID == event.ID
	})

	if event.Unavailable {
		client.Caches().Guilds().SetUnavailable(event.ID)
	}

	genericGuildEvent := &events.GenericGuild{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		GuildID:      event.ID,
		Guild:        guild,
	}

	if event.Unavailable {
		client.EventManager().DispatchEvent(&events.GuildUnavailable{
			GenericGuild: genericGuildEvent,
		})
	} else {
		client.EventManager().DispatchEvent(&events.GuildLeave{
			GenericGuild: genericGuildEvent,
		})
	}
}
