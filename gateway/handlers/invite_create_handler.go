package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

// InviteCreateHandler handles discord.GatewayEventTypeInviteDelete
type InviteCreateHandler struct{}

// EventType returns the api.GatewayGatewayEventType
func (h *InviteCreateHandler) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeInviteDelete
}

// New constructs a new payload receiver for the raw gateway event
func (h *InviteCreateHandler) New() interface{} {
	return discord.Invite{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *InviteCreateHandler) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	invite, ok := v.(discord.Invite)
	if !ok {
		return
	}

	bot.EventManager.Dispatch(&events.GuildInviteCreateEvent{
		GenericGuildInviteEvent: &events.GenericGuildInviteEvent{
			GenericGuildEvent: &events.GenericGuildEvent{
				GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
				GuildID:      *invite.GuildID,
				Guild:        bot.Caches.GuildCache().Get(*invite.GuildID),
			},
			Code:      invite.Code,
			ChannelID: invite.ChannelID,
		},
		Invite: bot.EntityBuilder.CreateInvite(invite, core.CacheStrategyYes),
	})
}
