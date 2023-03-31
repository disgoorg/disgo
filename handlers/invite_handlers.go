package handlers

import (
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerInviteCreate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventInviteCreate) {
	var guildID *snowflake.ID
	if event.Guild != nil {
		guildID = &event.Guild.ID
	}

	client.EventManager().DispatchEvent(&events.InviteCreate{
		GenericInvite: &events.GenericInvite{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			GuildID:      guildID,
			Code:         event.Code,
			ChannelID:    event.ChannelID,
		},
		Invite: event.Invite,
	})
}

func gatewayHandlerInviteDelete(client bot.Client, sequenceNumber int, shardID int, event gateway.EventInviteDelete) {
	client.EventManager().DispatchEvent(&events.InviteDelete{
		GenericInvite: &events.GenericInvite{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			GuildID:      event.GuildID,
			ChannelID:    event.ChannelID,
			Code:         event.Code,
		},
	})
}
