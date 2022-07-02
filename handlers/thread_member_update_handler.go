package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
)

type gatewayHandlerThreadMemberUpdate struct{}

func (h *gatewayHandlerThreadMemberUpdate) EventType() gateway.EventType {
	return gateway.EventTypeThreadMemberUpdate
}

func (h *gatewayHandlerThreadMemberUpdate) New() any {
	return &discord.ThreadMember{}
}

func (h *gatewayHandlerThreadMemberUpdate) HandleGatewayEvent(_ bot.Client, _ int, _ int, _ any) {
	// ThreadMembersUpdate kinda handles this already?
}
