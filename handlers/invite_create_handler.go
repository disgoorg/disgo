package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/snowflake/v2"
)

type gatewayHandlerInviteCreate struct{}

func (h *gatewayHandlerInviteCreate) EventType() gateway.EventType {
	return gateway.EventTypeInviteCreate
}

func (h *gatewayHandlerInviteCreate) New() any {
	return &discord.Invite{}
}

func (h *gatewayHandlerInviteCreate) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	invite := *v.(*discord.Invite)

	var guildID *snowflake.ID
	if invite.Guild != nil {
		guildID = &invite.Guild.ID
	}

	client.EventManager().DispatchEvent(&events.InviteCreate{
		GenericInvite: &events.GenericInvite{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			GuildID:      guildID,
			Code:         invite.Code,
			ChannelID:    invite.ChannelID,
		},
		Invite: invite,
	})
}
