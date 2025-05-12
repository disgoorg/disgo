package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerGuildScheduledEventCreate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildScheduledEventCreate) {
	client.Caches().AddGuildScheduledEvent(event.GuildScheduledEvent)

	client.EventManager().DispatchEvent(&events.GuildScheduledEventCreate{
		GenericGuildScheduledEvent: &events.GenericGuildScheduledEvent{
			GenericEvent:   events.NewGenericEvent(client, sequenceNumber, shardID),
			GuildScheduled: event.GuildScheduledEvent,
		},
	})
}

func gatewayHandlerGuildScheduledEventUpdate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildScheduledEventUpdate) {
	oldGuildScheduledEvent, _ := client.Caches().GuildScheduledEvent(event.GuildID, event.ID)
	client.Caches().AddGuildScheduledEvent(event.GuildScheduledEvent)

	client.EventManager().DispatchEvent(&events.GuildScheduledEventUpdate{
		GenericGuildScheduledEvent: &events.GenericGuildScheduledEvent{
			GenericEvent:   events.NewGenericEvent(client, sequenceNumber, shardID),
			GuildScheduled: event.GuildScheduledEvent,
		},
		OldGuildScheduled: oldGuildScheduledEvent,
	})
}

func gatewayHandlerGuildScheduledEventDelete(client bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildScheduledEventDelete) {
	client.Caches().RemoveGuildScheduledEvent(event.GuildID, event.ID)

	client.EventManager().DispatchEvent(&events.GuildScheduledEventDelete{
		GenericGuildScheduledEvent: &events.GenericGuildScheduledEvent{
			GenericEvent:   events.NewGenericEvent(client, sequenceNumber, shardID),
			GuildScheduled: event.GuildScheduledEvent,
		},
	})
}

func gatewayHandlerGuildScheduledEventUserAdd(client bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildScheduledEventUserAdd) {
	client.EventManager().DispatchEvent(&events.GuildScheduledEventUserAdd{
		GenericGuildScheduledEventUser: &events.GenericGuildScheduledEventUser{
			GenericEvent:          events.NewGenericEvent(client, sequenceNumber, shardID),
			GuildScheduledEventID: event.GuildScheduledEventID,
			UserID:                event.UserID,
			GuildID:               event.GuildID,
		},
	})
}

func gatewayHandlerGuildScheduledEventUserRemove(client bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildScheduledEventUserRemove) {
	client.EventManager().DispatchEvent(&events.GuildScheduledEventUserRemove{
		GenericGuildScheduledEventUser: &events.GenericGuildScheduledEventUser{
			GenericEvent:          events.NewGenericEvent(client, sequenceNumber, shardID),
			GuildScheduledEventID: event.GuildScheduledEventID,
			UserID:                event.UserID,
			GuildID:               event.GuildID,
		},
	})
}
