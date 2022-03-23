package gateway

import (
	"context"
	"io"
	"time"

	"github.com/DisgoOrg/log"
	"github.com/disgoorg/disgo/discord"
)

// Status is the state that the client is currently in
type Status int

// IsConnected returns whether the Gateway is connected
func (s Status) IsConnected() bool {
	switch s {
	case StatusWaitingForHello, StatusIdentifying, StatusWaitingForReady, StatusReady:
		return true
	default:
		return false
	}
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

type EventHandlerFunc func(gatewayEventType discord.GatewayEventType, sequenceNumber discord.GatewaySequence, payload io.Reader)

type CreateFunc func(token string, opts ...ConfigOpt) Gateway

// Gateway is what is used to connect to discord
type Gateway interface {
	Logger() log.Logger
	ShardID() int
	ShardCount() int
	GatewayIntents() discord.GatewayIntents

	Open(ctx context.Context) error
	ReOpen(ctx context.Context, delay time.Duration) error
	Close(ctx context.Context)
	CloseWithCode(ctx context.Context, code int, message string)
	Status() Status
	Send(ctx context.Context, op discord.GatewayOpcode, data discord.GatewayCommandData) error
	Latency() time.Duration
}
