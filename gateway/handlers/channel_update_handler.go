package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
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
func (h *gatewayHandlerChannelUpdate) HandleGatewayEvent(client bot.Client, sequenceNumber discord.GatewaySequence, v any) {
	channel := v.(*discord.UnmarshalChannel).Channel

	if guildChannel, ok := channel.(discord.GuildChannel); ok {
		oldGuildChannel, _ := client.Caches().Channels().GetGuildChannel(channel.ID())
		client.Caches().Channels().Put(channel.ID(), channel)

		client.EventManager().Dispatch(&events.GuildChannelUpdateEvent{
			GenericGuildChannelEvent: &events.GenericGuildChannelEvent{
				GenericEvent: events.NewGenericEvent(client, sequenceNumber),
				ChannelID:    channel.ID(),
				Channel:      guildChannel,
				GuildID:      guildChannel.GuildID(),
			},
			OldChannel: oldGuildChannel,
		})

		if channel.Type() == discord.ChannelTypeGuildText || channel.Type() == discord.ChannelTypeGuildNews {
			if member, ok := client.Caches().Members().Get(guildChannel.GuildID(), client.ID()); ok &&
				client.Caches().GetMemberPermissionsInChannel(guildChannel, member).Missing(discord.PermissionViewChannel) {
				for _, guildThread := range client.Caches().Channels().GuildThreadsInChannel(channel.ID()) {
					client.Caches().ThreadMembers().RemoveAll(guildThread.ID())
					client.Caches().Channels().Remove(guildThread.ID())
					client.EventManager().Dispatch(&events.ThreadHideEvent{
						GenericThreadEvent: &events.GenericThreadEvent{
							GenericEvent: events.NewGenericEvent(client, sequenceNumber),
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
		oldDMChannel, _ := client.Caches().Channels().GetDMChannel(channel.ID())
		client.Caches().Channels().Put(channel.ID(), channel)

		client.EventManager().Dispatch(&events.DMChannelUpdateEvent{
			GenericDMChannelEvent: &events.GenericDMChannelEvent{
				GenericEvent: events.NewGenericEvent(client, sequenceNumber),
				ChannelID:    channel.ID(),
				Channel:      dmChannel,
			},
			OldChannel: oldDMChannel,
		})
	}
}
