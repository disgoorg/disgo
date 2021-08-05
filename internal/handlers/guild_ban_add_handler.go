package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

type guildBanAddPayload struct {
	GuildID api.Snowflake `json:"guild_id"`
	User    *api.User     `json:"user"`
}

// GuildBanAddHandler handles api.GatewayEventGuildBanAdd
type GuildBanAddHandler struct{}

// Event returns the api.GatewayEventType
func (h *GuildBanAddHandler) Event() api.GatewayEventType {
	return api.GatewayEventChannelUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *GuildBanAddHandler) New() interface{} {
	return &guildBanAddPayload{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *GuildBanAddHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
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
		User: disgo.EntityBuilder().CreateUser(payload.User, api.CacheStrategyNo),
	})
}
