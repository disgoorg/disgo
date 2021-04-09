package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

// GuildMemberAddHandler handles api.GuildMemberAddGatewayEvent
type GuildMemberAddHandler struct{}

// Event returns the raw gateway event Event
func (h GuildMemberAddHandler) Event() api.GatewayEventType {
	return api.GatewayEventGuildMemberAdd
}

// New constructs a new payload receiver for the raw gateway event
func (h GuildMemberAddHandler) New() interface{} {
	return &api.Member{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h GuildMemberAddHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
	member, ok := i.(*api.Member)
	if !ok {
		return
	}

	disgo.Cache().CacheMember(member)

	genericGuildEvent := events.GenericGuildEvent{
		GenericEvent: events.NewEvent(disgo, sequenceNumber),
		GuildID:      member.GuildID,
	}
	eventManager.Dispatch(genericGuildEvent)

	genericGuildMemberEvent := events.GenericGuildMemberEvent{
		GenericGuildEvent: genericGuildEvent,
		UserID:            member.User.ID,
	}
	eventManager.Dispatch(genericGuildMemberEvent)

	eventManager.Dispatch(events.GuildMemberJoinEvent{
		GenericGuildMemberEvent: genericGuildMemberEvent,
		Member:                  member,
	})
}
