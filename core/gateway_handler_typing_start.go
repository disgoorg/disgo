package core

import (
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
func (h *gatewayHandlerTypingStart) HandleGatewayEvent(bot *Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*discord.TypingStartGatewayEvent)

	user := bot.EntityBuilder.CreateUser(payload.User, CacheStrategyYes)

	bot.EventManager.Dispatch(&UserTypingEvent{
		GenericUserEvent: &GenericUserEvent{
			GenericEvent: NewGenericEvent(bot, sequenceNumber),
			UserID:       payload.UserID,
			User:         user,
		},
		ChannelID: payload.ChannelID,
	})

	if payload.GuildID == nil {
		bot.EventManager.Dispatch(&DMChannelUserTypingEvent{
			GenericUserEvent: &GenericUserEvent{
				GenericEvent: NewGenericEvent(bot, sequenceNumber),
				User:         user,
				UserID:       payload.UserID,
			},
			ChannelID: payload.ChannelID,
		})
	} else {
		bot.EventManager.Dispatch(&GuildMemberTypingEvent{
			GenericGuildMemberEvent: &GenericGuildMemberEvent{
				GenericGuildEvent: &GenericGuildEvent{
					GenericEvent: NewGenericEvent(bot, sequenceNumber),
					GuildID:      *payload.GuildID,
					Guild:        bot.Caches.GuildCache().Get(*payload.GuildID),
				},
				Member: bot.EntityBuilder.CreateMember(*payload.GuildID, *payload.Member, CacheStrategyYes),
			},
		})
	}
}
