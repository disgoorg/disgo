package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

type gatewayHandlerGuildScheduledEventDelete struct{}

func (h *gatewayHandlerGuildScheduledEventDelete) EventType() gateway.EventType {
	return gateway.EventTypeGuildScheduledEventDelete
}

func (h *gatewayHandlerGuildScheduledEventDelete) New() any {
	return &discord.GuildScheduledEvent{}
}

func (h *gatewayHandlerGuildScheduledEventDelete) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	guildScheduledEvent := *v.(*discord.GuildScheduledEvent)

	client.Caches().GuildScheduledEvents().Remove(guildScheduledEvent.GuildID, guildScheduledEvent.ID)

	client.EventManager().DispatchEvent(&events.GuildScheduledEventDelete{
		GenericGuildScheduledEvent: &events.GenericGuildScheduledEvent{
			GenericEvent:   events.NewGenericEvent(client, sequenceNumber, shardID),
			GuildScheduled: guildScheduledEvent,
		},
	})
}
