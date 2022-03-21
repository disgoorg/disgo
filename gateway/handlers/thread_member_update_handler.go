package handlers

import (
	"github.com/DisgoOrg/disgo/bot"
	"github.com/DisgoOrg/disgo/discord"
)

type gatewayHandlerThreadMemberUpdate struct{}

func (h *gatewayHandlerThreadMemberUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeThreadMemberUpdate
}

func (h *gatewayHandlerThreadMemberUpdate) New() any {
	return &discord.ThreadMember{}
}

func (h *gatewayHandlerThreadMemberUpdate) HandleGatewayEvent(_ bot.Client, _ discord.GatewaySequence, _ any) {
	// ThreadMembersUpdate kinda handles this already?
}
