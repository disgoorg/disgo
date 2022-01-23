package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/snowflake"
)

// Ban represents a banned User from a Guild (https://discord.com/developers/docs/resources/guild#ban-object)
type Ban struct {
	discord.Ban
	Bot     *Bot
	User    *User
	GuildID snowflake.Snowflake
}

// Unban unbans the User associated with this Ban from the Guild
func (b *Ban) Unban(opts ...rest.RequestOpt) error {
	return b.Bot.RestServices.GuildService().DeleteBan(b.GuildID, b.User.ID, opts...)
}
