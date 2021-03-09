package constants

// ConnectionStatus is the state that the client is currently in
type ConnectionStatus int

// Indicates how far along the client is to connecting
const (
	Ready ConnectionStatus = iota
	Connecting
	Reconnecting
	WaitingForHello
	WaitingForReady
	Disconnected
	WaitingForGuilds
	Identifying
	Resuming
)
