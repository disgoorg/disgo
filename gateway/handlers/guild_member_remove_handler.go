package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
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
func (h *gatewayHandlerGuildMemberRemove) HandleGatewayEvent(client bot.Client, sequenceNumber discord.GatewaySequence, v any) {
	payload := *v.(*discord.GuildMemberRemoveGatewayEvent)

	if guild, ok := client.Caches().Guilds().Get(payload.GuildID); ok {
		guild.MemberCount--
		client.Caches().Guilds().Put(guild.ID, guild)
	}

	member, _ := client.Caches().Members().Remove(payload.GuildID, payload.User.ID)

	client.EventManager().Dispatch(&events.GuildMemberLeaveEvent{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber),
		GuildID:      payload.GuildID,
		User:         payload.User,
		Member:       member,
	})
}
