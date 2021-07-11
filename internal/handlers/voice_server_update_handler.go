package handlers

import "github.com/DisgoOrg/disgo/api"

// VoiceServerUpdateHandler handles api.GatewayEventVoiceServerUpdate
type VoiceServerUpdateHandler struct{}

// Event returns the raw gateway event Event
func (h *VoiceServerUpdateHandler) Event() api.GatewayEventType {
	return api.GatewayEventVoiceServerUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *VoiceServerUpdateHandler) New() interface{} {
	return &api.VoiceServerUpdate{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *VoiceServerUpdateHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
	voiceServerUpdate, ok := i.(*api.VoiceServerUpdate)
	if !ok {
		return
	}

	if interceptor := disgo.VoiceDispatchInterceptor(); interceptor != nil {
		interceptor.OnVoiceServerUpdate(&api.VoiceServerUpdateEvent{
			VoiceServerUpdate: voiceServerUpdate,
			Disgo:             disgo,
		})
	}
}
