package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/gateway"
)

type gatewayHandlerGuildMembersChunk struct{}

func (h *gatewayHandlerGuildMembersChunk) EventType() gateway.EventType {
	return gateway.EventTypeGuildMembersChunk
}

func (h *gatewayHandlerGuildMembersChunk) New() any {
	return &gateway.EventGuildMembersChunk{}
}

func (h *gatewayHandlerGuildMembersChunk) HandleGatewayEvent(client bot.Client, _ int, _ int, v any) {
	payload := *v.(*gateway.EventGuildMembersChunk)

	for i := range payload.Members {
		payload.Members[i].GuildID = payload.GuildID
	}

	if client.MemberChunkingManager() != nil {
		client.MemberChunkingManager().HandleChunk(payload)
	}
}
