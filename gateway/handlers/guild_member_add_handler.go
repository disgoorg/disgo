package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

// gatewayHandlerGuildMemberAdd handles discord.GatewayEventTypeGuildMemberAdd
type gatewayHandlerGuildMemberAdd struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerGuildMemberAdd) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildMemberAdd
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildMemberAdd) New() any {
	return &discord.GatewayEventGuildMemberAdd{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildMemberAdd) HandleGatewayEvent(client bot.Client, sequenceNumber int, v any) {
	member := *v.(*discord.GatewayEventGuildMemberAdd)

	if guild, ok := client.Caches().Guilds().Get(member.GuildID); ok {
		guild.MemberCount++
		client.Caches().Guilds().Put(guild.ID, guild)
	}

	client.Caches().Members().Put(member.GuildID, member.User.ID, member.Member)

	client.EventManager().DispatchEvent(&events.GuildMemberJoinEvent{
		GenericGuildMemberEvent: &events.GenericGuildMemberEvent{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber),
			GuildID:      member.GuildID,
			Member:       member.Member,
		},
	})
}
