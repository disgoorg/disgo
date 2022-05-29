package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

type gatewayHandlerUserUpdate struct{}

func (h *gatewayHandlerUserUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeUserUpdate
}

func (h *gatewayHandlerUserUpdate) New() any {
	return &discord.OAuth2User{}
}

func (h *gatewayHandlerUserUpdate) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	user := *v.(*discord.OAuth2User)

	oldUser, _ := client.Caches().GetSelfUser()
	client.Caches().PutSelfUser(user)

	client.EventManager().DispatchEvent(&events.SelfUpdate{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		SelfUser:     user,
		OldSelfUser:  oldUser,
	})

}
