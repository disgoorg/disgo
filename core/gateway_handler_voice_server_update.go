package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

// gatewayHandlerVoiceServerUpdate handles core.GatewayEventVoiceServerUpdate
type gatewayHandlerVoiceServerUpdate struct{}

// EventType returns the core.GatewayGatewayEventType
func (h *gatewayHandlerVoiceServerUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeVoiceServerUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerVoiceServerUpdate) New() interface{} {
	return &discord.VoiceServerUpdate{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerVoiceServerUpdate) HandleGatewayEvent(bot *Bot, _ int, v interface{}) {
	voiceServerUpdate := *v.(*discord.VoiceServerUpdate)

	if interceptor := bot.VoiceDispatchInterceptor; interceptor != nil {
		interceptor.OnVoiceServerUpdate(&VoiceServerUpdateEvent{
			VoiceServerUpdate: voiceServerUpdate,
			Bot:               bot,
		})
	}
}
