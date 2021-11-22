package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	events2 "github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

// gatewayHandlerTypingStart handles discord.GatewayEventTypeInviteDelete
type gatewayHandlerTypingStart struct{}

// EventType returns the core.GatewayGatewayEventType
func (h *gatewayHandlerTypingStart) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeTypingStart
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerTypingStart) New() interface{} {
	return &discord.TypingStartGatewayEvent{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerTypingStart) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*discord.TypingStartGatewayEvent)

	bot.EventManager.Dispatch(&events2.UserTypingStartEvent{
		GenericEvent: events2.NewGenericEvent(bot, sequenceNumber),
		UserID:       payload.UserID,
		ChannelID:    payload.ChannelID,
	})

	if payload.GuildID == nil {
		bot.EventManager.Dispatch(&events2.DMChannelUserTypingStartEvent{
			GenericEvent: events2.NewGenericEvent(bot, sequenceNumber),
			UserID:       payload.UserID,
			ChannelID:    payload.ChannelID,
		})
	} else {
		bot.EventManager.Dispatch(&events2.GuildMemberTypingStartEvent{
			GenericEvent: events2.NewGenericEvent(bot, sequenceNumber),
			ChannelID:    payload.ChannelID,
			UserID:       payload.UserID,
			GuildID:      *payload.GuildID,
			Timestamp:    payload.Timestamp,
			Member:       bot.EntityBuilder.CreateMember(*payload.GuildID, *payload.Member, core.CacheStrategyYes),
		})
	}
}
