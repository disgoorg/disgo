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

	oldThread := disgo.Cache().Thread(channel.ID())
	if oldThread != nil {
		oldThread = &*oldThread.(*api.ChannelImpl)
	}

	eventManager.Dispatch(&events.ThreadUpdateEvent{
		GenericThreadEvent: &events.GenericThreadEvent{
			GenericChannelEvent: &events.GenericChannelEvent{
				GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
				ChannelID:    channel.ID(),
			},
			Thread: disgo.EntityBuilder().CreateThread(channel, api.CacheStrategyYes),
		},
		OldThread: oldThread,
	})
}
