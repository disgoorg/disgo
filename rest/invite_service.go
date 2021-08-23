package rest

import (
	"context"

	"github.com/DisgoOrg/disgo/discord"
)

func NewInviteService(client Client) InviteService {
	return nil
}

type InviteService interface {
	Service
	GetInvite(code string) (*discord.Invite, Error)
	CreateInvite(channelID discord.Snowflake, inviteCreate discord.InviteCreate)
	DeleteInvite(code string) (*discord.Invite, Error)
	GetGuildInvites(guildID discord.Snowflake) ([]discord.Invite, Error)
	GetChannelInvites(channelID discord.Snowflake) ([]discord.Invite, Error)
}
