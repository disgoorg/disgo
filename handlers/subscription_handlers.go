package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerSubscriptionCreate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventSubscriptionCreate) {
	client.EventManager().DispatchEvent(&events.SubscriptionCreate{
		GenericSubscriptionEvent: &events.GenericSubscriptionEvent{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			Subscription: event.Subscription,
		},
	})
}

func gatewayHandlerSubscriptionUpdate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventSubscriptionUpdate) {
	client.EventManager().DispatchEvent(&events.SubscriptionUpdate{
		GenericSubscriptionEvent: &events.GenericSubscriptionEvent{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			Subscription: event.Subscription,
		},
	})
}

func gatewayHandlerSubscriptionDelete(client bot.Client, sequenceNumber int, shardID int, event gateway.EventSubscriptionDelete) {
	client.EventManager().DispatchEvent(&events.SubscriptionDelete{
		GenericSubscriptionEvent: &events.GenericSubscriptionEvent{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			Subscription: event.Subscription,
		},
	})
}
