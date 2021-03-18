package handlers

import (
	log "github.com/sirupsen/logrus"

	"github.com/DiscoOrg/disgo/api"
)

// ReadyEventData is the ReadyEvent.D payload
type ReadyEventData struct {
	V               int                    `json:"v"`
	User            api.User               `json:"user"`
	PrivateChannels []api.DMChannel        `json:"private_channels"`
	Guilds          []api.UnavailableGuild `json:"guilds"`
	SessionID       string                 `json:"session_id"`
	Shard           *[2]int                `json:"shard,omitempty"`
}

type ReadyHandler struct{}

func (h ReadyHandler) New() interface{} {
	return &ReadyEventData{}
}

func (h ReadyHandler) Handle(disgo api.Disgo, eventManager api.EventManager, i interface{}) {
	readyData, ok := i.(*ReadyEventData)
	if !ok {
		return
	}
	for i := range readyData.Guilds {
		log.Infof("added unavail guild: %#v", readyData.Guilds[i])
		disgo.Cache().CacheUnavailableGuild(&readyData.Guilds[i])
	}

}
