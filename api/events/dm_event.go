package events

import "github.com/DiscoOrg/disgo/api"

type GenericDMEvent struct {
	api.Event
	UserID      api.Snowflake
	DMChannelID api.Snowflake
}

func (e GenericDMEvent) User() *api.DMChannel {
	return e.Disgo.Cache().DMChannel(e.DMChannelID)
}

func (e GenericDMEvent) DMChannel() *api.DMChannel {
	return e.Disgo.Cache().DMChannel(e.DMChannelID)
}