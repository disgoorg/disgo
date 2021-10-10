package gateway

import (
	"context"
	"io"
	"time"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/log"
)

// Status is the state that the client is currently in
type Status int

// IsConnected returns whether you can send payloads to the Gateway
func (s Status) IsConnected() bool {
	return s == StatusReady
}

// Indicates how far along the client is to connecting
const (
	StatusUnconnected Status = iota
	StatusConnecting
	StatusWaitingForHello
	StatusIdentifying
	StatusResuming
	StatusWaitingForReady
	StatusReady
	StatusDisconnected
)

type EventHandlerFunc func(gatewayEventType discord.GatewayEventType, sequenceNumber int, payload io.Reader)

// Gateway is what is used to connect to discord
type Gateway interface {
	Logger() log.Logger
	Config() Config
	ShardID() int
	ShardCount() int
	Open() error
	OpenCtx(ctx context.Context) error
	Close()
	CloseWithCode(code int)
	Status() Status
	Send(command discord.GatewayCommand) error
	SendCtx(ctx context.Context, command discord.GatewayCommand) error
	Latency() time.Duration
}
