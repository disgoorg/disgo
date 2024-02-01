package handlers

import (
	"github.com/snekROmonoro/snowflake"

	"github.com/snekROmonoro/disgo/bot"
	"github.com/snekROmonoro/disgo/events"
	"github.com/snekROmonoro/disgo/gateway"
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
