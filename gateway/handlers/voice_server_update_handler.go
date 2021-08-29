package handlers

import (
	"github.com/DisgoOrg/disgo/core"
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
	return discord.VoiceServerUpdate{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *VoiceServerUpdateHandler) HandleGatewayEvent(disgo core.Disgo, _ core.EventManager, _ int, v interface{}) {
	voiceServerUpdate, ok := v.(discord.VoiceServerUpdate)
	if !ok {
		return
	}

	if interceptor := disgo.VoiceDispatchInterceptor(); interceptor != nil {
		interceptor.OnVoiceServerUpdate(&core.VoiceServerUpdateEvent{
			VoiceServerUpdate: voiceServerUpdate,
			Disgo:             disgo,
		})
	}
}
