package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

// gatewayHandlerGuildMemberRemove handles discord.GatewayEventTypeGuildMemberRemove
type gatewayHandlerGuildMemberRemove struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerGuildMemberRemove) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildMemberRemove
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildMemberRemove) New() any {
	return &discord.GuildMemberRemoveGatewayEvent{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildMemberRemove) HandleGatewayEvent(bot core.Bot, sequenceNumber discord.GatewaySequence, v any) {
	payload := *v.(*discord.GuildMemberRemoveGatewayEvent)

	if guild, ok := bot.Caches().Guilds().Get(payload.GuildID); ok {
		guild.MemberCount--
		bot.Caches().Guilds().Put(guild.ID, guild)
	}

	member, _ := bot.Caches().Members().Remove(payload.GuildID, payload.User.ID)

	bot.EventManager().Dispatch(&events.GuildMemberLeaveEvent{
		GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
		GuildID:      payload.GuildID,
		User:         payload.User,
		Member:       member,
	})
}
