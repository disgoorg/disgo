package gatewayhandlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
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
		bot.EventManager.Dispatch(&events.GuildChannelCreateEvent{
			GenericGuildChannelEvent: &events.GenericGuildChannelEvent{
				GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
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
		bot.EventManager.Dispatch(&events.DMChannelCreateEvent{
			GenericDMChannelEvent: &events.GenericDMChannelEvent{
				GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
				ChannelID:    channel.ID(),
				Channel:      dmChannel,
			},
		})
	}
}
