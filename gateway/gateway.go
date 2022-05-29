package gateway

import (
	"context"
	"io"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/log"
)

// Version defines which discord API version disgo should use to connect to discord.
const Version = 10

// Status is the state that the client is currently in.
type Status int

// IsConnected returns whether the Gateway is connected.
func (s Status) IsConnected() bool {
	switch s {
	case StatusWaitingForHello, StatusIdentifying, StatusWaitingForReady, StatusReady:
		return true
	default:
		return false
	}
}

// Indicates how far along the client is too connecting.
const (
	// StatusUnconnected is the initial state when a new Gateway is created.
	StatusUnconnected Status = iota

	// StatusConnecting is the state when the client is connecting to the Discord gateway.
	StatusConnecting

	// StatusWaitingForHello is the state when the Gateway is waiting for the first discord.GatewayOpcodeHello packet.
	StatusWaitingForHello

	// StatusIdentifying is the state when the Gateway received its first discord.GatewayOpcodeHello packet and now sends a discord.GatewayOpcodeIdentify packet.
	StatusIdentifying

	// StatusResuming is the state when the Gateway received its first discord.GatewayOpcodeHello packet and now sends a discord.GatewayOpcodeResume packet.
	StatusResuming

	// StatusWaitingForReady is the state when the Gateway received sent a discord.GatewayOpcodeIdentify or discord.GatewayOpcodeResume packet and now waits for a discord.GatewayOpcodeDispatch with discord.GatewayEventTypeReady packet.
	StatusWaitingForReady

	// StatusReady is the state when the Gateway received a discord.GatewayOpcodeDispatch with discord.GatewayEventTypeReady packet.
	StatusReady

	// StatusDisconnected is the state when the Gateway is disconnected.
	// Either due to an error or because the Gateway was closed gracefully.
	StatusDisconnected
)

type (
	// EventHandlerFunc is a function that is called when an event is received.
	EventHandlerFunc func(gatewayEventType discord.GatewayEventType, sequenceNumber int, shardID int, payload io.Reader)

	// CreateFunc is a type that is used to create a new Gateway(s).
	CreateFunc func(token string, eventHandlerFunc EventHandlerFunc, closeHandlerFUnc CloseHandlerFunc, opts ...ConfigOpt) Gateway

	// CloseHandlerFunc is a function that is called when the Gateway is closed.
	CloseHandlerFunc func(gateway Gateway, err error)
)

// Gateway is what is used to connect to discord.
type Gateway interface {
	// Logger returns the logger that is used by the Gateway.
	Logger() log.Logger

	// ShardID returns the shard ID that this Gateway is configured to use.
	ShardID() int

	// ShardCount returns the total number of shards that this Gateway is configured to use.
	ShardCount() int

	// SessionID returns the session ID that is used by this Gateway.
	// This may be nil if the Gateway was never connected to Discord, was gracefully closed with websocket.CloseNormalClosure or websocket.CloseGoingAway.
	SessionID() *string

	// LastSequenceReceived returns the last sequence number that was received by the Gateway.
	// This may be nil if the Gateway was never connected to Discord, was gracefully closed with websocket.CloseNormalClosure or websocket.CloseGoingAway.
	LastSequenceReceived() *int

	// GatewayIntents returns the discord.GatewayIntents that are used by this Gateway.
	GatewayIntents() discord.GatewayIntents

	// Open connects this Gateway to the Discord API.
	Open(ctx context.Context) error

	// Close gracefully closes the Gateway with the websocket.CloseNormalClosure code.
	// If the context is done, the Gateway connection will be killed.
	Close(ctx context.Context)

	// CloseWithCode closes the Gateway with the given code & message.
	// If the context is done, the Gateway connection will be killed.
	CloseWithCode(ctx context.Context, code int, message string)

	// Status returns the Status of the Gateway.
	Status() Status

	// Send sends a message to the Discord gateway with the opCode and data.
	// If context is deadline exceeds, the message sending will be aborted.
	Send(ctx context.Context, op discord.GatewayOpcode, data discord.GatewayMessageData) error

	// Latency returns the latency of the Gateway.
	// This is calculated by the time it takes to send a heartbeat and receive a heartbeat ack by discord.
	Latency() time.Duration
}
