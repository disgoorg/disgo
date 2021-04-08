package events

import (
	"github.com/DisgoOrg/disgo/api"
)

type SelfUpdateEvent struct {
	GenericEvent
	NewSelf *api.User
	OldSelf *api.User
}
