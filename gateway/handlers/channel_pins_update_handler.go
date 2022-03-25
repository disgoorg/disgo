package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

// gatewayHandlerChannelPinsUpdate handles discord.GatewayEventTypeChannelPinsUpdate
type gatewayHandlerChannelPinsUpdate struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerChannelPinsUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeChannelPinsUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerChannelPinsUpdate) New() any {
	return &discord.GatewayEventChannelPinsUpdate{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerChannelPinsUpdate) HandleGatewayEvent(client bot.Client, sequenceNumber int, v any) {
	payload := *v.(*discord.GatewayEventChannelPinsUpdate)

	var oldTime *discord.Time
	channel, ok := client.Caches().Channels().GetMessageChannel(payload.ChannelID)
	if ok {
		// TODO: update channels last pinned timestamp
		oldTime = channel.LastPinTimestamp()
	}

	if payload.GuildID == nil {
		client.EventManager().DispatchEvent(&events.DMChannelPinsUpdateEvent{
			GenericEvent:        events.NewGenericEvent(client, sequenceNumber),
			ChannelID:           payload.ChannelID,
			OldLastPinTimestamp: oldTime,
			NewLastPinTimestamp: payload.LastPinTimestamp,
		})
	} else {
		client.EventManager().DispatchEvent(&events.GuildChannelPinsUpdateEvent{
			GenericEvent:        events.NewGenericEvent(client, sequenceNumber),
			GuildID:             *payload.GuildID,
			ChannelID:           payload.ChannelID,
			OldLastPinTimestamp: oldTime,
			NewLastPinTimestamp: payload.LastPinTimestamp,
		})
	}

}
