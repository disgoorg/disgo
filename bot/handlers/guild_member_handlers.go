package handlers

import (
	"log/slog"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerGuildMemberAdd(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildMemberAdd) {
	guild, err := client.Caches.Guild(event.GuildID)
	if err != nil && err != cache.ErrNotFound {
		client.Logger.Error("failed to get guild from cache", slog.Any("err", err), slog.String("guild_id", event.GuildID.String()))
	}
	if err == nil {
		guild.MemberCount++
		client.Caches.AddGuild(guild)
	}

	client.Caches.AddMember(event.Member)

	client.EventManager.DispatchEvent(&events.GuildMemberJoin{
		GenericGuildMember: &events.GenericGuildMember{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			GuildID:      event.GuildID,
			Member:       event.Member,
		},
	})
}

func gatewayHandlerGuildMemberUpdate(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildMemberUpdate) {
	oldMember, err := client.Caches.Member(event.GuildID, event.User.ID)
	if err != nil && err != cache.ErrNotFound {
		client.Logger.Error("failed to get member from cache", slog.Any("err", err), slog.String("guild_id", event.GuildID.String()), slog.String("user_id", event.User.ID.String()))
	}
	client.Caches.AddMember(event.Member)

	client.EventManager.DispatchEvent(&events.GuildMemberUpdate{
		GenericGuildMember: &events.GenericGuildMember{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			GuildID:      event.GuildID,
			Member:       event.Member,
		},
		OldMember: oldMember,
	})
}

func gatewayHandlerGuildMemberRemove(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildMemberRemove) {
	guild, err := client.Caches.Guild(event.GuildID)
	if err != nil && err != cache.ErrNotFound {
		client.Logger.Error("failed to get guild from cache", slog.Any("err", err), slog.String("guild_id", event.GuildID.String()))
	}
	if err == nil {
		guild.MemberCount--
		client.Caches.AddGuild(guild)
	}

	member, err := client.Caches.RemoveMember(event.GuildID, event.User.ID)
	if err != nil && err != cache.ErrNotFound {
		client.Logger.Error("failed to remove member from cache", slog.Any("err", err), slog.String("guild_id", event.GuildID.String()), slog.String("user_id", event.User.ID.String()))
	}

	client.EventManager.DispatchEvent(&events.GuildMemberLeave{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		GuildID:      event.GuildID,
		User:         event.User,
		Member:       member,
	})
}

func gatewayHandlerGuildMembersChunk(client *bot.Client, _ int, _ int, event gateway.EventGuildMembersChunk) {
	for i := range event.Members {
		event.Members[i].GuildID = event.GuildID
	}

	if client.MemberChunkingManager != nil {
		client.MemberChunkingManager.HandleChunk(event)
	}
}
