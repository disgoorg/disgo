package events

import (
	
	"github.com/DisgoOrg/disgo/discord"
)

// GenericGuildInviteEvent is called upon receiving GuildInviteCreateEvent or GuildInviteDeleteEvent(requires api.GatewayIntentsGuildInvites)
type GenericGuildInviteEvent struct {
	*GenericGuildEvent
	Code      string
	ChannelID discord.Snowflake
}

// Channel returns the api.Channel the GenericGuildInviteEvent happened in(returns nil if the api.Channel is uncached or api.Cache is disabled)
func (e GenericGuildInviteEvent) Channel() *core.Channel {
	return e.Disgo().Cache().Channel(e.ChannelID)
}

// GuildChannel returns the api.GuildChannel the GenericGuildInviteEvent happened in(returns nil if the api.GuildChannel is uncached or api.Cache is disabled)
func (e GenericGuildInviteEvent) GuildChannel() *core.GuildChannel {
	return e.Disgo().Cache().GuildChannel(e.ChannelID)
}

// MessageChannel returns the api.MessageChannel the GenericGuildInviteEvent happened in(returns nil if the api.MessageChannel is uncached or api.Cache is disabled)
func (e GenericGuildInviteEvent) MessageChannel() *core.TextChannel {
	return e.Disgo().Cache().TextChannel(e.ChannelID)
}

// TextChannel returns the api.TextChannel the GenericGuildInviteEvent happened in(returns nil if the api.TextChannel is uncached or api.CacheFlagTextChannels is disabled)
func (e GenericGuildInviteEvent) TextChannel() *core.TextChannel {
	return e.Disgo().Cache().TextChannel(e.ChannelID)
}

// VoiceChannel returns the api.VoiceChannel the GenericGuildInviteEvent happened in(returns nil if the api.VoiceChannel is uncached or api.CacheFlagVoiceChannels is disabled)
func (e GenericGuildInviteEvent) VoiceChannel() *core.VoiceChannel {
	return e.Disgo().Cache().VoiceChannel(e.ChannelID)
}

// StoreChannel returns the api.StoreChannel the GenericGuildInviteEvent happened in(returns nil if the api.StoreChannel is uncached or api.CacheFlagStoreChannels is disabled)
func (e GenericGuildInviteEvent) StoreChannel() *core.StoreChannel {
	return e.Disgo().Cache().StoreChannel(e.ChannelID)
}

// Category returns the api.Category the GenericGuildInviteEvent happened in(returns nil if the api.Category is uncached or api.CacheFlagCategories is disabled)
func (e GenericGuildInviteEvent) Category() *core.Category {
	return e.Disgo().Cache().Category(e.ChannelID)
}

// GuildInviteCreateEvent is called upon creation of a new api.Invite in an api.Guild(requires api.GatewayIntentsGuildInvites)
type GuildInviteCreateEvent struct {
	*GenericGuildInviteEvent
	Invite *core.Invite
}

// GuildInviteDeleteEvent is called upon deletion of a new api.Invite in an api.Guild(requires api.GatewayIntentsGuildInvites)
type GuildInviteDeleteEvent struct {
	*GenericGuildInviteEvent
}
