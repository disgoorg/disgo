package core

import "github.com/DisgoOrg/disgo/discord"

type Webhook struct {
	discord.Webhook
	Disgo   Disgo
	GuildID discord.Snowflake
}
