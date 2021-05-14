package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

type ThreadUpdateHandler struct{}

func (h ThreadUpdateHandler) Event() api.GatewayEventType {
	return api.GatewayEventThreadUpdate
}

func (h ThreadUpdateHandler) New() interface{} {
	return &api.ChannelImpl{}
}

func (h ThreadUpdateHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
	channel, ok := i.(*api.ChannelImpl)
	if !ok {
		return
	}

	genericChannelEvent := events.GenericChannelEvent{
		GenericEvent: events.NewEvent(disgo, sequenceNumber),
		ChannelID:    channel.ID(),
	}
	eventManager.Dispatch(genericChannelEvent)

	oldThread := disgo.Cache().Thread(channel.ID())
	if oldThread != nil {
		oldThread = &*oldThread.(*api.ChannelImpl)
	}

	genericThreadEvent := events.GenericThreadEvent{
		GenericChannelEvent: genericChannelEvent,
		Thread:        disgo.EntityBuilder().CreateThread(channel, api.CacheStrategyYes),
	}
	eventManager.Dispatch(genericThreadEvent)

	eventManager.Dispatch(events.ThreadUpdateEvent{
		GenericThreadEvent: genericThreadEvent,
		OldThread:          oldThread,
	})
}