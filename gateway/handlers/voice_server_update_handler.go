package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
)

// VoiceServerUpdateHandler handles api.GatewayEventVoiceServerUpdate
type VoiceServerUpdateHandler struct{}

// EventType returns the api.GatewayEventType
func (h *VoiceServerUpdateHandler) EventType() gateway.EventType {
	return gateway.EventTypeVoiceServerUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *VoiceServerUpdateHandler) New() interface{} {
	return &discord.VoiceServerUpdate{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *VoiceServerUpdateHandler) HandleGatewayEvent(disgo core.Disgo, _ core.EventManager, _ int, i interface{}) {
	voiceServerUpdate, ok := i.(*discord.VoiceServerUpdate)
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
