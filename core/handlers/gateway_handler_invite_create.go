package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/snowflake"
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
func (h *gatewayHandlerInviteCreate) HandleGatewayEvent(bot core.Bot, sequenceNumber discord.GatewaySequence, v any) {
	invite := *v.(*discord.Invite)

	var guildID *snowflake.Snowflake
	if invite.Guild != nil {
		guildID = &invite.Guild.ID
	}

	bot.EventManager().Dispatch(&events.InviteCreateEvent{
		GenericInviteEvent: &events.GenericInviteEvent{
			GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
			GuildID:      guildID,
			Code:         invite.Code,
			ChannelID:    invite.ChannelID,
		},
		Invite: invite,
	})
}
