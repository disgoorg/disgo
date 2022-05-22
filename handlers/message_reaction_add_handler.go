package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

// gatewayHandlerMessageUpdate handles discord.GatewayEventTypeMessageReactionAdd
type gatewayHandlerMessageReactionAdd struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerMessageReactionAdd) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeMessageReactionAdd
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerMessageReactionAdd) New() any {
	return &discord.GatewayEventMessageReactionAdd{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerMessageReactionAdd) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*discord.GatewayEventMessageReactionAdd)

	genericEvent := events.NewGenericEvent(client, sequenceNumber, shardID)

	if payload.Member != nil {
		client.Caches().Members().Put(*payload.GuildID, payload.UserID, *payload.Member)
	}

	client.EventManager().DispatchEvent(&events.MessageReactionAdd{
		GenericReaction: &events.GenericReaction{
			GenericEvent: genericEvent,
			MessageID:    payload.MessageID,
			ChannelID:    payload.ChannelID,
			GuildID:      payload.GuildID,
			UserID:       payload.UserID,
			Emoji:        payload.Emoji,
		},
		Member: payload.Member,
	})

	if payload.GuildID == nil {
		client.EventManager().DispatchEvent(&events.DMMessageReactionAdd{
			GenericDMMessageReaction: &events.GenericDMMessageReaction{
				GenericEvent: genericEvent,
				MessageID:    payload.MessageID,
				ChannelID:    payload.ChannelID,
				UserID:       payload.UserID,
				Emoji:        payload.Emoji,
			},
		})
	} else {
		client.EventManager().DispatchEvent(&events.GuildMessageReactionAdd{
			GenericGuildMessageReaction: &events.GenericGuildMessageReaction{
				GenericEvent: genericEvent,
				MessageID:    payload.MessageID,
				ChannelID:    payload.ChannelID,
				GuildID:      *payload.GuildID,
				UserID:       payload.UserID,
				Emoji:        payload.Emoji,
			},
			Member: *payload.Member,
		})
	}
}
