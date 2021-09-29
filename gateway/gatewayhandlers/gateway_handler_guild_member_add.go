package gatewayhandlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
)

// gatewayHandlerGuildMemberAdd handles discord.GatewayEventTypeGuildMemberAdd
type gatewayHandlerGuildMemberAdd struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerGuildMemberAdd) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildMemberAdd
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildMemberAdd) New() interface{} {
	return &discord.Member{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildMemberAdd) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	member := *v.(*discord.Member)

	bot.EventManager.Dispatch(&events.GuildMemberJoinEvent{
		GenericGuildMemberEvent: &events.GenericGuildMemberEvent{
			GenericGuildEvent: &events.GenericGuildEvent{
				GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
				Guild:        bot.Caches.GuildCache().Get(member.GuildID),
			},
			Member: bot.EntityBuilder.CreateMember(member.GuildID, member, core.CacheStrategyYes),
		},
	})
}
