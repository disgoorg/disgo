package handlers

import (
	"errors"
	"log/slog"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerGuildRoleCreate(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildRoleCreate) {
	if err := client.Caches.AddRole(event.Role); err != nil {
		client.Logger.Error("failed to add role to cache", slog.Any("err", err), slog.String("guild_id", event.GuildID.String()), slog.String("role_id", event.Role.ID.String()))
	}

	client.EventManager.DispatchEvent(&events.RoleCreate{
		GenericRole: &events.GenericRole{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			GuildID:      event.GuildID,
			RoleID:       event.Role.ID,
			Role:         event.Role,
		},
	})
}

func gatewayHandlerGuildRoleUpdate(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildRoleUpdate) {
	oldRole, err := client.Caches.Role(event.GuildID, event.Role.ID)
	if err != nil && !errors.Is(err, cache.ErrNotFound) {
		client.Logger.Error("failed to get role from cache", slog.Any("err", err), slog.String("guild_id", event.GuildID.String()), slog.String("role_id", event.Role.ID.String()))
	}
	if err := client.Caches.AddRole(event.Role); err != nil {
		client.Logger.Error("failed to add role to cache", slog.Any("err", err), slog.String("guild_id", event.GuildID.String()), slog.String("role_id", event.Role.ID.String()))
	}

	client.EventManager.DispatchEvent(&events.RoleUpdate{
		GenericRole: &events.GenericRole{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			GuildID:      event.GuildID,
			RoleID:       event.Role.ID,
			Role:         event.Role,
		},
		OldRole: oldRole,
	})
}

func gatewayHandlerGuildRoleDelete(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildRoleDelete) {
	role, err := client.Caches.RemoveRole(event.GuildID, event.RoleID)
	if err != nil && !errors.Is(err, cache.ErrNotFound) {
		client.Logger.Error("failed to remove role from cache", slog.Any("err", err), slog.String("guild_id", event.GuildID.String()), slog.String("role_id", event.RoleID.String()))
	}

	client.EventManager.DispatchEvent(&events.RoleDelete{
		GenericRole: &events.GenericRole{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			GuildID:      event.GuildID,
			RoleID:       event.RoleID,
			Role:         role,
		},
	})
}
