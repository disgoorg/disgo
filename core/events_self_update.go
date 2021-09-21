package core

// SelfUpdateEvent is called when something about this core.User updates
type SelfUpdateEvent struct {
	*GenericEvent
	SelfUser    *SelfUser
	OldSelfUser *SelfUser
}
