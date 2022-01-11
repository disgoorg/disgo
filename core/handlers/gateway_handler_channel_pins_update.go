package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
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
func (h *gatewayHandlerChannelPinsUpdate) HandleGatewayEvent(bot core.Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*discord.ChannelPinsUpdateGatewayEvent)

	channel := bot.Caches().Channels().Get(payload.ChannelID)
	var oldTime *discord.Time
	if channel != nil {
		oldTime = core.LastPinTimestamp(channel.(core.MessageChannel))
		switch ch := channel.(type) {
		case *core.GuildTextChannel:
			ch.LastPinTimestamp = payload.LastPinTimestamp

		case *core.DMChannel:
			ch.LastPinTimestamp = payload.LastPinTimestamp

		case *core.GroupDMChannel:
			ch.LastPinTimestamp = payload.LastPinTimestamp

		case *core.GuildNewsChannel:
			ch.LastPinTimestamp = payload.LastPinTimestamp

		case *core.GuildNewsThread:
			ch.LastPinTimestamp = payload.LastPinTimestamp

		case *core.GuildPrivateThread:
			ch.LastPinTimestamp = payload.LastPinTimestamp

		case *core.GuildPublicThread:
			ch.LastPinTimestamp = payload.LastPinTimestamp
		}
	}

	if payload.GuildID == nil {
		bot.EventManager().Dispatch(&events.DMChannelPinsUpdateEvent{
			GenericEvent:        events.NewGenericEvent(bot, sequenceNumber),
			ChannelID:           payload.ChannelID,
			OldLastPinTimestamp: oldTime,
			NewLastPinTimestamp: payload.LastPinTimestamp,
		})
	} else {
		bot.EventManager().Dispatch(&events.GuildChannelPinsUpdateEvent{
			GenericEvent:        events.NewGenericEvent(bot, sequenceNumber),
			GuildID:             *payload.GuildID,
			ChannelID:           payload.ChannelID,
			OldLastPinTimestamp: oldTime,
			NewLastPinTimestamp: payload.LastPinTimestamp,
		})
	}

}
