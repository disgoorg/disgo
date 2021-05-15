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

	genericChannelEvent := events.GenericChannelEvent{
		GenericEvent: events.NewEvent(disgo, sequenceNumber),
		ChannelID:    threadMember.ThreadID,
	}
	eventManager.Dispatch(genericChannelEvent)

	genericThreadEvent := events.GenericThreadEvent{
		GenericChannelEvent: genericChannelEvent,
		Thread:              disgo.Cache().Thread(threadMember.ThreadID),
	}
	eventManager.Dispatch(genericThreadEvent)

	genericThreadMemberEvent := events.GenericThreadMemberEvent{
		GenericThreadEvent: genericThreadEvent,
		ThreadMember:       threadMember,
	}
	eventManager.Dispatch(genericThreadMemberEvent)

	eventManager.Dispatch(events.ThreadMemberUpdateEvent{
		GenericThreadMemberEvent: genericThreadMemberEvent,
		OldThreadMember:          oldThreadMember,
	})
}
