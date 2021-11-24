package core

import "github.com/DisgoOrg/disgo/discord"

type GuildScheduledEvent struct {
	discord.GuildScheduledEvent
	Bot *Bot
}
