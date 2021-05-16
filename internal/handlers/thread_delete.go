package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

type ThreadDeleteHandler struct{}

func (h ThreadDeleteHandler) Event() api.GatewayEventType {
	return api.GatewayEventThreadDelete
}

func (h ThreadDeleteHandler) New() interface{} {
	return &api.ChannelImpl{}
}

func (h ThreadDeleteHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
	channel, ok := i.(*api.ChannelImpl)
	if !ok {
		return
	}

	genericChannelEvent := events.GenericChannelEvent{
		GenericEvent: events.NewEvent(disgo, sequenceNumber),
		ChannelID:    channel.ID(),
	}
	eventManager.Dispatch(genericChannelEvent)

	disgo.Cache().UncacheThread(channel.GuildID(), channel.ID())
	genericThreadEvent := events.GenericThreadEvent{
		GenericChannelEvent: genericChannelEvent,
		Thread:              disgo.EntityBuilder().CreateThread(channel, api.CacheStrategyNo),
	}
	eventManager.Dispatch(genericThreadEvent)

	eventManager.Dispatch(events.ThreadDeleteEvent{
		GenericThreadEvent: genericThreadEvent,
	})
}
