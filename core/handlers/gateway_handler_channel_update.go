package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	events2 "github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

// gatewayHandlerChannelUpdate handles core.GatewayEventChannelUpdate
type gatewayHandlerChannelUpdate struct{}

// EventType returns the core.GatewayGatewayEventType
func (h *gatewayHandlerChannelUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeChannelUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerChannelUpdate) New() interface{} {
	return &discord.UnmarshalChannel{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerChannelUpdate) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	channel := v.(*discord.UnmarshalChannel).Channel

	oldChannel := bot.Caches.Channels().GetCopy(channel.ID())

	if ch, ok := channel.(discord.GuildChannel); ok {
		var (
			oldGuildChannel core.GuildChannel
			guildChannel    core.GuildChannel
		)
		if c, ok := oldChannel.(core.GuildChannel); ok {
			oldGuildChannel = c
		}
		if c, ok := bot.EntityBuilder.CreateChannel(channel, core.CacheStrategyNo).(core.GuildChannel); ok {
			guildChannel = c
		}
		bot.Caches.Channels().Remove(channel.ID())
		bot.EventManager.Dispatch(&events2.GuildChannelUpdateEvent{
			GenericGuildChannelEvent: &events2.GenericGuildChannelEvent{
				GenericEvent: events2.NewGenericEvent(bot, sequenceNumber),
				ChannelID:    channel.ID(),
				Channel:      guildChannel,
				GuildID:      ch.GuildID(),
			},
			OldChannel: oldGuildChannel,
		})
	} else {
		var (
			oldDmChannel *core.DMChannel
			dmChannel    *core.DMChannel
		)
		if c, ok := oldChannel.(*core.DMChannel); ok {
			oldDmChannel = c
		}
		if c, ok := bot.EntityBuilder.CreateChannel(channel, core.CacheStrategyYes).(*core.DMChannel); ok {
			dmChannel = c
		}
		bot.EventManager.Dispatch(&events2.DMChannelUpdateEvent{
			GenericDMChannelEvent: &events2.GenericDMChannelEvent{
				GenericEvent: events2.NewGenericEvent(bot, sequenceNumber),
				ChannelID:    channel.ID(),
				Channel:      dmChannel,
			},
			OldChannel: oldDmChannel,
		})
	}
}
