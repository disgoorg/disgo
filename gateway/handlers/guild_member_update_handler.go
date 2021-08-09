package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/gateway"
)

// GuildMemberUpdateHandler handles api.GuildMemberUpdateGatewayEvent
type GuildMemberUpdateHandler struct{}

// Event returns the api.GatewayEventType
func (h *GuildMemberUpdateHandler) EventType() gateway.EventType {
	return gateway.EventTypeGuildMemberUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *GuildMemberUpdateHandler) New() interface{} {
	return &discord.Member{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *GuildMemberUpdateHandler) HandleGatewayEvent(disgo core.Disgo, eventManager core.EventManager, sequenceNumber int, i interface{}) {
	member, ok := i.(*discord.Member)
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
	member = disgo.EntityBuilder().CreateMember(member.GuildID, member, core.CacheStrategyYes)

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
