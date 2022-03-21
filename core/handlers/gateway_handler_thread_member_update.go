package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

type gatewayHandlerThreadMemberUpdate struct{}

func (h *gatewayHandlerThreadMemberUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeThreadMemberUpdate
}

func (h *gatewayHandlerThreadMemberUpdate) New() any {
	return &discord.ThreadMember{}
}

func (h *gatewayHandlerThreadMemberUpdate) HandleGatewayEvent(_ core.Bot, _ discord.GatewaySequence, _ any) {
	// ThreadMembersUpdate kinda handles this already?
}
