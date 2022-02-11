package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
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
func (h *gatewayHandlerMessageReactionAdd) HandleGatewayEvent(bot *core.Bot, sequenceNumber discord.GatewaySequence, v interface{}) {
	payload := *v.(*discord.GatewayEventMessageReactionAdd)

	genericEvent := events.NewGenericEvent(bot, sequenceNumber)

	var member *core.Member
	if payload.Member != nil {
		member = bot.EntityBuilder.CreateMember(*payload.GuildID, *payload.Member, core.CacheStrategyYes)
	}

	bot.EventManager.Dispatch(&events.MessageReactionAddEvent{
		GenericReactionEvent: &events.GenericReactionEvent{
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
		bot.EventManager.Dispatch(&events.DMMessageReactionAddEvent{
			GenericDMMessageReactionEvent: &events.GenericDMMessageReactionEvent{
				GenericEvent: genericEvent,
				MessageID:    payload.MessageID,
				ChannelID:    payload.ChannelID,
				UserID:       payload.UserID,
				Emoji:        payload.Emoji,
			},
		})
	} else {
		bot.EventManager.Dispatch(&events.GuildMessageReactionAddEvent{
			GenericGuildMessageReactionEvent: &events.GenericGuildMessageReactionEvent{
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
