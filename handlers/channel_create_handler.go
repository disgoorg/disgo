package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

type gatewayHandlerChannelCreate struct{}

func (h *gatewayHandlerChannelCreate) EventType() gateway.EventType {
	return gateway.EventTypeChannelCreate
}

func (h *gatewayHandlerChannelCreate) New() any {
	return &discord.UnmarshalChannel{}
}

func (h *gatewayHandlerChannelCreate) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	channel := v.(*discord.UnmarshalChannel).Channel
	client.Caches().Channels().Put(channel.ID(), channel)

	if guildChannel, ok := channel.(discord.GuildChannel); ok {
		client.EventManager().DispatchEvent(&events.GuildChannelCreate{
			GenericGuildChannel: &events.GenericGuildChannel{
				GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
				ChannelID:    channel.ID(),
				Channel:      guildChannel,
				GuildID:      guildChannel.GuildID(),
			},
		})
	} else if dmChannel, ok := channel.(discord.DMChannel); ok {
		client.EventManager().DispatchEvent(&events.DMChannelCreate{
			GenericDMChannel: &events.GenericDMChannel{
				GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
				ChannelID:    channel.ID(),
				Channel:      dmChannel,
			},
		})
	}
}
