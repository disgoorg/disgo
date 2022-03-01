package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

type gatewayHandlerThreadUpdate struct{}

func (h *gatewayHandlerThreadUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeThreadUpdate
}

func (h *gatewayHandlerThreadUpdate) New() interface{} {
	return &discord.UnmarshalChannel{}
}

func (h *gatewayHandlerThreadUpdate) HandleGatewayEvent(bot *core.Bot, sequenceNumber discord.GatewaySequence, shardID int, v interface{}) {
	payload := v.(*discord.UnmarshalChannel).Channel

	var oldThread core.GuildThread
	if ch, ok := bot.Caches.Channels().Get(payload.ID()).(core.GuildThread); ok {
		oldThread = ch
	}

	thread := bot.EntityBuilder.CreateChannel(payload, core.CacheStrategyYes).(core.GuildThread)

	bot.EventManager.Dispatch(&events.ThreadUpdateEvent{
		GenericThreadEvent: &events.GenericThreadEvent{
			GenericEvent: events.NewGenericEvent(bot, sequenceNumber, shardID),
			Thread:       thread,
			ThreadID:     thread.ID(),
			GuildID:      thread.GuildID(),
			ParentID:     thread.ParentID(),
		},
		OldThread: oldThread,
	})
}
