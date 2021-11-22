package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	events2 "github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

// gatewayHandlerChannelCreate handles core.GatewayEventChannelCreate
type gatewayHandlerChannelCreate struct{}

// EventType returns the core.GatewayGatewayEventType
func (h *gatewayHandlerChannelCreate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeChannelCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerChannelCreate) New() interface{} {
	return &discord.UnmarshalChannel{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerChannelCreate) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	channel := v.(*discord.UnmarshalChannel).Channel

	if ch, ok := channel.(discord.GuildChannel); ok {
		var guildChannel core.GuildChannel
		if c, ok := bot.EntityBuilder.CreateChannel(channel, core.CacheStrategyYes).(core.GuildChannel); ok {
			guildChannel = c
		}
		bot.EventManager.Dispatch(&events2.GuildChannelCreateEvent{
			GenericGuildChannelEvent: &events2.GenericGuildChannelEvent{
				GenericEvent: events2.NewGenericEvent(bot, sequenceNumber),
				ChannelID:    channel.ID(),
				Channel:      guildChannel,
				GuildID:      ch.GuildID(),
			},
		})
	} else {
		var dmChannel *core.DMChannel
		if c, ok := bot.EntityBuilder.CreateChannel(channel, core.CacheStrategyYes).(*core.DMChannel); ok {
			dmChannel = c
		}
		bot.EventManager.Dispatch(&events2.DMChannelCreateEvent{
			GenericDMChannelEvent: &events2.GenericDMChannelEvent{
				GenericEvent: events2.NewGenericEvent(bot, sequenceNumber),
				ChannelID:    channel.ID(),
				Channel:      dmChannel,
			},
		})
	}
}
