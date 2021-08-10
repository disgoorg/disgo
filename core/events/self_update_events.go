package events

import (
	"github.com/DisgoOrg/disgo/core"
)

// SelfUpdateEvent is called when something about this api.User updates
type SelfUpdateEvent struct {
	*GenericEvent
	SelfUser    *core.SelfUser
	OldSelfUser *core.SelfUser
}
