package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

type threadMembersUpdatePayload struct {
	ThreadID         api.Snowflake       `json:"id"`
	GuildID          api.Snowflake       `json:"guild_id"`
	MemberCount      int                 `json:"member_count"`
	AddedMembers     []*api.ThreadMember `json:"added_members,omitempty"`
	RemovedMemberIDs []api.Snowflake     `json:"removed_member_ids,omitempty"`
}

type ThreadMembersUpdateHandler struct{}

func (h ThreadMembersUpdateHandler) Event() api.GatewayEventType {
	return api.GatewayEventThreadMembersUpdate
}

func (h ThreadMembersUpdateHandler) New() interface{} {
	return &threadMembersUpdatePayload{}
}

func (h ThreadMembersUpdateHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
	payload, ok := i.(*threadMembersUpdatePayload)
	if !ok {
		return
	}

	oldThread := disgo.Cache().Thread(payload.ThreadID)
	var thread api.Thread
	if oldThread != nil {
		oldThread = &*oldThread.(*api.ChannelImpl)

		channelImpl := &*oldThread.(*api.ChannelImpl)
		channelImpl.MemberCount_ = payload.MemberCount
		thread = disgo.EntityBuilder().CreateThread(channelImpl, api.CacheStrategyYes)
	}

	genericChannelEvent := &events.GenericChannelEvent{
		GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
		ChannelID:    payload.ThreadID,
	}

	genericThreadEvent := &events.GenericThreadEvent{
		GenericChannelEvent: genericChannelEvent,
		Thread:              thread,
	}

	eventManager.Dispatch(&events.ThreadUpdateEvent{
		GenericThreadEvent: genericThreadEvent,
		OldThread:          oldThread,
	})

	for _, threadMember := range payload.AddedMembers {
		eventManager.Dispatch(&events.ThreadMemberAddEvent{
			GenericThreadMemberEvent: &events.GenericThreadMemberEvent{
				GenericThreadEvent: genericThreadEvent,
				ThreadMember:       disgo.EntityBuilder().CreateThreadMember(payload.GuildID, threadMember, api.CacheStrategyYes),
			},
		})
	}

	for _, threadMemberID := range payload.RemovedMemberIDs {

		threadMember := disgo.Cache().ThreadMember(payload.GuildID, payload.ThreadID, threadMemberID)

		disgo.Cache().UncacheThreadMember(payload.GuildID, payload.ThreadID, threadMemberID)

		eventManager.Dispatch(&events.ThreadMemberRemoveEvent{
			GenericThreadMemberEvent: &events.GenericThreadMemberEvent{
				GenericThreadEvent: genericThreadEvent,
				ThreadMember:       threadMember,
			},
		})
	}

}
