package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
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

		bot.EventManager.Dispatch(&events.GuildChannelUpdateEvent{
			GenericGuildChannelEvent: &events.GenericGuildChannelEvent{
				GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
				ChannelID:    channel.ID(),
				Channel:      guildChannel,
				GuildID:      ch.GuildID(),
			},
			OldChannel: oldGuildChannel,
		})

		if guild := guildChannel.Guild(); guild != nil {
			if guildMessageChannel, ok := guildChannel.(core.GuildMessageChannel); ok && guild.SelfMember().ChannelPermissions(guildChannel).Has(discord.PermissionViewChannel) {
				for _, guildThread := range guildMessageChannel.Threads() {
					bot.Caches.Channels().Remove(guildThread.ID())
					bot.EventManager.Dispatch(&events.ThreadHideEvent{
						GenericThreadEvent: &events.GenericThreadEvent{
							GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
							Thread:       guildThread,
							ThreadID:     guildThread.ID(),
							GuildID:      guildThread.GuildID(),
							ParentID:     guildThread.ParentID(),
						},
					})
				}
			}
		}
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
		bot.EventManager.Dispatch(&events.DMChannelUpdateEvent{
			GenericDMChannelEvent: &events.GenericDMChannelEvent{
				GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
				ChannelID:    channel.ID(),
				Channel:      dmChannel,
			},
			OldChannel: oldDmChannel,
		})
	}
}
