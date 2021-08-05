package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

type guildBanRemovePayload struct {
	GuildID api.Snowflake `json:"guild_id"`
	User    *api.User     `json:"user"`
}

// GuildBanRemoveHandler handles api.GatewayEventGuildBanRemove
type GuildBanRemoveHandler struct{}

// Event returns the api.GatewayEventType
func (h *GuildBanRemoveHandler) Event() api.GatewayEventType {
	return api.GatewayEventGuildBanRemove
}

// New constructs a new payload receiver for the raw gateway event
func (h *GuildBanRemoveHandler) New() interface{} {
	return &guildBanRemovePayload{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *GuildBanRemoveHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
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
		User: disgo.EntityBuilder().CreateUser(payload.User, api.CacheStrategyNo),
	})
}
