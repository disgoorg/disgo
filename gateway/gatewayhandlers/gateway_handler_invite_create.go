package gatewayhandlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
)

// gatewayHandlerInviteCreate handles discord.GatewayEventTypeInviteCreate
type gatewayHandlerInviteCreate struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerInviteCreate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeInviteCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerInviteCreate) New() interface{} {
	return &discord.Invite{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerInviteCreate) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	invite := *v.(*discord.Invite)

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
