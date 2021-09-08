package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

// GuildMemberAddHandler handles api.GuildMemberAddGatewayEvent
type GuildMemberAddHandler struct{}

// EventType returns the api.GatewayGatewayEventType
func (h *GuildMemberAddHandler) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildMemberAdd
}

// New constructs a new payload receiver for the raw gateway event
func (h *GuildMemberAddHandler) New() interface{} {
	return &discord.Member{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *GuildMemberAddHandler) HandleGatewayEvent(bot *Bot, sequenceNumber int, v interface{}) {
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
