package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

// gatewayHandlerChannelDelete handles discord.GatewayEventTypeChannelDelete
type gatewayHandlerChannelDelete struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerChannelDelete) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeChannelDelete
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerChannelDelete) New() any {
	return &discord.UnmarshalChannel{}
}

// HandleGatewayEvent handles the specific raw gateway event
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
