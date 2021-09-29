package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

// gatewayHandlerChannelPinsUpdate handles discord.GatewayEventTypeChannelPinsUpdate
type gatewayHandlerChannelPinsUpdate struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerChannelPinsUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeChannelPinsUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerChannelPinsUpdate) New() interface{} {
	return &discord.ChannelPinsUpdateGatewayEvent{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerChannelPinsUpdate) HandleGatewayEvent(bot *Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*discord.ChannelPinsUpdateGatewayEvent)

	channel := bot.Caches.ChannelCache().Get(payload.ChannelID)
	var oldTime *discord.Time
	if channel != nil {
		oldTime = channel.LastPinTimestamp
		channel.LastPinTimestamp = payload.LastPinTimestamp
	}

	genericChannelEvent := &GenericChannelEvent{
		GenericEvent: NewGenericEvent(bot, sequenceNumber),
		ChannelID:    payload.ChannelID,
		Channel:      channel,
	}

	if payload.GuildID == nil {
		bot.EventManager.Dispatch(&DMChannelPinsUpdateEvent{
			GenericDMChannelEvent: &GenericDMChannelEvent{
				GenericChannelEvent: genericChannelEvent,
			},
			OldLastPinTimestamp: oldTime,
			NewLastPinTimestamp: payload.LastPinTimestamp,
		})
	} else {
		bot.EventManager.Dispatch(&GuildChannelPinsUpdateEvent{
			GenericGuildChannelEvent: &GenericGuildChannelEvent{
				GenericChannelEvent: genericChannelEvent,
				GuildID:             *payload.GuildID,
			},
			OldLastPinTimestamp: oldTime,
			NewLastPinTimestamp: payload.LastPinTimestamp,
		})
	}

}
