package handlers

import (
	"github.com/DisgoOrg/disgo/api"
)

// VoiceStateUpdateHandler handles api.VoiceStateUpdateGatewayEvent
type VoiceStateUpdateHandler struct{}

// Event returns the raw gateway event Event
func (h VoiceStateUpdateHandler) Event() api.GatewayEventType {
	return api.GatewayEventVoiceStateUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h VoiceStateUpdateHandler) New() interface{} {
	return &api.VoiceStateUpdateEvent{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h VoiceStateUpdateHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
	voiceStateUpdate, ok := i.(*api.VoiceStateUpdateEvent)
	if !ok {
		return
	}

	voiceStateUpdate.Disgo = disgo

	//guild := disgo.Cache().Guild(voiceStateUpdate.GuildID)

	//oldMember := disgo.Cache().Member(voiceStateUpdate.GuildID, voiceStateUpdate.UserID)
	//newMember := voiceStateUpdate.Member

	// TODO update voice state cache
	// TODO fire several events

	if disgo.ApplicationID() == voiceStateUpdate.UserID {
		if interceptor := disgo.VoiceDispatchInterceptor(); interceptor != nil {
			interceptor.OnVoiceStateUpdate(voiceStateUpdate)
		}
	}

}
