package gatewayhandlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
)

type gatewayHandlerThreadDelete struct{}

func (h *gatewayHandlerThreadDelete) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeThreadDelete
}

func (h *gatewayHandlerThreadDelete) New() interface{} {
	return &discord.GatewayEventThreadDelete{}
}

func (h *gatewayHandlerThreadDelete) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*discord.GatewayEventThreadDelete)

	channel := bot.Caches.ChannelCache().GetCopy(payload.ID)
	var thread core.GuildThread
	if  {

	}
	bot.EventManager.Dispatch(&events.ThreadDeleteEvent{
		GenericThreadEvent: &events.GenericThreadEvent{
			GenericGuildChannelEvent: &events.GenericGuildChannelEvent{
				GenericChannelEvent: &events.GenericChannelEvent{
					GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
					ChannelID:    payload.ID,
					Channel:      channel,
				},
				GuildID: payload.GuildID(),
				Channel: channel,
			},
			Thread: channel.(core.GuildThread),
		},
		ThreadID: payload.ID(),
		ParentID: payload.ParentID,
	})
}
