package events

import (
	"github.com/DisgoOrg/disgo/discord"
)

// SelfUpdateEvent is called when something about this discord.User updates
type SelfUpdateEvent struct {
	*GenericEvent
	SelfUser    discord.OAuth2User
	OldSelfUser discord.OAuth2User
}
