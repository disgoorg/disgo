package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
)

// gatewayHandlerGuildMemberUpdate handles discord.GatewayEventTypeGuildMembersChunk
type gatewayHandlerGuildMembersChunk struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerGuildMembersChunk) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildMembersChunk
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildMembersChunk) New() any {
	return &discord.GatewayEventGuildMembersChunk{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildMembersChunk) HandleGatewayEvent(client bot.Client, _ int, _ int, v any) {
	payload := *v.(*discord.GatewayEventGuildMembersChunk)

	for i := range payload.Members {
		payload.Members[i].GuildID = payload.GuildID
	}

	if client.MemberChunkingManager() != nil {
		client.MemberChunkingManager().HandleChunk(payload)
	}
}
