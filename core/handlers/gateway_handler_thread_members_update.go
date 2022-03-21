package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

type gatewayHandlerThreadMembersUpdate struct{}

func (h *gatewayHandlerThreadMembersUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeThreadMembersUpdate
}

func (h *gatewayHandlerThreadMembersUpdate) New() any {
	return &discord.GatewayEventThreadMembersUpdate{}
}

func (h *gatewayHandlerThreadMembersUpdate) HandleGatewayEvent(bot core.Bot, sequenceNumber discord.GatewaySequence, v any) {
	payload := *v.(*discord.GatewayEventThreadMembersUpdate)

	genericEvent := events.NewGenericEvent(bot, sequenceNumber)

	if thread, ok := bot.Caches().Channels().GetGuildThread(payload.ID); ok {
		thread.MemberCount = payload.MemberCount
		bot.Caches().Channels().Put(thread.ID(), thread)
	}

	for _, addedMember := range payload.AddedMembers {
		bot.Caches().ThreadMembers().Put(payload.ID, addedMember.UserID, addedMember.ThreadMember)
		bot.Caches().Members().Put(payload.GuildID, addedMember.UserID, addedMember.Member)

		if addedMember.Presence != nil {
			bot.Caches().Presences().Put(payload.GuildID, addedMember.UserID, *addedMember.Presence)
		}

		bot.EventManager().Dispatch(&events.ThreadMemberAddEvent{
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
		threadMember, _ := bot.Caches().ThreadMembers().Remove(payload.ID, removedMemberID)

		bot.EventManager().Dispatch(&events.ThreadMemberRemoveEvent{
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
