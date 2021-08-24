package gateway

import (
	"io"
	"time"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/log"
)

// Status is the state that the client is currently in
type Status int

// IsConnected returns whether you can send payloads to the Gateway
func (s Status) IsConnected() bool {
	switch s {
	case StatusWaitingForGuilds, StatusReady:
		return true
	default:
		return false
	}
}

// Indicates how far along the client is to connecting
const (
	StatusUnconnected Status = iota
	StatusConnecting
	StatusReconnecting
	StatusIdentifying
	StatusWaitingForHello
	StatusWaitingForReady
	StatusWaitingForGuilds
	StatusReady
	StatusDisconnected
	StatusResuming
)

type EventHandlerFunc func(gatewayEventType discord.GatewayEventType, sequenceNumber int, payload io.Reader)

// Gateway is what is used to connect to discord
type Gateway interface {
	Logger() log.Logger
	Config() Config
	Open() error
	Close()
	Status() Status
	Send(command GatewayCommand) error
	Latency() time.Duration
}
