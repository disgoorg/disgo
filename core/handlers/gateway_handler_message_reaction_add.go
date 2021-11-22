package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	events2 "github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

// gatewayHandlerMessageUpdate handles discord.GatewayEventTypeMessageReactionAdd
type gatewayHandlerMessageReactionAdd struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerMessageReactionAdd) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeMessageReactionAdd
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerMessageReactionAdd) New() interface{} {
	return &discord.GatewayEventMessageReactionAdd{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerMessageReactionAdd) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*discord.GatewayEventMessageReactionAdd)

	genericEvent := events2.NewGenericEvent(bot, sequenceNumber)

	var member *core.Member
	if payload.Member != nil {
		member = bot.EntityBuilder.CreateMember(*payload.GuildID, *payload.Member, core.CacheStrategyYes)
	}

	bot.EventManager.Dispatch(&events2.MessageReactionAddEvent{
		GenericReactionEvent: &events2.GenericReactionEvent{
			GenericEvent: genericEvent,
			MessageID:    payload.MessageID,
			ChannelID:    payload.ChannelID,
			GuildID:      payload.GuildID,
			UserID:       payload.UserID,
			Emoji:        payload.Emoji,
		},
		Member: member,
	})

	if payload.GuildID == nil {
		bot.EventManager.Dispatch(&events2.DMMessageReactionAddEvent{
			GenericDMMessageReactionEvent: &events2.GenericDMMessageReactionEvent{
				GenericEvent: genericEvent,
				MessageID:    payload.MessageID,
				ChannelID:    payload.ChannelID,
				UserID:       payload.UserID,
				Emoji:        payload.Emoji,
			},
		})
	} else {
		bot.EventManager.Dispatch(&events2.GuildMessageReactionAddEvent{
			GenericGuildMessageReactionEvent: &events2.GenericGuildMessageReactionEvent{
				GenericEvent: genericEvent,
				MessageID:    payload.MessageID,
				ChannelID:    payload.ChannelID,
				GuildID:      *payload.GuildID,
				UserID:       payload.UserID,
				Emoji:        payload.Emoji,
			},
			Member: member,
		})
	}
}
