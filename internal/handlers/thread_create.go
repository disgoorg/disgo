package handlers

import (
	"github.com/DisgoOrg/disgo/api"
)

type ThreadCreateHandler struct{}

func (h ThreadCreateHandler) Event() api.GatewayEventType {
	return api.GatewayEventThreadCreate
}

func (h ThreadCreateHandler) New() interface{} {
	return &api.Thread{}
}

func (h ThreadCreateHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {

}

