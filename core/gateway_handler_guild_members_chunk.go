package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

// gatewayHandlerGuildMemberUpdate handles discord.GatewayEventTypeGuildMembersChunk
type gatewayHandlerGuildMembersChunk struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerGuildMembersChunk) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildMembersChunk
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildMembersChunk) New() interface{} {
	return &discord.GuildMembersChunkGatewayEvent{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildMembersChunk) HandleGatewayEvent(bot *Bot, _ int, v interface{}) {
	payload := *v.(*discord.GuildMembersChunkGatewayEvent)

	if bot.MembersChunkingManager != nil {
		bot.MembersChunkingManager.HandleChunk(payload)
	}
}
