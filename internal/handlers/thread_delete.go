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

	disgo.Cache().UncacheThread(channel.GuildID(), channel.ID())

	eventManager.Dispatch(&events.ThreadDeleteEvent{
		GenericThreadEvent: &events.GenericThreadEvent{
			GenericChannelEvent: &events.GenericChannelEvent{
				GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
				ChannelID:    channel.ID(),
			},
			Thread: disgo.EntityBuilder().CreateThread(channel, api.CacheStrategyNo),
		},
	})
}
