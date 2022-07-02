package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerGuildRoleCreate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildRoleCreate) {
	client.Caches().Roles().Put(event.GuildID, event.Role.ID, event.Role)

	client.EventManager().DispatchEvent(&events.RoleCreate{
		GenericRole: &events.GenericRole{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			GuildID:      event.GuildID,
			RoleID:       event.Role.ID,
			Role:         event.Role,
		},
	})
}

func gatewayHandlerGuildRoleUpdate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildRoleUpdate) {
	oldRole, _ := client.Caches().Roles().Get(event.GuildID, event.Role.ID)
	client.Caches().Roles().Put(event.GuildID, event.Role.ID, event.Role)

	client.EventManager().DispatchEvent(&events.RoleUpdate{
		GenericRole: &events.GenericRole{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			GuildID:      event.GuildID,
			RoleID:       event.Role.ID,
			Role:         event.Role,
		},
		OldRole: oldRole,
	})
}

func gatewayHandlerGuildRoleDelete(client bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildRoleDelete) {
	role, _ := client.Caches().Roles().Remove(event.GuildID, event.RoleID)

	client.EventManager().DispatchEvent(&events.RoleDelete{
		GenericRole: &events.GenericRole{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			GuildID:      event.GuildID,
			RoleID:       event.RoleID,
			Role:         role,
		},
	})
}
