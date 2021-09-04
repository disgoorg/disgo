package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

type guildBanAddPayload struct {
	GuildID discord.Snowflake `json:"guild_id"`
	User    discord.User      `json:"user"`
}

// GuildBanAddHandler handles api.GatewayEventGuildBanAdd
type GuildBanAddHandler struct{}

// EventType returns the api.GatewayGatewayEventType
func (h *GuildBanAddHandler) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildBanAdd
}

// New constructs a new payload receiver for the raw gateway event
func (h *GuildBanAddHandler) New() interface{} {
	return guildBanAddPayload{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *GuildBanAddHandler) HandleGatewayEvent(disgo core.Disgo, eventManager core.EventManager, sequenceNumber int, v interface{}) {
	payload, ok := v.(guildBanAddPayload)
	if !ok {
		return
	}

	eventManager.Dispatch(&events.GuildBanEvent{
		GenericGuildEvent: &events.GenericGuildEvent{
			GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
			GuildID:      payload.GuildID,
			Guild:        disgo.Caches().GuildCache().Get(payload.GuildID),
		},
		User: disgo.EntityBuilder().CreateUser(payload.User, core.CacheStrategyNo),
	})
}
