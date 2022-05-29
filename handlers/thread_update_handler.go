package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

type gatewayHandlerThreadUpdate struct{}

func (h *gatewayHandlerThreadUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeThreadUpdate
}

func (h *gatewayHandlerThreadUpdate) New() any {
	return &discord.GuildThread{}
}

func (h *gatewayHandlerThreadUpdate) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	guildThread := *v.(*discord.GuildThread)

	oldGuildThread, _ := client.Caches().Channels().GetGuildThread(guildThread.ID())
	client.Caches().Channels().Put(guildThread.ID(), guildThread)

	client.EventManager().DispatchEvent(&events.ThreadUpdate{
		GenericThread: &events.GenericThread{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			Thread:       guildThread,
			ThreadID:     guildThread.ID(),
			GuildID:      guildThread.GuildID(),
			ParentID:     *guildThread.ParentID(),
		},
		OldThread: oldGuildThread,
	})
}
