package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

// GuildMemberAddHandler handles api.GuildMemberAddGatewayEvent
type GuildMemberAddHandler struct{}

// Name returns the raw gateway event name
func (h GuildMemberAddHandler) Name() string {
	return api.GuildMemberAddGatewayEvent
}

// New constructs a new payload receiver for the raw gateway event
func (h GuildMemberAddHandler) New() interface{} {
	return &api.Member{}
}

// Handle handles the specific raw gateway event
func (h GuildMemberAddHandler) Handle(disgo api.Disgo, eventManager api.EventManager, i interface{}) {
	member, ok := i.(*api.Member)
	if !ok {
		return
	}

	disgo.Cache().CacheMember(member)

	genericGuildEvent := events.GenericGuildEvent{
		Event: api.Event{
			Disgo: disgo,
		},
		GuildID: member.GuildID,
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
