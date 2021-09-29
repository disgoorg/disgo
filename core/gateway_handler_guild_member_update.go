package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

// gatewayHandlerGuildMemberUpdate handles discord.GatewayEventTypeGuildMemberUpdate
type gatewayHandlerGuildMemberUpdate struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerGuildMemberUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildMemberUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildMemberUpdate) New() interface{} {
	return &discord.Member{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildMemberUpdate) HandleGatewayEvent(bot *Bot, sequenceNumber int, v interface{}) {
	member := *v.(*discord.Member)

	oldCoreMember := bot.Caches.MemberCache().GetCopy(member.GuildID, member.User.ID)

	bot.EventManager.Dispatch(&GuildMemberUpdateEvent{
		GenericGuildMemberEvent: &GenericGuildMemberEvent{
			GenericGuildEvent: &GenericGuildEvent{
				GenericEvent: NewGenericEvent(bot, sequenceNumber),
				Guild:        bot.Caches.GuildCache().Get(member.GuildID),
			},
			Member: bot.EntityBuilder.CreateMember(member.GuildID, member, CacheStrategyYes),
		},
		OldMember: oldCoreMember,
	})
}
