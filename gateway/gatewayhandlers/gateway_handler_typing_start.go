package gatewayhandlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
)

// gatewayHandlerTypingStart handles discord.GatewayEventTypeInviteDelete
type gatewayHandlerTypingStart struct{}

// EventType returns the discord.GatewayEventType
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

	user := bot.EntityBuilder.CreateUser(payload.User, core.CacheStrategyYes)

	bot.EventManager.Dispatch(&events.UserTypingEvent{
		GenericUserEvent: &events.GenericUserEvent{
			GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
			UserID:       payload.UserID,
			User:         user,
		},
		ChannelID: payload.ChannelID,
	})

	if payload.GuildID == nil {
		bot.EventManager.Dispatch(&events.DMChannelUserTypingEvent{
			GenericUserEvent: &events.GenericUserEvent{
				GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
				User:         user,
				UserID:       payload.UserID,
			},
			ChannelID: payload.ChannelID,
		})
	} else {
		bot.EventManager.Dispatch(&events.GuildMemberTypingEvent{
			GenericGuildMemberEvent: &events.GenericGuildMemberEvent{
				GenericGuildEvent: &events.GenericGuildEvent{
					GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
					GuildID:      *payload.GuildID,
					Guild:        bot.Caches.GuildCache().Get(*payload.GuildID),
				},
				Member: bot.EntityBuilder.CreateMember(*payload.GuildID, *payload.Member, core.CacheStrategyYes),
			},
		})
	}
}
