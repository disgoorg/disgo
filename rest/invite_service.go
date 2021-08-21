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
	GetInvite(ctx context.Context, code string) (*discord.Invite, Error)
	CreateInvite(ctx context.Context, channelID discord.Snowflake, inviteCreate discord.InviteCreate)
	DeleteInvite(ctx context.Context, code string) (*discord.Invite, Error)
	GetGuildInvites(ctx context.Context, guildID discord.Snowflake) ([]discord.Invite, Error)
	GetChannelInvites(ctx context.Context, channelID discord.Snowflake) ([]discord.Invite, Error)
}
