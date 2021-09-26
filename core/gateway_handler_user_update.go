package core

import (
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
func (h *gatewayHandlerUserUpdate) HandleGatewayEvent(bot *Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*discord.OAuth2User)

	var oldSelfUser *SelfUser
	if bot.SelfUser != nil {
		selfUser := *bot.SelfUser
		oldSelfUser = &selfUser
	}

	bot.EventManager.Dispatch(&SelfUpdateEvent{
		GenericEvent: NewGenericEvent(bot, sequenceNumber),
		SelfUser:     bot.EntityBuilder.CreateSelfUser(payload, CacheStrategyYes),
		OldSelfUser:  oldSelfUser,
	})

}
