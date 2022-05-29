package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

type gatewayHandlerGuildMemberRemove struct{}

func (h *gatewayHandlerGuildMemberRemove) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildMemberRemove
}

func (h *gatewayHandlerGuildMemberRemove) New() any {
	return &discord.GatewayEventGuildMemberRemove{}
}

func (h *gatewayHandlerGuildMemberRemove) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*discord.GatewayEventGuildMemberRemove)

	if guild, ok := client.Caches().Guilds().Get(payload.GuildID); ok {
		guild.MemberCount--
		client.Caches().Guilds().Put(guild.ID, guild)
	}

	member, _ := client.Caches().Members().Remove(payload.GuildID, payload.User.ID)

	client.EventManager().DispatchEvent(&events.GuildMemberLeave{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		GuildID:      payload.GuildID,
		User:         payload.User,
		Member:       member,
	})
}
