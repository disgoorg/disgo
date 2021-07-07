package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

type guildMemberRemoveData struct {
	GuildID api.Snowflake `json:"guild_id"`
	User    *api.User     `json:"user"`
}

// GuildMemberRemoveHandler handles api.GuildMemberRemoveGatewayEvent
type GuildMemberRemoveHandler struct{}

// Event returns the raw gateway event Event
func (h *GuildMemberRemoveHandler) Event() api.GatewayEventType {
	return api.GatewayEventGuildMemberRemove
}

// New constructs a new payload receiver for the raw gateway event
func (h *GuildMemberRemoveHandler) New() interface{} {
	return &guildMemberRemoveData{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *GuildMemberRemoveHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
	memberData, ok := i.(*guildMemberRemoveData)
	if !ok {
		return
	}

	guild := disgo.Cache().Guild(memberData.GuildID)
	if guild == nil {
		// todo: replay event later. maybe guild is not cached yet but in a few seconds
		return
	}
	memberData.User = disgo.EntityBuilder().CreateUser(memberData.User, api.CacheStrategyYes)

	member := disgo.Cache().Member(memberData.GuildID, memberData.User.ID)
	disgo.Cache().UncacheMember(memberData.GuildID, memberData.User.ID)

	eventManager.Dispatch(&events.GuildMemberLeaveEvent{
		GenericGuildMemberEvent: &events.GenericGuildMemberEvent{
			GenericGuildEvent: &events.GenericGuildEvent{
				GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
				Guild:        guild,
			},
			Member: member,
		},
	})
}
