package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

type gatewayHandlerThreadMemberUpdate struct{}

func (h *gatewayHandlerThreadMemberUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeThreadMemberUpdate
}

func (h *gatewayHandlerThreadMemberUpdate) New() interface{} {
	return &discord.ThreadMember{}
}

func (h *gatewayHandlerThreadMemberUpdate) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*discord.ThreadMember)

	bot.EntityBuilder.CreateThreadMember(payload, core.CacheStrategyYes)
}
