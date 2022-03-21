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

func (h *gatewayHandlerThreadUpdate) New() any {
	return &discord.GuildThread{}
}

func (h *gatewayHandlerThreadUpdate) HandleGatewayEvent(bot core.Bot, sequenceNumber discord.GatewaySequence, v any) {
	guildThread := *v.(*discord.GuildThread)

	oldGuildThread, _ := bot.Caches().Channels().GetGuildThread(guildThread.ID())
	bot.Caches().Channels().Put(guildThread.ID(), guildThread)

	bot.EventManager().Dispatch(&events.ThreadUpdateEvent{
		GenericThreadEvent: &events.GenericThreadEvent{
			GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
			Thread:       guildThread,
			ThreadID:     guildThread.ID(),
			GuildID:      guildThread.GuildID(),
			ParentID:     *guildThread.ParentID(),
		},
		OldThread: oldGuildThread,
	})
}
