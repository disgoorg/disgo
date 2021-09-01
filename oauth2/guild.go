package oauth2

import "github.com/DisgoOrg/disgo/discord"

type Guild struct {
	discord.PartialGuild
	Client Client
}
