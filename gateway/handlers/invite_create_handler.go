package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake"
)

// gatewayHandlerInviteCreate handles discord.GatewayEventTypeInviteCreate
type gatewayHandlerInviteCreate struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerInviteCreate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeInviteCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerInviteCreate) New() any {
	return &discord.Invite{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerInviteCreate) HandleGatewayEvent(client bot.Client, sequenceNumber discord.GatewaySequence, v any) {
	invite := *v.(*discord.Invite)

	var guildID *snowflake.Snowflake
	if invite.Guild != nil {
		guildID = &invite.Guild.ID
	}

	client.EventManager().DispatchEvent(&events.InviteCreateEvent{
		GenericInviteEvent: &events.GenericInviteEvent{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber),
			GuildID:      guildID,
			Code:         invite.Code,
			ChannelID:    invite.ChannelID,
		},
		Invite: invite,
	})
}
