package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

type gatewayHandlerGuildBanAdd struct{}

func (h *gatewayHandlerGuildBanAdd) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildBanAdd
}

func (h *gatewayHandlerGuildBanAdd) New() any {
	return &discord.GatewayEventGuildBanAdd{}
}

func (h *gatewayHandlerGuildBanAdd) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*discord.GatewayEventGuildBanAdd)

	client.EventManager().DispatchEvent(&events.GuildBan{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		GuildID:      payload.GuildID,
		User:         payload.User,
	})
}
