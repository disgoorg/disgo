package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

type guildBanRemovePayload struct {
	GuildID discord.Snowflake `json:"guild_id"`
	User    discord.User      `json:"user"`
}

// GuildBanRemoveHandler handles api.GatewayEventGuildBanRemove
type GuildBanRemoveHandler struct{}

// EventType returns the api.GatewayGatewayEventType
func (h *GuildBanRemoveHandler) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildBanRemove
}

// New constructs a new payload receiver for the raw gateway event
func (h *GuildBanRemoveHandler) New() interface{} {
	return guildBanRemovePayload{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *GuildBanRemoveHandler) HandleGatewayEvent(disgo core.Disgo, eventManager core.EventManager, sequenceNumber int, v interface{}) {
	payload, ok := v.(guildBanRemovePayload)
	if !ok {
		return
	}

	eventManager.Dispatch(&events.GuildUnbanEvent{
		GenericGuildEvent: &events.GenericGuildEvent{
			GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
			GuildID:      payload.GuildID,
			Guild:        disgo.Caches().GuildCache().Get(payload.GuildID),
		},
		User: disgo.EntityBuilder().CreateUser(payload.User, core.CacheStrategyNo),
	})
}
