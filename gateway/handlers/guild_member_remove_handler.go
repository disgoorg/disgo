package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

type guildMemberRemoveData struct {
	GuildID discord.Snowflake `json:"guild_id"`
	User    discord.User      `json:"user"`
}

// GuildMemberRemoveHandler handles api.GuildMemberRemoveGatewayEvent
type GuildMemberRemoveHandler struct{}

// EventType returns the api.GatewayGatewayEventType
func (h *GuildMemberRemoveHandler) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildMemberRemove
}

// New constructs a new payload receiver for the raw gateway event
func (h *GuildMemberRemoveHandler) New() interface{} {
	return guildMemberRemoveData{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *GuildMemberRemoveHandler) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	memberData, ok := v.(guildMemberRemoveData)
	if !ok {
		return
	}

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
