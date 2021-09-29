package events

import "time"

// HeartbeatEvent is called upon sending a heartbeat to the core.Gateway
type HeartbeatEvent struct {
	*GenericEvent
	NewPing time.Duration
	OldPing time.Duration
}
