package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

// gatewayHandlerGuildMemberUpdate handles discord.GatewayEventTypeGuildMemberUpdate
type gatewayHandlerGuildMemberUpdate struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerGuildMemberUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildMemberUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildMemberUpdate) New() any {
	return &discord.GatewayEventGuildMemberAdd{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildMemberUpdate) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	member := *v.(*discord.GatewayEventGuildMemberAdd)

	oldMember, _ := client.Caches().Members().Get(member.GuildID, member.User.ID)
	client.Caches().Members().Put(member.GuildID, member.User.ID, member.Member)

	client.EventManager().DispatchEvent(&events.GuildMemberUpdateEvent{
		GenericGuildMemberEvent: &events.GenericGuildMemberEvent{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			GuildID:      member.GuildID,
			Member:       member.Member,
		},
		OldMember: oldMember,
	})
}
