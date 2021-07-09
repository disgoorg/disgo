package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

type ThreadCreateHandler struct{}

func (h ThreadCreateHandler) Event() api.GatewayEventType {
	return api.GatewayEventThreadCreate
}

func (h ThreadCreateHandler) New() interface{} {
	return &api.ChannelImpl{}
}

func (h ThreadCreateHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
	channel, ok := i.(*api.ChannelImpl)
	if !ok {
		return
	}

	genericThreadEvent := &events.GenericThreadEvent{
		GenericChannelEvent: &events.GenericChannelEvent{
			GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
			ChannelID:    channel.ID(),
		},
		Thread: disgo.EntityBuilder().CreateThread(channel, api.CacheStrategyYes),
	}

	eventManager.Dispatch(&events.ThreadCreateEvent{
		GenericThreadEvent: genericThreadEvent,
	})

	eventManager.Dispatch(&events.ThreadJoinEvent{
		GenericThreadEvent: genericThreadEvent,
	})
}
