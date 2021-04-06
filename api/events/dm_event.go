package events

import "github.com/DisgoOrg/disgo/api"

// GenericDMEvent is a generic dm channel event
type GenericDMEvent struct {
	api.GenericEvent
	UserID      api.Snowflake
	DMChannelID api.Snowflake
}

// User gets the user from the api.Cache
func (e GenericDMEvent) User() *api.DMChannel {
	return e.Disgo().Cache().DMChannel(e.DMChannelID)
}

// DMChannel returns the api.DMChannel from the api.Cache
func (e GenericDMEvent) DMChannel() *api.DMChannel {
	return e.Disgo().Cache().DMChannel(e.DMChannelID)
}
