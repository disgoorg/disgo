package gatewayhandlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

type gatewayHandlerThreadMembersUpdate struct{}

func (h *gatewayHandlerThreadMembersUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeThreadMembersUpdate
}

func (h *gatewayHandlerThreadMembersUpdate) New() interface{} {
	return &discord.GatewayEventThreadMembersUpdate{}
}

func (h *gatewayHandlerThreadMembersUpdate) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*discord.GatewayEventThreadMembersUpdate)

}
