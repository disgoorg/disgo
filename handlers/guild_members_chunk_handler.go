package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
)

type gatewayHandlerGuildMembersChunk struct{}

func (h *gatewayHandlerGuildMembersChunk) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildMembersChunk
}

func (h *gatewayHandlerGuildMembersChunk) New() any {
	return &discord.GatewayEventGuildMembersChunk{}
}

func (h *gatewayHandlerGuildMembersChunk) HandleGatewayEvent(client bot.Client, _ int, _ int, v any) {
	payload := *v.(*discord.GatewayEventGuildMembersChunk)

	for i := range payload.Members {
		payload.Members[i].GuildID = payload.GuildID
	}

	if client.MemberChunkingManager() != nil {
		client.MemberChunkingManager().HandleChunk(payload)
	}
}
