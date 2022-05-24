package events

import (
	"github.com/disgoorg/disgo/discord"
)

// SelfUpdate is called when something about this discord.User updates
type SelfUpdate struct {
	*GenericEvent
	SelfUser    discord.OAuth2User
	OldSelfUser discord.OAuth2User
}
