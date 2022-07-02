package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

type gatewayHandlerChannelDelete struct{}

func (h *gatewayHandlerChannelDelete) EventType() gateway.EventType {
	return gateway.EventTypeChannelDelete
}

func (h *gatewayHandlerChannelDelete) New() any {
	return &discord.UnmarshalChannel{}
}

func (h *gatewayHandlerChannelDelete) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	channel := v.(*discord.UnmarshalChannel).Channel

	client.Caches().Channels().Remove(channel.ID())

	if guildChannel, ok := channel.(discord.GuildChannel); ok {
		client.EventManager().DispatchEvent(&events.GuildChannelDelete{
			GenericGuildChannel: &events.GenericGuildChannel{
				GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
				ChannelID:    channel.ID(),
				Channel:      guildChannel,
				GuildID:      guildChannel.GuildID(),
			},
		})
	} else if dmChannel, ok := channel.(discord.DMChannel); ok {
		client.EventManager().DispatchEvent(&events.DMChannelDelete{
			GenericDMChannel: &events.GenericDMChannel{
				GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
				ChannelID:    channel.ID(),
				Channel:      dmChannel,
			},
		})
	}
}
