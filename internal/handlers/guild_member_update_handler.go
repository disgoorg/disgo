package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

// GuildMemberUpdateHandler handles api.GuildMemberUpdateGatewayEvent
type GuildMemberUpdateHandler struct{}

// Event returns the raw gateway event Event
func (h GuildMemberUpdateHandler) Event() api.GatewayEventType {
	return api.GatewayEventGuildMemberUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h GuildMemberUpdateHandler) New() interface{} {
	return &api.Member{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h GuildMemberUpdateHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
	newMember, ok := i.(*api.Member)
	if !ok {
		return
	}

	oldMember := disgo.Cache().Member(newMember.GuildID, newMember.User.ID)
	if oldMember != nil {
		oldMember = &*oldMember
	}
	newMember = disgo.EntityBuilder().CreateMember(newMember.GuildID, newMember, api.CacheStrategyYes)

	genericGuildEvent := events.GenericGuildEvent{
		GenericEvent: events.NewEvent(disgo, sequenceNumber),
		GuildID:      newMember.GuildID,
	}
	eventManager.Dispatch(genericGuildEvent)

	genericGuildMemberEvent := events.GenericGuildMemberEvent{
		GenericGuildEvent: genericGuildEvent,
		UserID:            newMember.User.ID,
	}
	eventManager.Dispatch(genericGuildMemberEvent)

	eventManager.Dispatch(events.GuildMemberUpdateEvent{
		GenericGuildMemberEvent: genericGuildMemberEvent,
		NewMember:               newMember,
		OldMember:               oldMember,
	})
}
