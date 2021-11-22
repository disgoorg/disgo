package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	events2 "github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

// gatewayHandlerGuildMemberUpdate handles core.GuildMemberUpdateGatewayEvent
type gatewayHandlerGuildMemberUpdate struct{}

// EventType returns the core.GatewayGatewayEventType
func (h *gatewayHandlerGuildMemberUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildMemberUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildMemberUpdate) New() interface{} {
	return &discord.Member{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildMemberUpdate) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*discord.Member)

	oldCoreMember := bot.Caches.Members().GetCopy(payload.GuildID, payload.User.ID)

	bot.EventManager.Dispatch(&events2.GuildMemberUpdateEvent{
		GenericGuildMemberEvent: &events2.GenericGuildMemberEvent{
			GenericEvent: events2.NewGenericEvent(bot, sequenceNumber),
			GuildID:      payload.GuildID,
			Member:       bot.EntityBuilder.CreateMember(payload.GuildID, payload, core.CacheStrategyYes),
		},
		OldMember: oldCoreMember,
	})
}
