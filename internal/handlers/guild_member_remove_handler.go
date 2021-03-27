package handlers

import (
	"github.com/DiscoOrg/disgo/api"
	"github.com/DiscoOrg/disgo/api/events"
)

type GuildMemberRemoveData struct {
	GuildID api.Snowflake `json:"guild_id"`
	User    api.User      `json:"user"`
}

type GuildMemberRemoveHandler struct{}

// Name returns the raw gateway event name
func (h GuildMemberRemoveHandler) Name() string {
	return api.GuildMemberAddGatewayEvent
}

// New constructs a new payload receiver for the raw gateway event
func (h GuildMemberRemoveHandler) New() interface{} {
	return &GuildMemberRemoveData{}
}

// Handle handles the specific raw gateway event
func (h GuildMemberRemoveHandler) Handle(disgo api.Disgo, eventManager api.EventManager, i interface{}) {
	member, ok := i.(*GuildMemberRemoveData)
	if !ok {
		return
	}

	oldMember := *disgo.Cache().Member(member.GuildID, member.User.ID)
	disgo.Cache().UncacheMember(member.GuildID, member.User.ID)

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

	eventManager.Dispatch(events.GuildMemberLeaveEvent{
		GenericGuildMemberEvent: genericGuildMemberEvent,
		Member:                  oldMember,
	})
}
