package gatewayhandlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
)

type gatewayHandlerThreadCreate struct{}

func (h *gatewayHandlerThreadCreate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeThreadCreate
}

func (h *gatewayHandlerThreadCreate) New() interface{} {
	return &discord.GatewayEventThreadCreate{}
}

func (h *gatewayHandlerThreadCreate) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*discord.GatewayEventThreadCreate)

	channel := bot.EntityBuilder.CreateChannel(payload.GuildThread, core.CacheStrategyYes)

	bot.EventManager.Dispatch(&events.ThreadCreateEvent{
		GenericThreadEvent: &events.GenericThreadEvent{
			GenericGuildChannelEvent: &events.GenericGuildChannelEvent{
				GenericChannelEvent: &events.GenericChannelEvent{
					GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
					ChannelID:    payload.ID(),
					Channel:      channel,
				},
				GuildID: payload.GuildID(),
				Channel: channel.(core.GuildChannel),
			},
			Thread: channel.(core.GuildThread),
		},
	})
}
