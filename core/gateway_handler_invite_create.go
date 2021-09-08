package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

// InviteCreateHandler handles discord.GatewayEventTypeInviteDelete
type InviteCreateHandler struct{}

// EventType returns the core.GatewayGatewayEventType
func (h *InviteCreateHandler) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeInviteDelete
}

// New constructs a new payload receiver for the raw gateway event
func (h *InviteCreateHandler) New() interface{} {
	return &discord.Invite{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *InviteCreateHandler) HandleGatewayEvent(bot *Bot, sequenceNumber int, v interface{}) {
	invite := *v.(*discord.Invite)

	bot.EventManager.Dispatch(&GuildInviteCreateEvent{
		GenericGuildInviteEvent: &GenericGuildInviteEvent{
			GenericGuildEvent: &GenericGuildEvent{
				GenericEvent: NewGenericEvent(bot, sequenceNumber),
				GuildID:      *invite.GuildID,
				Guild:        bot.Caches.GuildCache().Get(*invite.GuildID),
			},
			Code:      invite.Code,
			ChannelID: invite.ChannelID,
		},
		Invite: bot.EntityBuilder.CreateInvite(invite, CacheStrategyYes),
	})
}
