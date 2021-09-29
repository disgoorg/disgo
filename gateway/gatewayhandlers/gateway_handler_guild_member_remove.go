package gatewayhandlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
)

// gatewayHandlerGuildMemberRemove handles core.GuildMemberRemoveGatewayEvent
type gatewayHandlerGuildMemberRemove struct{}

// EventType returns the core.GatewayGatewayEventType
func (h *gatewayHandlerGuildMemberRemove) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildMemberRemove
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildMemberRemove) New() interface{} {
	return &discord.GuildMemberRemoveGatewayEvent{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildMemberRemove) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	memberData := *v.(*discord.GuildMemberRemoveGatewayEvent)

	bot.EntityBuilder.CreateUser(memberData.User, core.CacheStrategyYes)

	member := bot.Caches.MemberCache().GetCopy(memberData.GuildID, memberData.User.ID)

	bot.Caches.MemberCache().Remove(memberData.GuildID, memberData.User.ID)

	bot.EventManager.Dispatch(&events.GuildMemberLeaveEvent{
		GenericGuildMemberEvent: &events.GenericGuildMemberEvent{
			GenericGuildEvent: &events.GenericGuildEvent{
				GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
				Guild:        bot.Caches.GuildCache().Get(memberData.GuildID),
			},
			Member: member,
		},
	})
}
