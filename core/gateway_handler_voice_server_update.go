package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

// VoiceServerUpdateHandler handles api.GatewayEventVoiceServerUpdate
type VoiceServerUpdateHandler struct{}

// EventType returns the api.GatewayGatewayEventType
func (h *VoiceServerUpdateHandler) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeVoiceServerUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *VoiceServerUpdateHandler) New() interface{} {
	return &discord.VoiceServerUpdate{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *VoiceServerUpdateHandler) HandleGatewayEvent(bot *Bot, _ int, v interface{}) {
	voiceServerUpdate := *v.(*discord.VoiceServerUpdate)

	if interceptor := bot.VoiceDispatchInterceptor; interceptor != nil {
		interceptor.OnVoiceServerUpdate(&VoiceServerUpdateEvent{
			VoiceServerUpdate: voiceServerUpdate,
			Bot:               bot,
		})
	}
}
