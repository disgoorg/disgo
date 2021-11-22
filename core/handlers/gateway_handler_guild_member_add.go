package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	events2 "github.com/DisgoOrg/disgo/core/events"
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
func (h *gatewayHandlerGuildMemberAdd) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*discord.Member)

	if guild := bot.Caches.Guilds().Get(payload.GuildID); guild != nil {
		guild.ApproximateMemberCount++
	}

	bot.EventManager.Dispatch(&events2.GuildMemberJoinEvent{
		GenericGuildMemberEvent: &events2.GenericGuildMemberEvent{
			GenericEvent: events2.NewGenericEvent(bot, sequenceNumber),
			GuildID:      payload.GuildID,
			Member:       bot.EntityBuilder.CreateMember(payload.GuildID, payload, core.CacheStrategyYes),
		},
	})
}
