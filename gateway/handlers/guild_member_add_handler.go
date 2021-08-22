package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
)

// GuildMemberAddHandler handles api.GuildMemberAddGatewayEvent
type GuildMemberAddHandler struct{}

// EventType returns the api.GatewayEventType
func (h *GuildMemberAddHandler) EventType() gateway.EventType {
	return gateway.EventTypeGuildMemberAdd
}

// New constructs a new payload receiver for the raw gateway event
func (h *GuildMemberAddHandler) New() interface{} {
	return discord.Member{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *GuildMemberAddHandler) HandleGatewayEvent(disgo core.Disgo, eventManager core.EventManager, sequenceNumber int, i interface{}) {
	member, ok := i.(discord.Member)
	if !ok {
		return
	}

	eventManager.Dispatch(&events.GuildMemberJoinEvent{
		GenericGuildMemberEvent: &events.GenericGuildMemberEvent{
			GenericGuildEvent: &events.GenericGuildEvent{
				GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
				Guild:        disgo.Cache().GuildCache().Get(member.GuildID),
			},
			Member: disgo.EntityBuilder().CreateMember(member.GuildID, member, core.CacheStrategyYes),
		},
	})
}
