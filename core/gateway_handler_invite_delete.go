package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

type inviteDeletePayload struct {
	ChannelID discord.Snowflake  `json:"channel_id"`
	GuildID   *discord.Snowflake `json:"guild_id"`
	Code      string             `json:"code"`
}

// gatewayHandlerInviteDelete handles core.GatewayEventChannelCreate
type gatewayHandlerInviteDelete struct{}

// EventType returns the core.GatewayGatewayEventType
func (h *gatewayHandlerInviteDelete) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeInviteCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerInviteDelete) New() interface{} {
	return &inviteDeletePayload{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerInviteDelete) HandleGatewayEvent(bot *Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*inviteDeletePayload)

	bot.EventManager.Dispatch(&GuildInviteDeleteEvent{
		GenericGuildInviteEvent: &GenericGuildInviteEvent{
			GenericGuildEvent: &GenericGuildEvent{
				GenericEvent: NewGenericEvent(bot, sequenceNumber),
				GuildID:      *payload.GuildID,
				Guild:        bot.Caches.GuildCache().Get(*payload.GuildID),
			},
			Code:      payload.Code,
			ChannelID: payload.ChannelID,
		},
	})
}
