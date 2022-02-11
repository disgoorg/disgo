package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
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
func (h *gatewayHandlerGuildMemberUpdate) HandleGatewayEvent(bot *core.Bot, sequenceNumber discord.GatewaySequence, v interface{}) {
	payload := *v.(*discord.Member)

	oldCoreMember := bot.Caches.Members().GetCopy(payload.GuildID, payload.User.ID)

	bot.EventManager.Dispatch(&events.GuildMemberUpdateEvent{
		GenericGuildMemberEvent: &events.GenericGuildMemberEvent{
			GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
			GuildID:      payload.GuildID,
			Member:       bot.EntityBuilder.CreateMember(payload.GuildID, payload, core.CacheStrategyYes),
		},
		OldMember: oldCoreMember,
	})
}
