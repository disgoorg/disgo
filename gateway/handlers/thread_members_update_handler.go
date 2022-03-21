package handlers

import (
	"github.com/DisgoOrg/disgo/bot"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
)

type gatewayHandlerThreadMembersUpdate struct{}

func (h *gatewayHandlerThreadMembersUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeThreadMembersUpdate
}

func (h *gatewayHandlerThreadMembersUpdate) New() any {
	return &discord.GatewayEventThreadMembersUpdate{}
}

func (h *gatewayHandlerThreadMembersUpdate) HandleGatewayEvent(client bot.Client, sequenceNumber discord.GatewaySequence, v any) {
	payload := *v.(*discord.GatewayEventThreadMembersUpdate)

	genericEvent := events.NewGenericEvent(client, sequenceNumber)

	if thread, ok := client.Caches().Channels().GetGuildThread(payload.ID); ok {
		thread.MemberCount = payload.MemberCount
		client.Caches().Channels().Put(thread.ID(), thread)
	}

	for _, addedMember := range payload.AddedMembers {
		client.Caches().ThreadMembers().Put(payload.ID, addedMember.UserID, addedMember.ThreadMember)
		client.Caches().Members().Put(payload.GuildID, addedMember.UserID, addedMember.Member)

		if addedMember.Presence != nil {
			client.Caches().Presences().Put(payload.GuildID, addedMember.UserID, *addedMember.Presence)
		}

		client.EventManager().Dispatch(&events.ThreadMemberAddEvent{
			GenericThreadMemberEvent: &events.GenericThreadMemberEvent{
				GenericEvent:   genericEvent,
				GuildID:        payload.GuildID,
				ThreadID:       payload.ID,
				ThreadMemberID: addedMember.UserID,
				ThreadMember:   addedMember.ThreadMember,
			},
		})
	}

	for _, removedMemberID := range payload.RemovedMemberIDs {
		threadMember, _ := client.Caches().ThreadMembers().Remove(payload.ID, removedMemberID)

		client.EventManager().Dispatch(&events.ThreadMemberRemoveEvent{
			GenericThreadMemberEvent: &events.GenericThreadMemberEvent{
				GenericEvent:   genericEvent,
				GuildID:        payload.GuildID,
				ThreadID:       payload.ID,
				ThreadMemberID: removedMemberID,
				ThreadMember:   threadMember,
			},
		})
	}

}
