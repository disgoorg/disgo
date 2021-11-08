package gatewayhandlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

type gatewayHandlerThreadUpdate struct{}

func (h *gatewayHandlerThreadUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeThreadUpdate
}

func (h *gatewayHandlerThreadUpdate) New() interface{} {
	return &discord.UnmarshalChannel{}
}

func (h *gatewayHandlerThreadUpdate) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	_ = v.(*discord.UnmarshalChannel).Channel

}
