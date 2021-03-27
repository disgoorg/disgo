package handlers

import (
	"github.com/DiscoOrg/disgo/api"
	"github.com/DiscoOrg/disgo/api/events"
)

type GuildMemberUpdateHandler struct{}

// Name returns the raw gateway event name
func (h GuildMemberUpdateHandler) Name() string {
	return api.GuildMemberAddGatewayEvent
}

// New constructs a new payload receiver for the raw gateway event
func (h GuildMemberUpdateHandler) New() interface{} {
	return &api.Member{}
}

// Handle handles the specific raw gateway event
func (h GuildMemberUpdateHandler) Handle(disgo api.Disgo, eventManager api.EventManager, i interface{}) {
	member, ok := i.(*api.Member)
	if !ok {
		return
	}

	oldMember := *disgo.Cache().Member(member.GuildID, member.User.ID)
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

	eventManager.Dispatch(events.GuildMemberUpdateEvent{
		GenericGuildMemberEvent: genericGuildMemberEvent,
		NewMember:               member,
		OldMember:               oldMember,
	})
}
