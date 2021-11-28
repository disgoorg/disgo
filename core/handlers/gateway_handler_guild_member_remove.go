package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
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

	if guild := bot.Caches.Guilds().Get(payload.GuildID); guild != nil {
		guild.ApproximateMemberCount--
	}

	member := bot.Caches.Members().GetCopy(payload.GuildID, payload.User.ID)

	bot.Caches.Members().Remove(payload.GuildID, payload.User.ID)

	user := bot.EntityBuilder.CreateUser(payload.User, core.CacheStrategyYes)

	bot.EventManager.Dispatch(&events.GuildMemberLeaveEvent{
		GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
		GuildID:      payload.GuildID,
		User:         user,
		Member:       member,
	})
}
