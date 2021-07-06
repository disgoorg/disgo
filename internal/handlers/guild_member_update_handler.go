package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

// GuildMemberUpdateHandler handles api.GuildMemberUpdateGatewayEvent
type GuildMemberUpdateHandler struct{}

// Event returns the raw gateway event Event
func (h *GuildMemberUpdateHandler) Event() api.GatewayEventType {
	return api.GatewayEventGuildMemberUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *GuildMemberUpdateHandler) New() interface{} {
	return &api.Member{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *GuildMemberUpdateHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
	member, ok := i.(*api.Member)
	if !ok {
		return
	}

	guild := disgo.Cache().Guild(member.GuildID)
	if guild == nil {
		// todo: replay event later. maybe guild is not cached yet but in a few seconds
		return
	}

	oldMember := disgo.Cache().Member(member.GuildID, member.User.ID)
	if oldMember != nil {
		oldMember = &*oldMember
	}
	member = disgo.EntityBuilder().CreateMember(member.GuildID, member, api.CacheStrategyYes)

	eventManager.Dispatch(&events.GuildMemberUpdateEvent{
		GenericGuildMemberEvent: &events.GenericGuildMemberEvent{
			GenericGuildEvent: &events.GenericGuildEvent{
				GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
				Guild:        guild,
			},
			Member: member,
		},
		OldMember: oldMember,
	})
}
