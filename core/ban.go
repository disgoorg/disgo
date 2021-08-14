package core

import (
	"context"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

// Ban represents a banned User from a Guild (https://discord.com/developers/docs/resources/guild#ban-object)
type Ban struct {
	discord.Ban
	Disgo   Disgo
	User    *User
	GuildID discord.Snowflake
}

func (b *Ban) Unban(ctx context.Context) rest.Error {
	return b.Disgo.RestServices().GuildService().DeleteBan(ctx, b.GuildID, b.User.ID)
}
