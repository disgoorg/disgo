package events

import (
	
)

// SelfUpdateEvent is called when something about this api.User updates
type SelfUpdateEvent struct {
	*GenericEvent
	Self    *core.User
	OldSelf *core.User
}
