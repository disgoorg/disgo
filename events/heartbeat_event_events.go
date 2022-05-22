package events

import "time"

// Heartbeat is called upon sending a heartbeat to the gateway.Gateway
type Heartbeat struct {
	*GenericEvent
	NewPing time.Duration
	OldPing time.Duration
}
