package events

import (
	"github.com/disgoorg/disgo/discord"
)

type GenericEntitlementEvent struct {
	*Event
	*GatewayEvent
	discord.Entitlement
}

type EntitlementCreate struct {
	*Event
	*GatewayEvent
	*WebhookEvent
	discord.Entitlement
}

type EntitlementUpdate struct {
	*GenericEntitlementEvent
}

type EntitlementDelete struct {
	*GenericEntitlementEvent
}
