package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

// gatewayHandlerGuildUpdate handles discord.GatewayEventTypeGuildUpdate
type gatewayHandlerGuildUpdate struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerGuildUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildUpdate) New() any {
	return &discord.Guild{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildUpdate) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	guild := *v.(*discord.Guild)

	oldGuild, _ := client.Caches().Guilds().Get(guild.ID)
	client.Caches().Guilds().Put(guild.ID, guild)

	client.EventManager().DispatchEvent(&events.GuildUpdate{
		GenericGuild: &events.GenericGuild{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			Guild:        guild,
		},
		OldGuild: oldGuild,
	})

}
