package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

type channelPinsUpdatePayload struct {
	GuildID          *discord.Snowflake `json:"guild_id"`
	ChannelID        discord.Snowflake  `json:"channel_id"`
	LastPinTimestamp *discord.Time      `json:"last_pin_timestamp"`
}

// ChannelPinsUpdateHandler handles api.GatewayEventChannelUpdate
type ChannelPinsUpdateHandler struct{}

// EventType returns the api.GatewayGatewayEventType
func (h *ChannelPinsUpdateHandler) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeChannelPinsUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *ChannelPinsUpdateHandler) New() interface{} {
	return discord.Channel{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *ChannelPinsUpdateHandler) HandleGatewayEvent(disgo core.Disgo, eventManager core.EventManager, sequenceNumber int, v interface{}) {
	payload, ok := v.(channelPinsUpdatePayload)
	if !ok {
		return
	}

	channel := disgo.Caches().ChannelCache().GetMessageChannel(payload.ChannelID).(*core.ChannelImpl)
	oldTime := channel.LastPinTimestamp()
	channel.Channel.LastPinTimestamp = payload.LastPinTimestamp

	genericChannelEvent := &events.GenericChannelEvent{
		GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
		ChannelID:    payload.ChannelID,
		Channel:      channel,
	}

	eventManager.Dispatch(&events.ChannelPinsUpdateEvent{
		GenericChannelEvent: genericChannelEvent,
		GuildID:             payload.GuildID,
		OldLastPinTimestamp: oldTime,
		NewLastPinTimestamp: payload.LastPinTimestamp,
	})

	if payload.GuildID == nil {
		eventManager.Dispatch(&events.DMChannelPinsUpdateEvent{
			GenericDMChannelEvent: &events.GenericDMChannelEvent{
				GenericChannelEvent: genericChannelEvent,
				DMChannel:           channel,
			},
			OldLastPinTimestamp: oldTime,
			NewLastPinTimestamp: payload.LastPinTimestamp,
		})
	} else {
		eventManager.Dispatch(&events.GuildChannelPinsUpdateEvent{
			GenericGuildChannelEvent: &events.GenericGuildChannelEvent{
				GenericChannelEvent: genericChannelEvent,
				GuildID:             *payload.GuildID,
				GuildChannel:        channel,
			},
			OldLastPinTimestamp: oldTime,
			NewLastPinTimestamp: payload.LastPinTimestamp,
		})
	}

}
