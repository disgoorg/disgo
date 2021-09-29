package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

// gatewayHandlerGuildMemberRemove handles discord.GatewayEventTypeGuildMemberRemove
type gatewayHandlerGuildMemberRemove struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerGuildMemberRemove) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildMemberRemove
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildMemberRemove) New() interface{} {
	return &discord.GuildMemberRemoveGatewayEvent{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildMemberRemove) HandleGatewayEvent(bot *Bot, sequenceNumber int, v interface{}) {
	memberData := *v.(*discord.GuildMemberRemoveGatewayEvent)

	bot.EntityBuilder.CreateUser(memberData.User, CacheStrategyYes)

	member := bot.Caches.MemberCache().GetCopy(memberData.GuildID, memberData.User.ID)

	bot.Caches.MemberCache().Remove(memberData.GuildID, memberData.User.ID)

	bot.EventManager.Dispatch(&GuildMemberLeaveEvent{
		GenericGuildMemberEvent: &GenericGuildMemberEvent{
			GenericGuildEvent: &GenericGuildEvent{
				GenericEvent: NewGenericEvent(bot, sequenceNumber),
				Guild:        bot.Caches.GuildCache().Get(memberData.GuildID),
			},
			Member: member,
		},
	})
}
