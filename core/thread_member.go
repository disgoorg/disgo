package core

import "github.com/DisgoOrg/disgo/discord"

type ThreadMember struct {
	discord.ThreadMember
	Bot     *Bot
	GuildID discord.Snowflake
}
