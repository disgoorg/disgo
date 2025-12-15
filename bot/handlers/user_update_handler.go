package handlers

import (
	"errors"
	"log/slog"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerUserUpdate(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventUserUpdate) {
	oldUser, err := client.Caches.SelfUser()
	if err != nil && !errors.Is(err, cache.ErrNotFound) {
		client.Logger.Error("failed to get self user from cache", slog.Any("err", err))
	}
	if err := client.Caches.SetSelfUser(event.OAuth2User); err != nil {
		client.Logger.Error("failed to set self user to cache", slog.Any("err", err), slog.String("user_id", event.OAuth2User.ID.String()))
	}

	client.EventManager.DispatchEvent(&events.SelfUpdate{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		SelfUser:     event.OAuth2User,
		OldSelfUser:  oldUser,
	})

}
