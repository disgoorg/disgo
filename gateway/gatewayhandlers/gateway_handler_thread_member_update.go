package gatewayhandlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

type gatewayHandlerThreadMemberUpdate struct{}

func (h *gatewayHandlerThreadMemberUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeThreadMemberUpdate
}

func (h *gatewayHandlerThreadMemberUpdate) New() interface{} {
	return &discord.GatewayEventThreadMemberUpdate{}
}

func (h *gatewayHandlerThreadMemberUpdate) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	_ = *v.(*discord.GatewayEventThreadMemberUpdate)

}
