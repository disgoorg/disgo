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
func (h *gatewayHandlerChannelUpdate) New() any {
	return &discord.UnmarshalChannel{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerChannelUpdate) HandleGatewayEvent(bot core.Bot, sequenceNumber discord.GatewaySequence, v any) {
	channel := v.(*discord.UnmarshalChannel).Channel

	if guildChannel, ok := channel.(discord.GuildChannel); ok {
		oldGuildChannel, _ := bot.Caches().Channels().GetGuildChannel(channel.ID())
		bot.Caches().Channels().Put(channel.ID(), channel)

		bot.EventManager().Dispatch(&events.GuildChannelUpdateEvent{
			GenericGuildChannelEvent: &events.GenericGuildChannelEvent{
				GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
				ChannelID:    channel.ID(),
				Channel:      guildChannel,
				GuildID:      guildChannel.GuildID(),
			},
			OldChannel: oldGuildChannel,
		})

		if channel.Type() == discord.ChannelTypeGuildText || channel.Type() == discord.ChannelTypeGuildNews {
			if member, ok := bot.Caches().Members().Get(guildChannel.GuildID(), bot.ClientID()); ok &&
				bot.Caches().GetMemberPermissionsInChannel(guildChannel, member).Missing(discord.PermissionViewChannel) {
				for _, guildThread := range bot.Caches().Channels().GuildThreadsInChannel(channel.ID()) {
					bot.Caches().ThreadMembers().RemoveAll(guildThread.ID())
					bot.Caches().Channels().Remove(guildThread.ID())
					bot.EventManager().Dispatch(&events.ThreadHideEvent{
						GenericThreadEvent: &events.GenericThreadEvent{
							GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
							Thread:       guildThread,
							ThreadID:     guildThread.ID(),
							GuildID:      guildThread.GuildID(),
							ParentID:     *guildThread.ParentID(),
						},
					})
				}
			}

		}
	} else if dmChannel, ok := channel.(discord.DMChannel); ok {
		oldDMChannel, _ := bot.Caches().Channels().GetDMChannel(channel.ID())
		bot.Caches().Channels().Put(channel.ID(), channel)

		bot.EventManager().Dispatch(&events.DMChannelUpdateEvent{
			GenericDMChannelEvent: &events.GenericDMChannelEvent{
				GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
				ChannelID:    channel.ID(),
				Channel:      dmChannel,
			},
			OldChannel: oldDMChannel,
		})
	}
}
