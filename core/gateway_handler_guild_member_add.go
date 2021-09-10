package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

// gatewayHandlerGuildMemberAdd handles core.GuildMemberAddGatewayEvent
type gatewayHandlerGuildMemberAdd struct{}

// EventType returns the core.GatewayGatewayEventType
func (h *gatewayHandlerGuildMemberAdd) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildMemberAdd
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildMemberAdd) New() interface{} {
	return &discord.Member{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildMemberAdd) HandleGatewayEvent(bot *Bot, sequenceNumber int, v interface{}) {
	member := *v.(*discord.Member)

	bot.EventManager.Dispatch(&GuildMemberJoinEvent{
		GenericGuildMemberEvent: &GenericGuildMemberEvent{
			GenericGuildEvent: &GenericGuildEvent{
				GenericEvent: NewGenericEvent(bot, sequenceNumber),
				Guild:        bot.Caches.GuildCache().Get(member.GuildID),
			},
			Member: bot.EntityBuilder.CreateMember(member.GuildID, member, CacheStrategyYes),
		},
	})
}
