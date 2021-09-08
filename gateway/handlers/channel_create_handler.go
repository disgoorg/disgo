package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

// ChannelCreateHandler handles api.GatewayEventChannelCreate
type ChannelCreateHandler struct{}

// EventType returns the api.GatewayGatewayEventType
func (h *ChannelCreateHandler) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeChannelCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h *ChannelCreateHandler) New() interface{} {
	return discord.Channel{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *ChannelCreateHandler) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	channel, ok := v.(discord.Channel)
	if !ok {
		return
	}

	genericChannelEvent := &events.GenericChannelEvent{
		GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
		ChannelID:    channel.ID,
		Channel:      bot.EntityBuilder.CreateChannel(channel, core.CacheStrategyYes),
	}

	if channel.GuildID != nil {
		bot.EventManager.Dispatch(&events.GuildChannelCreateEvent{
			GenericGuildChannelEvent: &events.GenericGuildChannelEvent{
				GuildID:             *channel.GuildID,
				GenericChannelEvent: genericChannelEvent,
			},
		})
	} else {
		bot.EventManager.Dispatch(&events.DMChannelCreateEvent{
			GenericDMChannelEvent: &events.GenericDMChannelEvent{
				GenericChannelEvent: genericChannelEvent,
			},
		})
	}
}
