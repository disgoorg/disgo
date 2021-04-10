package events

import "time"

type HeartbeatEvent struct {
	GenericEvent
	NewPing time.Duration
	OldPing time.Duration
}
