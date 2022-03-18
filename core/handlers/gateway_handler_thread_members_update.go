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

func (h *gatewayHandlerThreadMembersUpdate) New() interface{} {
	return &discord.GatewayEventThreadMembersUpdate{}
}

func (h *gatewayHandlerThreadMembersUpdate) HandleGatewayEvent(bot core.Bot, sequenceNumber discord.GatewaySequence, v interface{}) {
	payload := *v.(*discord.GatewayEventThreadMembersUpdate)

	genericEvent := events.NewGenericEvent(bot, sequenceNumber)

	if channel := bot.Caches.Channels().Get(payload.ID); channel != nil {
		switch thread := channel.(type) {
		case *core.GuildNewsThread:
			thread.MemberCount = payload.MemberCount

		case *core.GuildPublicThread:
			thread.MemberCount = payload.MemberCount

		case *core.GuildPrivateThread:
			thread.MemberCount = payload.MemberCount
		}
	}

	for i := range payload.AddedMembers {
		threadMember := bot.EntityBuilder.CreateThreadMember(payload.AddedMembers[i].ThreadMember, core.CacheStrategyYes)

		bot.EntityBuilder.CreateMember(payload.GuildID, payload.AddedMembers[i].Member, core.CacheStrategyYes)
		if presence := payload.AddedMembers[i].Presence; presence != nil {
			bot.EntityBuilder.CreatePresence(*presence, core.CacheStrategyYes)
		}

		bot.EventManager().Dispatch(&events.ThreadMemberAddEvent{
			GenericThreadMemberEvent: &events.GenericThreadMemberEvent{
				GenericEvent:   genericEvent,
				GuildID:        payload.GuildID,
				ThreadID:       payload.ID,
				ThreadMemberID: payload.AddedMembers[i].UserID,
				ThreadMember:   threadMember,
			},
		})
	}

	for i := range payload.RemovedMemberIDs {
		threadMember := bot.Caches.ThreadMembers().GetCopy(payload.ID, payload.RemovedMemberIDs[i])
		bot.Caches.ThreadMembers().Remove(payload.ID, payload.RemovedMemberIDs[i])

		bot.EventManager().Dispatch(&events.ThreadMemberRemoveEvent{
			GenericThreadMemberEvent: &events.GenericThreadMemberEvent{
				GenericEvent:   genericEvent,
				GuildID:        payload.GuildID,
				ThreadID:       payload.ID,
				ThreadMemberID: payload.RemovedMemberIDs[i],
				ThreadMember:   threadMember,
			},
		})
	}

}
