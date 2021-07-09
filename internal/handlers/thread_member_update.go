package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

type ThreadMemberUpdateHandler struct{}

func (h ThreadMemberUpdateHandler) Event() api.GatewayEventType {
	return api.GatewayEventThreadMemberUpdate
}

func (h ThreadMemberUpdateHandler) New() interface{} {
	return &api.ThreadMember{}
}

func (h ThreadMemberUpdateHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
	threadMember, ok := i.(*api.ThreadMember)
	if !ok {
		return
	}

	var guildID *api.Snowflake
	for threadGuildID, guildThreads := range disgo.Cache().AllThreadCache() {
		if _, ok := guildThreads[threadMember.ThreadID]; ok {
			guildID = &threadGuildID
			break
		}
	}

	if guildID == nil {
		disgo.Logger().Warnf("ThreadMemberUpdate received for uncached thread")
		return
	}

	oldThreadMember := disgo.Cache().ThreadMember(*guildID, threadMember.ThreadID, threadMember.UserID)
	if oldThreadMember != nil {
		oldThreadMember = &*oldThreadMember
	}

	threadMember = disgo.EntityBuilder().CreateThreadMember(*guildID, threadMember, api.CacheStrategyYes)

	eventManager.Dispatch(&events.ThreadMemberUpdateEvent{
		GenericThreadMemberEvent: &events.GenericThreadMemberEvent{
			GenericThreadEvent: &events.GenericThreadEvent{
				GenericChannelEvent: &events.GenericChannelEvent{
					GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
					ChannelID:    threadMember.ThreadID,
				},
				Thread: disgo.Cache().Thread(threadMember.ThreadID),
			},
			ThreadMember: threadMember,
		},
		OldThreadMember: oldThreadMember,
	})
}
