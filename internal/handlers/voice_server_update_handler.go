package handlers

import "github.com/DisgoOrg/disgo/api"

// InteractionCreateHandler handles api.VoiceServerUpdateGatewayEvent
type VoiceServerUpdateHandler struct{}

// Name returns the raw gateway event name
func (h VoiceServerUpdateHandler) Name() string {
	return api.VoiceServerUpdateGatewayEvent
}

// New constructs a new payload receiver for the raw gateway event
func (h VoiceServerUpdateHandler) New() interface{} {
	return &api.VoiceServerUpdate{}
}

// Handle handles the specific raw gateway event
func (h VoiceServerUpdateHandler) Handle(disgo api.Disgo, eventManager api.EventManager, i interface{}) {
	voiceServerUpdate, ok := i.(*api.VoiceServerUpdate)
	if !ok {
		return
	}

	if voiceServerUpdate.Endpoint == nil {
		return
	}

	voiceServerUpdate.Disgo = disgo

	if interceptor := disgo.VoiceDispatchInterceptor(); interceptor != nil {
		interceptor.OnVoiceServerUpdate(*voiceServerUpdate)
	}
}
