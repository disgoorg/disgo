package handlers

import (
	"github.com/DisgoOrg/disgo/api"
)

// VoiceStateUpdateHandler handles api.VoiceStateUpdateGatewayEvent
type VoiceStateUpdateHandler struct{}

// Event returns the raw gateway event Event
func (h VoiceStateUpdateHandler) Event() api.GatewayEventName {
	return api.GatewayEventVoiceStateUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h VoiceStateUpdateHandler) New() interface{} {
	return &api.VoiceStateUpdate{}
}

// Handle handles the specific raw gateway event
func (h VoiceStateUpdateHandler) Handle(disgo api.Disgo, eventManager api.EventManager, i interface{}) {
	voiceStateUpdate, ok := i.(*api.VoiceStateUpdate)
	if !ok {
		return
	}

	voiceStateUpdate.Disgo = disgo

	//guild := disgo.Cache().Guild(voiceStateUpdate.GuildID)

	//oldMember := disgo.Cache().Member(voiceStateUpdate.GuildID, voiceStateUpdate.UserID)
	//newMember := voiceStateUpdate.Member

	// TODO update voice state cache
	// TODO fire several events

	if disgo.SelfUserID() == voiceStateUpdate.UserID {
		if interceptor := disgo.VoiceDispatchInterceptor(); interceptor != nil {
			interceptor.OnVoiceStateUpdate(*voiceStateUpdate)
		}
	}

}
