package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

type gatewayHandlerThreadMembersUpdate struct{}

func (h *gatewayHandlerThreadMembersUpdate) EventType() gateway.EventType {
	return gateway.EventTypeThreadMembersUpdate
}

func (h *gatewayHandlerThreadMembersUpdate) New() any {
	return &gateway.EventThreadMembersUpdate{}
}

func (h *gatewayHandlerThreadMembersUpdate) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*gateway.EventThreadMembersUpdate)

	genericEvent := events.NewGenericEvent(client, sequenceNumber, shardID)

	if thread, ok := client.Caches().Channels().GetGuildThread(payload.ID); ok {
		thread.MemberCount = payload.MemberCount
		client.Caches().Channels().Put(thread.ID(), thread)
	}

	for _, addedMember := range payload.AddedMembers {
		addedMember.Member.GuildID = payload.ID
		client.Caches().ThreadMembers().Put(payload.ID, addedMember.UserID, addedMember.ThreadMember)
		client.Caches().Members().Put(payload.GuildID, addedMember.UserID, addedMember.Member)

		if addedMember.Presence != nil {
			client.Caches().Presences().Put(payload.GuildID, addedMember.UserID, *addedMember.Presence)
		}

		client.EventManager().DispatchEvent(&events.ThreadMemberAdd{
			GenericThreadMember: &events.GenericThreadMember{
				GenericEvent:   genericEvent,
				GuildID:        payload.GuildID,
				ThreadID:       payload.ID,
				ThreadMemberID: addedMember.UserID,
				ThreadMember:   addedMember.ThreadMember,
			},
			Member:   addedMember.Member,
			Presence: addedMember.Presence,
		})
	}

	for _, removedMemberID := range payload.RemovedMemberIDs {
		threadMember, _ := client.Caches().ThreadMembers().Remove(payload.ID, removedMemberID)

		client.EventManager().DispatchEvent(&events.ThreadMemberRemove{
			GenericThreadMember: &events.GenericThreadMember{
				GenericEvent:   genericEvent,
				GuildID:        payload.GuildID,
				ThreadID:       payload.ID,
				ThreadMemberID: removedMemberID,
				ThreadMember:   threadMember,
			},
		})
	}

}
