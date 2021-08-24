package rest

import (
	"github.com/DisgoOrg/disgo/discord"
)

func NewInviteService(client Client) InviteService {
	return nil
}

type InviteService interface {
	Service
	GetInvite(code string, opts ...RequestOpt) (*discord.Invite, Error)
	CreateInvite(channelID discord.Snowflake, inviteCreate discord.InviteCreate, opts ...RequestOpt) (discord.Invite, Error)
	DeleteInvite(code string, opts ...RequestOpt) (*discord.Invite, Error)
	GetGuildInvites(guildID discord.Snowflake, opts ...RequestOpt) ([]discord.Invite, Error)
	GetChannelInvites(channelID discord.Snowflake, opts ...RequestOpt) ([]discord.Invite, Error)
}
