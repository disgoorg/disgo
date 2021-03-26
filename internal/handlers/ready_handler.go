package handlers

import (
	"github.com/DiscoOrg/disgo/api"
)

// ReadyEventData is the ReadyEvent.D payload
type ReadyEventData struct {
	Version         int             `json:"v"`
	SelfUser        api.User        `json:"user"`
	PrivateChannels []api.DMChannel `json:"private_channels"`
	Guilds          []api.Guild     `json:"guilds"`
	SessionID       string          `json:"session_id"`
	Shard           *[2]int         `json:"shard,omitempty"`
}

type ReadyHandler struct{}

func (h ReadyHandler) New() interface{} {
	return &ReadyEventData{}
}

func (h ReadyHandler) Handle(disgo api.Disgo, eventManager api.EventManager, i interface{}) {
	readyEvent, ok := i.(*ReadyEventData)
	if !ok {
		return
	}
	disgo.SetSelfUser(&readyEvent.SelfUser)
	for i := range readyEvent.Guilds {
		disgo.Cache().CacheGuild(&readyEvent.Guilds[i])
	}

}
