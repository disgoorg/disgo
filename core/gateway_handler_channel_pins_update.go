package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

type channelPinsUpdatePayload struct {
	GuildID          *discord.Snowflake `json:"guild_id"`
	ChannelID        discord.Snowflake  `json:"channel_id"`
	LastPinTimestamp *discord.Time      `json:"last_pin_timestamp"`
}

// gatewayHandlerChannelPinsUpdate handles core.GatewayEventChannelUpdate
type gatewayHandlerChannelPinsUpdate struct{}

// EventType returns the core.GatewayGatewayEventType
func (h *gatewayHandlerChannelPinsUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeChannelPinsUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerChannelPinsUpdate) New() interface{} {
	return &channelPinsUpdatePayload{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerChannelPinsUpdate) HandleGatewayEvent(bot *Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*channelPinsUpdatePayload)

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
