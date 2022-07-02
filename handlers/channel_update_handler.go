package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

type gatewayHandlerChannelUpdate struct{}

func (h *gatewayHandlerChannelUpdate) EventType() gateway.EventType {
	return gateway.EventTypeChannelUpdate
}

func (h *gatewayHandlerChannelUpdate) New() any {
	return &discord.UnmarshalChannel{}
}

func (h *gatewayHandlerChannelUpdate) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	channel := v.(*discord.UnmarshalChannel).Channel

	if guildChannel, ok := channel.(discord.GuildChannel); ok {
		oldGuildChannel, _ := client.Caches().Channels().GetGuildChannel(channel.ID())
		client.Caches().Channels().Put(channel.ID(), channel)

		client.EventManager().DispatchEvent(&events.GuildChannelUpdate{
			GenericGuildChannel: &events.GenericGuildChannel{
				GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
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
					client.EventManager().DispatchEvent(&events.ThreadHide{
						GenericThread: &events.GenericThread{
							GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
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

		client.EventManager().DispatchEvent(&events.DMChannelUpdate{
			GenericDMChannel: &events.GenericDMChannel{
				GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
				ChannelID:    channel.ID(),
				Channel:      dmChannel,
			},
			OldChannel: oldDMChannel,
		})
	}
}
