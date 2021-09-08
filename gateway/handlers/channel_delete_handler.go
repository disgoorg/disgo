package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

// ChannelDeleteHandler handles api.GatewayEventChannelDelete
type ChannelDeleteHandler struct{}

// EventType returns the api.GatewayGatewayEventType
func (h *ChannelDeleteHandler) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeChannelDelete
}

// New constructs a new payload receiver for the raw gateway event
func (h *ChannelDeleteHandler) New() interface{} {
	return discord.Channel{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *ChannelDeleteHandler) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	channel, ok := v.(discord.Channel)
	if !ok {
		return
	}

	genericChannelEvent := &events.GenericChannelEvent{
		GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
		ChannelID:    channel.ID,
		Channel:      bot.EntityBuilder.CreateChannel(channel, core.CacheStrategyNo),
	}

	if channel.GuildID != nil {
		bot.EventManager.Dispatch(&events.GuildChannelDeleteEvent{
			GenericGuildChannelEvent: &events.GenericGuildChannelEvent{
				GenericChannelEvent: genericChannelEvent,
				GuildID:             *channel.GuildID,
			},
		})
	} else {
		bot.EventManager.Dispatch(&events.DMChannelDeleteEvent{
			GenericDMChannelEvent: &events.GenericDMChannelEvent{
				GenericChannelEvent: genericChannelEvent,
			},
		})
	}
}
