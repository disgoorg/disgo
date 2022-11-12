package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerGuildDelete(client bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildDelete) {
	guild, _ := client.Caches().RemoveGuild(event.ID)
	client.Caches().RemoveVoiceStatesByGuildID(event.ID)
	client.Caches().RemovePresencesByGuildID(event.ID)
	// TODO: figure out a better way to remove thread members from cache via guild id without requiring cached GuildThreads
	client.Caches().ChannelsForEach(func(channel discord.GuildChannel) {
		if guildThread, ok := channel.(discord.GuildThread); ok && guildThread.GuildID() == event.ID {
			client.Caches().RemoveThreadMembersByThreadID(guildThread.ID())
		}
	})
	client.Caches().RemoveChannelsByGuildID(event.ID)
	client.Caches().RemoveEmojisByGuildID(event.ID)
	client.Caches().RemoveStickersByGuildID(event.ID)
	client.Caches().RemoveRolesByGuildID(event.ID)
	client.Caches().RemoveStageInstancesByGuildID(event.ID)
	client.Caches().RemoveMessagesByGuildID(event.ID)

	if event.Unavailable {
		client.Caches().SetGuildUnavailable(event.ID, true)
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
