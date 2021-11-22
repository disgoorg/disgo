package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	events2 "github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
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

	bot.EventManager.Dispatch(&events2.GuildInviteDeleteEvent{
		GenericGuildInviteEvent: &events2.GenericGuildInviteEvent{
			GenericEvent: events2.NewGenericEvent(bot, sequenceNumber),
			GuildID:      *payload.GuildID,
			ChannelID:    payload.ChannelID,
			Code:         payload.Code,
		},
	})
}
