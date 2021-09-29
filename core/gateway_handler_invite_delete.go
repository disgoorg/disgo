package core

import (
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
func (h *gatewayHandlerInviteDelete) HandleGatewayEvent(bot *Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*discord.InviteDeleteGatewayEvent)

	bot.EventManager.Dispatch(&GuildInviteDeleteEvent{
		GenericGuildInviteEvent: &GenericGuildInviteEvent{
			GenericGuildEvent: &GenericGuildEvent{
				GenericEvent: NewGenericEvent(bot, sequenceNumber),
				GuildID:      *payload.GuildID,
				Guild:        bot.Caches.GuildCache().Get(*payload.GuildID),
			},
			Code:      payload.Code,
			ChannelID: payload.ChannelID,
		},
	})
}
