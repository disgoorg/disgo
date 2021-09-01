package oauth2

import "github.com/DisgoOrg/disgo/discord"

type User struct {
	discord.OAuth2User
	Client Client
}
