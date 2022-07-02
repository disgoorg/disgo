package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerGuildUpdate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildUpdate) {
	oldGuild, _ := client.Caches().Guilds().Get(event.ID)
	client.Caches().Guilds().Put(event.ID, event.Guild)

	client.EventManager().DispatchEvent(&events.GuildUpdate{
		GenericGuild: &events.GenericGuild{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			Guild:        event.Guild,
		},
		OldGuild: oldGuild,
	})

}
