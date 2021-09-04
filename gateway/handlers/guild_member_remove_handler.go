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
func (h *GuildMemberRemoveHandler) HandleGatewayEvent(disgo core.Disgo, eventManager core.EventManager, sequenceNumber int, v interface{}) {
	memberData, ok := v.(guildMemberRemoveData)
	if !ok {
		return
	}

	disgo.EntityBuilder().CreateUser(memberData.User, core.CacheStrategyYes)

	member := disgo.Caches().MemberCache().GetCopy(memberData.GuildID, memberData.User.ID)

	disgo.Caches().MemberCache().Uncache(memberData.GuildID, memberData.User.ID)

	eventManager.Dispatch(&events.GuildMemberLeaveEvent{
		GenericGuildMemberEvent: &events.GenericGuildMemberEvent{
			GenericGuildEvent: &events.GenericGuildEvent{
				GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
				Guild:        disgo.Caches().GuildCache().Get(memberData.GuildID),
			},
			Member: member,
		},
	})
}
