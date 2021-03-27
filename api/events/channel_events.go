package events

import "github.com/DiscoOrg/disgo/api"

type ChannelEvent struct {
	api.Event
	ChannelID api.Snowflake
}
func (e ChannelEvent) Channel() *api.Channel {
	return e.Disgo.Cache().Channel(e.ChannelID)
}


type DMChannelEvent struct {
	ChannelEvent
}
func (e DMChannelEvent) DMChannel() *api.DMChannel {
	return e.Disgo.Cache().DMChannel(e.ChannelID)
}


type MessageChannelEvent struct {
	ChannelEvent
}
func (e MessageChannelEvent) MessageChannel() *api.MessageChannel {
	return e.Disgo.Cache().MessageChannel(e.ChannelID)
}


type TextChannelEvent struct {
	ChannelEvent
}
func (e TextChannelEvent) TextChannel() *api.TextChannel {
	return e.Disgo.Cache().TextChannel(e.ChannelID)
}


type VoiceChannelEvent struct {
	ChannelEvent
}
func (e VoiceChannelEvent) VoiceChannel() *api.VoiceChannel {
	return e.Disgo.Cache().VoiceChannel(e.ChannelID)
}

type CategoryChannelEvent struct {
	ChannelEvent
}
func (e CategoryChannelEvent) Category() *api.Category {
	return e.Disgo.Cache().Category(e.ChannelID)
}

type StoreChannelEvent struct {
	ChannelEvent
}
func (e StoreChannelEvent) StoreChannel() *api.StoreChannel {
	return e.Disgo.Cache().StoreChannel(e.ChannelID)
}
