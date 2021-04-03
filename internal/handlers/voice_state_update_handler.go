package handlers

import (
	log "github.com/sirupsen/logrus"

	"github.com/DisgoOrg/disgo/api"
)

// VoiceStateUpdateHandler handles api.VoiceStateUpdateGatewayEvent
type VoiceStateUpdateHandler struct{}

// Name returns the raw gateway event name
func (h VoiceStateUpdateHandler) Name() string {
	return api.VoiceStateUpdateGatewayEvent
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

	log.Printf("%+v", *voiceStateUpdate)

	voiceStateUpdate.Disgo = disgo

	//guild := disgo.Cache().Guild(voiceStateUpdate.GuildID)

	//oldMember := disgo.Cache().Member(voiceStateUpdate.GuildID, voiceStateUpdate.UserID)
	//newMember := voiceStateUpdate.Member

	// TODO update voice state cache
	// TODO fire several events

	if disgo.ApplicationID() == voiceStateUpdate.UserID {
		if interceptor := disgo.VoiceDispatchInterceptor(); interceptor != nil {
			interceptor.OnVoiceStateUpdate(*voiceStateUpdate)
		}
	}

}
