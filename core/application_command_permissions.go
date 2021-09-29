package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

type ApplicationCommandPermissions struct {
	discord.ApplicationCommandPermissions
	Bot *Bot
}

// TODO: implement methods to update
