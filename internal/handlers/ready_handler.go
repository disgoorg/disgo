package handlers

import (
	"github.com/DiscoOrg/disgo/api"
)

// ReadyEventData is the ReadyEvent.D payload
type ReadyEventData struct {
	User            api.User               `json:"user"`
	PrivateChannels []api.DMChannel        `json:"channel"`
	Guilds          []api.UnavailableGuild `json:"guild_events"`
	SessionID       string                 `json:"session_id"`
	Shard           [2]int                 `json:"shard,omitempty"`
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
		disgo.Cache().CacheUnavailableGuild(&readyData.Guilds[i])
	}

}
