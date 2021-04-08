package events

type HeartbeatEvent struct {
	GenericEvent
	NewPing int
	OldPing int
}
