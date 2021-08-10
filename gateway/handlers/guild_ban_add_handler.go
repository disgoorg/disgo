package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/gateway"
)

type guildBanAddPayload struct {
	GuildID discord.Snowflake `json:"guild_id"`
	User    *discord.User    `json:"user"`
}

// GuildBanAddHandler handles api.GatewayEventGuildBanAdd
type GuildBanAddHandler struct{}

// Event returns the api.GatewayEventType
func (h *GuildBanAddHandler) EventType() gateway.EventType {
	return gateway.EventTypeGuildBanAdd
}

// New constructs a new payload receiver for the raw gateway event
func (h *GuildBanAddHandler) New() interface{} {
	return &guildBanAddPayload{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *GuildBanAddHandler) HandleGatewayEvent(disgo core.Disgo, eventManager core.EventManager, sequenceNumber int, i interface{}) {
	payload, ok := i.(*guildBanAddPayload)
	if !ok {
		return
	}

	eventManager.Dispatch(&events.GuildBanEvent{
		GenericGuildEvent: &events.GenericGuildEvent{
			GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
			GuildID:      payload.GuildID,
			Guild:        disgo.Cache().Guild(payload.GuildID),
		},
		User: disgo.EntityBuilder().CreateUser(payload.User, core.CacheStrategyNo),
	})
}
