package handlers

import (
	"time"

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
func (h *gatewayHandlerChannelPinsUpdate) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*discord.GatewayEventChannelPinsUpdate)

	var oldTime *time.Time

	// update discord.MessageChannel.LastMessageID()
	if channel, ok := client.Caches().Channels().GetMessageChannel(payload.ChannelID); ok {
		oldTime = channel.LastPinTimestamp()
		client.Caches().Channels().Put(payload.ChannelID, discord.ApplyLastPinTimestamp(channel, payload.LastPinTimestamp))
	}

	if payload.GuildID == nil {
		client.EventManager().DispatchEvent(&events.DMChannelPinsUpdateEvent{
			GenericEvent:        events.NewGenericEvent(client, sequenceNumber, shardID),
			ChannelID:           payload.ChannelID,
			OldLastPinTimestamp: oldTime,
			NewLastPinTimestamp: payload.LastPinTimestamp,
		})
	} else {
		client.EventManager().DispatchEvent(&events.GuildChannelPinsUpdateEvent{
			GenericEvent:        events.NewGenericEvent(client, sequenceNumber, shardID),
			GuildID:             *payload.GuildID,
			ChannelID:           payload.ChannelID,
			OldLastPinTimestamp: oldTime,
			NewLastPinTimestamp: payload.LastPinTimestamp,
		})
	}

}
