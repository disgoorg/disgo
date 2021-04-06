package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

type guildMemberRemoveData struct {
	GuildID api.Snowflake `json:"guild_id"`
	User    api.User      `json:"user"`
}

// GuildMemberRemoveHandler handles api.GuildMemberRemoveGatewayEvent
type GuildMemberRemoveHandler struct{}

// Event returns the raw gateway event Event
func (h GuildMemberRemoveHandler) Event() api.GatewayEvent {
	return api.GatewayEventGuildMemberRemove
}

// New constructs a new payload receiver for the raw gateway event
func (h GuildMemberRemoveHandler) New() interface{} {
	return &guildMemberRemoveData{}
}

// Handle handles the specific raw gateway event
func (h GuildMemberRemoveHandler) Handle(disgo api.Disgo, eventManager api.EventManager, i interface{}) {
	member, ok := i.(*guildMemberRemoveData)
	if !ok {
		return
	}

	oldMember := disgo.Cache().Member(member.GuildID, member.User.ID)
	disgo.Cache().UncacheMember(member.GuildID, member.User.ID)

	genericGuildEvent := events.GenericGuildEvent{
		GenericEvent: api.NewEvent(disgo),
		GuildID:      member.GuildID,
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
