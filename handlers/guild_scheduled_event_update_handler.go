package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerGuildScheduledEventUpdate struct {}

func (h *gatewayHandlerGuildScheduledEventUpdate) EventType() gateway.EventType {
	return gateway.EventTypeGuildScheduledEventUpdate
}

func (h *gatewayHandlerGuildScheduledEventUpdate) New() any {
	return &discord.GuildScheduledEvent{}
}

func (h *gatewayHandlerGuildScheduledEventUpdate) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	guildScheduledEvent := *v.(*discord.GuildScheduledEvent)

	oldGuildScheduledEvent, _ := client.Caches().GuildScheduledEvents().Get(guildScheduledEvent.GuildID, guildScheduledEvent.ID)
	client.Caches().GuildScheduledEvents().Put(guildScheduledEvent.GuildID, guildScheduledEvent.ID, guildScheduledEvent)

	client.EventManager().DispatchEvent(&events.GuildScheduledEventUpdate{
		GenericGuildScheduledEvent: &events.GenericGuildScheduledEvent{
			GenericEvent:   events.NewGenericEvent(client, sequenceNumber, shardID),
			GuildScheduled: guildScheduledEvent,
		},
		OldGuildScheduled: oldGuildScheduledEvent,
	})
}
