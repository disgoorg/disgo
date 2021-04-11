package events

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/endpoints"
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

func (e GenericGuildInviteEvent) URL() string {
	url, err := endpoints.InviteURL.Compile(e.Code)
	if err != nil {
		return ""
	}
	return url.Route()
}

type GuildInviteCreateEvent struct {
	GenericGuildInviteEvent
	Invite *api.Invite
}

type GuildInviteDeleteEvent struct {
	GenericGuildInviteEvent
}
