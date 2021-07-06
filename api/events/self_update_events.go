package events

import (
	"github.com/DisgoOrg/disgo/api"
)

// SelfUpdateEvent is called when something about this api.User updates
type SelfUpdateEvent struct {
	*GenericEvent
	Self    *api.User
	OldSelf *api.User
}
