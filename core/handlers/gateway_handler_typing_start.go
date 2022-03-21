package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

// gatewayHandlerTypingStart handles discord.GatewayEventTypeInviteDelete
type gatewayHandlerTypingStart struct{}

// EventType returns the core.GatewayGatewayEventType
func (h *gatewayHandlerTypingStart) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeTypingStart
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerTypingStart) New() any {
	return &discord.TypingStartGatewayEvent{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerTypingStart) HandleGatewayEvent(bot core.Bot, sequenceNumber discord.GatewaySequence, v any) {
	payload := *v.(*discord.TypingStartGatewayEvent)

	bot.EventManager().Dispatch(&events.UserTypingStartEvent{
		GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
		ChannelID:    payload.ChannelID,
		GuildID:      payload.GuildID,
		UserID:       payload.UserID,
		Timestamp:    payload.Timestamp,
	})

	if payload.GuildID == nil {
		bot.EventManager().Dispatch(&events.DMUserTypingStartEvent{
			GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
			ChannelID:    payload.ChannelID,
			UserID:       payload.UserID,
			Timestamp:    payload.Timestamp,
		})
	} else {
		bot.Caches().Members().Put(*payload.GuildID, payload.UserID, *payload.Member)
		bot.EventManager().Dispatch(&events.GuildMemberTypingStartEvent{
			GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
			ChannelID:    payload.ChannelID,
			UserID:       payload.UserID,
			GuildID:      *payload.GuildID,
			Timestamp:    payload.Timestamp,
			Member:       *payload.Member,
		})
	}
}
