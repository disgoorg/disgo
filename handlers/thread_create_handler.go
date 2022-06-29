package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

type gatewayHandlerThreadCreate struct{}

func (h *gatewayHandlerThreadCreate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeThreadCreate
}

func (h *gatewayHandlerThreadCreate) New() any {
	return &discord.GatewayEventThreadCreate{}
}

func (h *gatewayHandlerThreadCreate) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*discord.GatewayEventThreadCreate)

	client.Caches().Channels().Put(payload.ID(), payload.GuildThread)
	client.Caches().ThreadMembers().Put(payload.ID(), payload.ThreadMember.UserID, payload.ThreadMember)

	// update discord.GuildForumChannel.LastThreadID()
	if channel, ok := client.Caches().Channels().GetGuildForumChannel(*payload.ParentID()); ok {
		threadID := payload.ID()
		channel.LastThreadID = &threadID
		client.Caches().Channels().Put(channel.ID(), channel)
	}

	client.EventManager().DispatchEvent(&events.ThreadCreate{
		GenericThread: &events.GenericThread{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			ThreadID:     payload.ID(),
			GuildID:      payload.GuildID(),
			Thread:       payload.GuildThread,
		},
		ThreadMember: payload.ThreadMember,
	})
}
