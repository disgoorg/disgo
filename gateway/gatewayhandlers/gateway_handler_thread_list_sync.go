package gatewayhandlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

type gatewayHandlerThreadListSync struct{}

func (h *gatewayHandlerThreadListSync) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeThreadListSync
}

func (h *gatewayHandlerThreadListSync) New() interface{} {
	return &discord.GatewayEventThreadListSync{}
}

func (h *gatewayHandlerThreadListSync) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*discord.GatewayEventThreadListSync)

}
