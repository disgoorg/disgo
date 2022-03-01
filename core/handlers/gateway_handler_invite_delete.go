package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

// gatewayHandlerInviteDelete handles discord.GatewayEventTypeInviteDelete
type gatewayHandlerInviteDelete struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerInviteDelete) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeInviteDelete
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerInviteDelete) New() interface{} {
	return &discord.InviteDeleteGatewayEvent{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerInviteDelete) HandleGatewayEvent(bot *core.Bot, sequenceNumber discord.GatewaySequence, shardID int, v interface{}) {
	payload := *v.(*discord.InviteDeleteGatewayEvent)

	bot.EventManager.Dispatch(&events.GuildInviteDeleteEvent{
		GenericGuildInviteEvent: &events.GenericGuildInviteEvent{
			GenericEvent: events.NewGenericEvent(bot, sequenceNumber, shardID),
			GuildID:      *payload.GuildID,
			ChannelID:    payload.ChannelID,
			Code:         payload.Code,
		},
	})
}
