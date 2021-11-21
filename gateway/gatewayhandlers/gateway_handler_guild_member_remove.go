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
	payload := *v.(*discord.GuildMemberRemoveGatewayEvent)

	member := bot.Caches.MemberCache().GetCopy(payload.GuildID, payload.User.ID)

	bot.Caches.MemberCache().Remove(payload.GuildID, payload.User.ID)

	user := bot.EntityBuilder.CreateUser(payload.User, core.CacheStrategyYes)

	bot.EventManager.Dispatch(&events.GuildMemberLeaveEvent{
		GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
		GuildID:      payload.GuildID,
		User:         user,
		Member:       member,
	})
}
