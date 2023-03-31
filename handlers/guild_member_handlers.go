package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerGuildMemberAdd(client bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildMemberAdd) {
	if guild, ok := client.Caches().Guild(event.GuildID); ok {
		guild.MemberCount++
		client.Caches().AddGuild(guild)
	}

	client.Caches().AddMember(event.Member)

	client.EventManager().DispatchEvent(&events.GuildMemberJoin{
		GenericGuildMember: &events.GenericGuildMember{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			GuildID:      event.GuildID,
			Member:       event.Member,
		},
	})
}

func gatewayHandlerGuildMemberUpdate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildMemberUpdate) {
	oldMember, _ := client.Caches().Member(event.GuildID, event.User.ID)
	client.Caches().AddMember(event.Member)

	client.EventManager().DispatchEvent(&events.GuildMemberUpdate{
		GenericGuildMember: &events.GenericGuildMember{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			GuildID:      event.GuildID,
			Member:       event.Member,
		},
		OldMember: oldMember,
	})
}

func gatewayHandlerGuildMemberRemove(client bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildMemberRemove) {
	if guild, ok := client.Caches().Guild(event.GuildID); ok {
		guild.MemberCount--
		client.Caches().AddGuild(guild)
	}

	member, _ := client.Caches().RemoveMember(event.GuildID, event.User.ID)

	client.EventManager().DispatchEvent(&events.GuildMemberLeave{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		GuildID:      event.GuildID,
		User:         event.User,
		Member:       member,
	})
}

func gatewayHandlerGuildMembersChunk(client bot.Client, _ int, _ int, event gateway.EventGuildMembersChunk) {
	for i := range event.Members {
		event.Members[i].GuildID = event.GuildID
	}

	if client.MemberChunkingManager() != nil {
		client.MemberChunkingManager().HandleChunk(event)
	}
}
