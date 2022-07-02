package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

type gatewayHandlerGuildMemberRemove struct{}

func (h *gatewayHandlerGuildMemberRemove) EventType() gateway.EventType {
	return gateway.EventTypeGuildMemberRemove
}

func (h *gatewayHandlerGuildMemberRemove) New() any {
	return &gateway.EventGuildMemberRemove{}
}

func (h *gatewayHandlerGuildMemberRemove) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*gateway.EventGuildMemberRemove)

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
