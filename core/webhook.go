package core

import "github.com/DisgoOrg/disgo/discord"

type Webhook struct {
	discord.Webhook
	Bot *Bot
}

// TODO: add update/delete
