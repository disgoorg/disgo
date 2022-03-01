package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

// gatewayHandlerVoiceServerUpdate handles discord.GatewayEventTypeUserUpdate
type gatewayHandlerUserUpdate struct{}

// EventType returns the discord.GatewayEventTypeUserUpdate
func (h *gatewayHandlerUserUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeUserUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerUserUpdate) New() interface{} {
	return &discord.OAuth2User{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerUserUpdate) HandleGatewayEvent(bot *core.Bot, sequenceNumber discord.GatewaySequence, shardID int, v interface{}) {
	payload := *v.(*discord.OAuth2User)

	var oldSelfUser *core.SelfUser
	if bot.SelfUser != nil {
		selfUser := *bot.SelfUser
		oldSelfUser = &selfUser
	}

	bot.EventManager.Dispatch(&events.SelfUpdateEvent{
		GenericEvent: events.NewGenericEvent(bot, sequenceNumber, shardID),
		SelfUser:     bot.EntityBuilder.CreateSelfUser(payload, core.CacheStrategyYes),
		OldSelfUser:  oldSelfUser,
	})

}
