package gatewayhandlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

type gatewayHandlerThreadCreate struct{}

func (h *gatewayHandlerThreadCreate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeThreadCreate
}

func (h *gatewayHandlerThreadCreate) New() interface{} {
	return &discord.Channel{}
}

func (h *gatewayHandlerThreadCreate) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	thread := *v.(*discord.Channel)

}