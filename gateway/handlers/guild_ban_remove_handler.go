package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
)

type guildBanRemovePayload struct {
	GuildID discord.Snowflake `json:"guild_id"`
	User    *discord.User     `json:"user"`
}

// GuildBanRemoveHandler handles api.GatewayEventGuildBanRemove
type GuildBanRemoveHandler struct{}

// Event returns the api.GatewayEventType
func (h *GuildBanRemoveHandler) EventType() gateway.EventType {
	return gateway.EventTypeGuildBanRemove
}

// New constructs a new payload receiver for the raw gateway event
func (h *GuildBanRemoveHandler) New() interface{} {
	return &guildBanRemovePayload{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *GuildBanRemoveHandler) HandleGatewayEvent(disgo core.Disgo, eventManager core.EventManager, sequenceNumber int, i interface{}) {
	payload, ok := i.(*guildBanRemovePayload)
	if !ok {
		return
	}

	eventManager.Dispatch(&events.GuildUnbanEvent{
		GenericGuildEvent: &events.GenericGuildEvent{
			GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
			GuildID:      payload.GuildID,
			Guild:        disgo.Cache().Guild(payload.GuildID),
		},
		User: disgo.EntityBuilder().CreateUser(payload.User, core.CacheStrategyNo),
	})
}
