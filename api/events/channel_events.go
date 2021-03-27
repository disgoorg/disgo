package events

import "github.com/DiscoOrg/disgo/api"

// GenericChannelEvent is called upon receiving an event in a api.Channel
type GenericChannelEvent struct {
	api.Event
	ChannelID api.Snowflake
}

// Channel returns the api.Channel from the api.Cache
func (e GenericChannelEvent) Channel() *api.Channel {
	return e.Disgo.Cache().Channel(e.ChannelID)
}

// GenericDMChannelEvent is called upon receiving an event in a api.DMChannel
type GenericDMChannelEvent struct {
	GenericChannelEvent
}

// DMChannel returns the api.DMChannel from the api.Cache
func (e GenericDMChannelEvent) DMChannel() *api.DMChannel {
	return e.Disgo.Cache().DMChannel(e.ChannelID)
}

// GenericMessageChannelEvent is called upon receiving an event in a api.MessageChannel
type GenericMessageChannelEvent struct {
	GenericChannelEvent
}

// MessageChannel returns the api.MessageChannel from the api.Cache
func (e GenericMessageChannelEvent) MessageChannel() *api.MessageChannel {
	return e.Disgo.Cache().MessageChannel(e.ChannelID)
}

// GenericTextChannelEvent is called upon receiving an event in a api.TextChannel
type GenericTextChannelEvent struct {
	GenericChannelEvent
}

// TextChannel returns the api.TextChannel from the api.Cache
func (e GenericTextChannelEvent) TextChannel() *api.TextChannel {
	return e.Disgo.Cache().TextChannel(e.ChannelID)
}

// GenericVoiceChannelEvent is called upon receiving an event in a api.VoiceChannel
type GenericVoiceChannelEvent struct {
	GenericChannelEvent
}

// VoiceChannel returns the api.VoiceChannel from the api.Cache
func (e GenericVoiceChannelEvent) VoiceChannel() *api.VoiceChannel {
	return e.Disgo.Cache().VoiceChannel(e.ChannelID)
}

// GenericCategoryEvent is called upon receiving an event in a api.Category
type GenericCategoryEvent struct {
	GenericChannelEvent
}

// Category returns the api.Category from the api.Cache
func (e GenericCategoryEvent) Category() *api.Category {
	return e.Disgo.Cache().Category(e.ChannelID)
}

// GenericStoreChannelEvent is called upon receiving an event in a api.StoreChannel
type GenericStoreChannelEvent struct {
	GenericChannelEvent
}

// StoreChannel returns the api.StoreChannel from the api.Cache
func (e GenericStoreChannelEvent) StoreChannel() *api.StoreChannel {
	return e.Disgo.Cache().StoreChannel(e.ChannelID)
}
