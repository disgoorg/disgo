package core

import (
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

func (b *Ban) Unban(opts ...rest.RequestOpt) rest.Error {
	return b.Disgo.RestServices().GuildService().DeleteBan(b.GuildID, b.User.ID, opts...)
}
