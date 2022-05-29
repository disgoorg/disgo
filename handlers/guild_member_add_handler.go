package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

type gatewayHandlerGuildMemberAdd struct{}

func (h *gatewayHandlerGuildMemberAdd) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildMemberAdd
}

func (h *gatewayHandlerGuildMemberAdd) New() any {
	return &discord.Member{}
}

func (h *gatewayHandlerGuildMemberAdd) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	member := *v.(*discord.Member)

	if guild, ok := client.Caches().Guilds().Get(member.GuildID); ok {
		guild.MemberCount++
		client.Caches().Guilds().Put(guild.ID, guild)
	}

	client.Caches().Members().Put(member.GuildID, member.User.ID, member)

	client.EventManager().DispatchEvent(&events.GuildMemberJoin{
		GenericGuildMember: &events.GenericGuildMember{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			GuildID:      member.GuildID,
			Member:       member,
		},
	})
}
