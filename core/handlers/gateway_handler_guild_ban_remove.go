package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	events2 "github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

// gatewayHandlerGuildBanRemove handles core.GatewayEventGuildBanRemove
type gatewayHandlerGuildBanRemove struct{}

// EventType returns the core.GatewayGatewayEventType
func (h *gatewayHandlerGuildBanRemove) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildBanRemove
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildBanRemove) New() interface{} {
	return &discord.GuildBanRemoveGatewayEvent{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildBanRemove) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*discord.GuildBanRemoveGatewayEvent)

	bot.EventManager.Dispatch(&events2.GuildUnbanEvent{
		GenericEvent: events2.NewGenericEvent(bot, sequenceNumber),
		GuildID:      payload.GuildID,
		User:         bot.EntityBuilder.CreateUser(payload.User, core.CacheStrategyNo),
	})
}
