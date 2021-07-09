package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

type threadListSyncPayload struct {
	GuildID       api.Snowflake       `json:"guild_id"`
	ChannelIDs    []api.Snowflake     `json:"channel_ids"`
	Threads       []*api.ChannelImpl  `json:"threads"`
	ThreadMembers []*api.ThreadMember `json:"members"`
}

type ThreadListSyncHandler struct{}

func (h ThreadListSyncHandler) Event() api.GatewayEventType {
	return api.GatewayEventThreadListSync
}

func (h ThreadListSyncHandler) New() interface{} {
	return &threadListSyncPayload{}
}

func (h ThreadListSyncHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
	threadListSyncPayload, ok := i.(*threadListSyncPayload)
	if !ok {
		return
	}

	disgo.Cache().UncacheThreads(threadListSyncPayload.GuildID)
	disgo.Cache().UncacheThreadMembers(threadListSyncPayload.GuildID)

	var threads []api.Thread

	for _, thread := range threadListSyncPayload.Threads {
		threads = append(threads, disgo.EntityBuilder().CreateThread(thread, api.CacheStrategyYes))
	}

	for _, threadMember := range threadListSyncPayload.ThreadMembers {
		disgo.EntityBuilder().CreateThreadMember(threadListSyncPayload.GuildID, threadMember, api.CacheStrategyYes)
	}

	for _, thread := range threads {
		eventManager.Dispatch(&events.ThreadJoinEvent{
			GenericThreadEvent: &events.GenericThreadEvent{
				GenericChannelEvent: &events.GenericChannelEvent{
					GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
					ChannelID:    thread.ID(),
				},
				Thread: thread,
			},
		})
	}
}
