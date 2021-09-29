package gatewayhandlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
)

// gatewayHandlerInviteDelete handles core.GatewayEventChannelCreate
type gatewayHandlerInviteDelete struct{}

// EventType returns the core.GatewayGatewayEventType
func (h *gatewayHandlerInviteDelete) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeInviteCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerInviteDelete) New() interface{} {
	return &discord.InviteDeleteGatewayEvent{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerInviteDelete) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*discord.InviteDeleteGatewayEvent)

	bot.EventManager.Dispatch(&events.GuildInviteDeleteEvent{
		GenericGuildInviteEvent: &events.GenericGuildInviteEvent{
			GenericGuildEvent: &events.GenericGuildEvent{
				GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
				GuildID:      *payload.GuildID,
				Guild:        bot.Caches.GuildCache().Get(*payload.GuildID),
			},
			Code:      payload.Code,
			ChannelID: payload.ChannelID,
		},
	})
}
