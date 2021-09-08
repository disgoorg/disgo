package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

// GuildMemberUpdateHandler handles api.GuildMemberUpdateGatewayEvent
type GuildMemberUpdateHandler struct{}

// EventType returns the api.GatewayGatewayEventType
func (h *GuildMemberUpdateHandler) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildMemberUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *GuildMemberUpdateHandler) New() interface{} {
	return discord.Member{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *GuildMemberUpdateHandler) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	member, ok := v.(discord.Member)
	if !ok {
		return
	}

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
