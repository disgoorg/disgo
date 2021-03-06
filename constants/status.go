package constants

type ConnectionStatus int

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