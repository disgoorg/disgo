package events

import "github.com/disgoorg/disgo/discord"

type GenericSubscriptionEvent struct {
	*GenericEvent
	discord.Subscription
}

type SubscriptionCreate struct {
	*GenericSubscriptionEvent
}

type SubscriptionUpdate struct {
	*GenericSubscriptionEvent
}

type SubscriptionDelete struct {
	*GenericSubscriptionEvent
}
