package events

import (
	"github.com/DisgoOrg/disgo/api"
)

type GenericGuildInviteEvent struct {
	GenericGuildEvent
	Code      string
	ChannelID api.Snowflake
}

func (e GenericGuildInviteEvent) GuildChannel() *api.GuildChannel {
	return e.Disgo().Cache().GuildChannel(e.ChannelID)
}

func (e GenericGuildInviteEvent) TextChannel() *api.TextChannel {
	return e.Disgo().Cache().TextChannel(e.ChannelID)
}

func (e GenericGuildInviteEvent) VoiceChannel() *api.VoiceChannel {
	return e.Disgo().Cache().VoiceChannel(e.ChannelID)
}

func (e GenericGuildInviteEvent) StoreChannel() *api.StoreChannel {
	return e.Disgo().Cache().StoreChannel(e.ChannelID)
}

func (e GenericGuildInviteEvent) Category() *api.Category {
	return e.Disgo().Cache().Category(e.ChannelID)
}

type GuildInviteCreateEvent struct {
	GenericGuildInviteEvent
	Invite *api.Invite
}

type GuildInviteDeleteEvent struct {
	GenericGuildInviteEvent
}
