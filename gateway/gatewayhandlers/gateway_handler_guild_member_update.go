package gatewayhandlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
)

// gatewayHandlerGuildMemberUpdate handles core.GuildMemberUpdateGatewayEvent
type gatewayHandlerGuildMemberUpdate struct{}

// EventType returns the core.GatewayGatewayEventType
func (h *gatewayHandlerGuildMemberUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildMemberUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildMemberUpdate) New() interface{} {
	return &discord.Member{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildMemberUpdate) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	member := *v.(*discord.Member)

	oldCoreMember := bot.Caches.MemberCache().GetCopy(member.GuildID, member.User.ID)

	bot.EventManager.Dispatch(&events.GuildMemberUpdateEvent{
		GenericGuildMemberEvent: &events.GenericGuildMemberEvent{
			GenericGuildEvent: &events.GenericGuildEvent{
				GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
				Guild:        bot.Caches.GuildCache().Get(member.GuildID),
			},
			Member: bot.EntityBuilder.CreateMember(member.GuildID, member, core.CacheStrategyYes),
		},
		OldMember: oldCoreMember,
	})
}
