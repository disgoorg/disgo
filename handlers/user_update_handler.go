package handlers

import (
	"github.com/snekROmonoro/disgo/bot"
	"github.com/snekROmonoro/disgo/events"
	"github.com/snekROmonoro/disgo/gateway"
)

func gatewayHandlerUserUpdate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventUserUpdate) {
	oldUser, _ := client.Caches().SelfUser()
	client.Caches().SetSelfUser(event.OAuth2User)

	client.EventManager().DispatchEvent(&events.SelfUpdate{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		SelfUser:     event.OAuth2User,
		OldSelfUser:  oldUser,
	})

}
