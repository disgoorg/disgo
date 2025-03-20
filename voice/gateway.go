package voice

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"sync"
	"syscall"
	"time"

	"github.com/disgoorg/json"
	"github.com/disgoorg/snowflake/v2"
	"github.com/gorilla/websocket"

	"github.com/disgoorg/disgo/discord"
)

var (
	// ErrGatewayNotConnected is returned when the gateway is not connected and a message is attempted to be sent.
	ErrGatewayNotConnected = fmt.Errorf("voice gateway not connected")

	// ErrGatewayAlreadyConnected is returned when the gateway is already connected and a connection is attempted to be opened.
	ErrGatewayAlreadyConnected = fmt.Errorf("voice gateway already connected")
)

// GatewayVersion is the version of the voice gateway we are using.
const GatewayVersion = 8

// Status returns the current status of the gateway.
type Status int

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

type (
	// EventHandlerFunc is a function that handles a voice gateway event.
	EventHandlerFunc func(opCode Opcode, data GatewayMessageData)

	// CloseHandlerFunc is a function that handles a voice gateway close.
	CloseHandlerFunc func(gateway Gateway, err error)

	// GatewayCreateFunc is a function that creates a new voice gateway.
	GatewayCreateFunc func(eventHandlerFunc EventHandlerFunc, closeHandlerFunc CloseHandlerFunc, opts ...GatewayConfigOpt) Gateway

	// StateProviderFunc is a function that provides the current conn state of the voice gateway.
	StateProviderFunc func() State
)

// State is the current state of the voice conn.
type State struct {
	GuildID snowflake.ID
	UserID  snowflake.ID

	ChannelID *snowflake.ID
	SessionID string
	Token     string
	Endpoint  string
}

// Gateway is a websocket connection to the Discord voice gateway.
type Gateway interface {
	// SSRC returns the SSRC of the current voice connection.
	SSRC() uint32

	// Latency returns the current latency of the voice gateway connection.
	Latency() time.Duration

	// Open opens a new websocket connection to the voice gateway.
	Open(ctx context.Context, state State) error

	// Close closes the websocket connection to the voice gateway.
	Close()

	// CloseWithCode closes the websocket connection to the voice gateway with a specific close code.
	CloseWithCode(code int, message string)

	// Send sends a message to the voice gateway.
	Send(ctx context.Context, opCode Opcode, data GatewayMessageData) error
}

// NewGateway creates a new voice Gateway.
func NewGateway(eventHandlerFunc EventHandlerFunc, closeHandlerFunc CloseHandlerFunc, opts ...GatewayConfigOpt) Gateway {
	config := DefaultGatewayConfig()
	config.Apply(opts)
	config.Logger = config.Logger.With(slog.String("name", "voice_conn_gateway"))

	return &gatewayImpl{
		config:           *config,
		eventHandlerFunc: eventHandlerFunc,
		closeHandlerFunc: closeHandlerFunc,
	}
}

type gatewayImpl struct {
	config           GatewayConfig
	eventHandlerFunc EventHandlerFunc
	closeHandlerFunc CloseHandlerFunc

	ssrc  uint32
	state State
	seq   int

	conn   *websocket.Conn
	connMu sync.Mutex
	status Status

	heartbeatTicker       *time.Ticker
	heartbeatInterval     time.Duration
	lastHeartbeatSent     time.Time
	lastHeartbeatReceived time.Time
	lastNonce             int64
}

func (g *gatewayImpl) SSRC() uint32 {
	return g.ssrc
}

func (g *gatewayImpl) Open(ctx context.Context, state State) error {
	g.config.Logger.Debug("opening voice gateway connection")
	g.connMu.Lock()
	defer g.connMu.Unlock()
	if g.conn != nil {
		return ErrGatewayAlreadyConnected
	}
	g.state = state
	g.status = StatusConnecting

	gatewayURL := fmt.Sprintf("wss://%s?v=%d", state.Endpoint, GatewayVersion)
	g.config.Logger.Debug("connecting to voice gateway at", slog.String("url", gatewayURL))
	g.lastHeartbeatSent = time.Now().UTC()
	conn, rs, err := g.config.Dialer.DialContext(ctx, gatewayURL, nil)
	if err != nil {
		g.Close()
		defer rs.Body.Close()
		return fmt.Errorf("error connecting to voice gateway: %w", err)
	}

	conn.SetCloseHandler(func(code int, text string) error {
		return nil
	})

	g.conn = conn
	g.status = StatusWaitingForHello

	go g.listen(g.conn)
	return nil
}

func (g *gatewayImpl) Close() {
	g.CloseWithCode(websocket.CloseNormalClosure, "Shutting down")
}

func (g *gatewayImpl) CloseWithCode(code int, message string) {
	if g.heartbeatTicker != nil {
		g.config.Logger.Debug("closing heartbeat goroutines")
		g.heartbeatTicker.Stop()
		g.heartbeatTicker = nil
	}

	g.connMu.Lock()
	defer g.connMu.Unlock()
	if g.conn != nil {
		g.config.Logger.Debug("closing voice gateway connection", slog.Int("code", code), slog.String("message", message))
		if err := g.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(code, message)); err != nil && !errors.Is(err, websocket.ErrCloseSent) {
			g.config.Logger.Debug("error writing close code", slog.Any("err", err))
		}
		_ = g.conn.Close()
		g.conn = nil

		// clear resume data as we closed gracefully
		if code == websocket.CloseNormalClosure || code == websocket.CloseGoingAway {
			g.ssrc = 0
			g.seq = 0
		}
	}
}

func (g *gatewayImpl) heartbeat() {
	g.heartbeatTicker = time.NewTicker(g.heartbeatInterval)
	defer g.heartbeatTicker.Stop()
	defer g.config.Logger.Debug("exiting voice heartbeat goroutine")

	for range g.heartbeatTicker.C {
		g.sendHeartbeat()
	}
}

func (g *gatewayImpl) sendHeartbeat() {
	g.lastNonce = time.Now().UnixMilli()
	ctx, cancel := context.WithTimeout(context.Background(), g.heartbeatInterval)
	defer cancel()

	if err := g.Send(ctx, OpcodeHeartbeat, GatewayMessageDataHeartbeat{
		T:      g.lastNonce,
		SeqAck: g.seq,
	}); err != nil {
		if !errors.Is(err, ErrGatewayNotConnected) || errors.Is(err, syscall.EPIPE) {
			return
		}
		g.config.Logger.Error("failed to send heartbeat", slog.Any("err", err))
		g.CloseWithCode(websocket.CloseServiceRestart, "heartbeat timeout")
		go g.reconnect()
		return
	}
	g.lastHeartbeatSent = time.Now().UTC()
}

func (g *gatewayImpl) listen(conn *websocket.Conn) {
	defer g.config.Logger.Debug("exiting listen goroutine")
loop:
	for {
		_, reader, err := conn.NextReader()
		if err != nil {
			g.connMu.Lock()
			sameConn := g.conn == conn
			g.connMu.Unlock()

			// if sameConn is false, it means the connection has been closed by the user, and we can just exit
			if !sameConn {
				return
			}

			reconnect := true
			var closeError *websocket.CloseError
			if errors.As(err, &closeError) {
				closeCode := GatewayCloseEventCodeByCode(closeError.Code)
				reconnect = closeCode.Reconnect
			}
			g.CloseWithCode(websocket.CloseServiceRestart, "listen error")
			if g.config.AutoReconnect && reconnect {
				go g.reconnect()
			} else if g.closeHandlerFunc != nil {
				go g.closeHandlerFunc(g, err)
			}
			break loop
		}

		message, err := g.parseMessage(reader)
		if err != nil {
			g.config.Logger.Error("error while parsing voice gateway event", slog.Any("err", err))
			continue
		}

		if message.Seq > 0 {
			g.seq = message.Seq
		}

		switch d := message.D.(type) {
		case GatewayMessageDataHello:
			g.status = StatusWaitingForReady
			g.lastHeartbeatReceived = time.Now().UTC()
			g.heartbeatInterval = time.Duration(d.HeartbeatInterval) * time.Millisecond
			go g.heartbeat()

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			if g.ssrc == 0 || g.seq == 0 {
				g.status = StatusIdentifying
				err = g.Send(ctx, OpcodeIdentify, GatewayMessageDataIdentify{
					GuildID:   g.state.GuildID,
					UserID:    g.state.UserID,
					SessionID: g.state.SessionID,
					Token:     g.state.Token,
				})
			} else {
				g.status = StatusResuming
				err = g.Send(ctx, OpcodeResume, GatewayMessageDataResume{
					GuildID:   g.state.GuildID,
					SessionID: g.state.SessionID,
					Token:     g.state.Token,
					SeqAck:    g.seq,
				})
			}
			cancel()
			if err != nil {
				g.CloseWithCode(websocket.CloseServiceRestart, "failed to send identify or resume")
				go g.reconnect()
				return
			}

		case GatewayMessageDataReady:
			g.status = StatusReady
			g.ssrc = d.SSRC

		case GatewayMessageDataHeartbeatACK:
			if int64(d) != g.lastNonce {
				g.config.Logger.Error("received heartbeat ack with nonce", slog.Int64("nonce", int64(d)), slog.Int64("last_nonce", g.lastNonce))
				go g.reconnect()
				break loop
			}
			g.lastHeartbeatReceived = time.Now().UTC()
		}
		g.eventHandlerFunc(message.Op, message.D)
	}
}

func (g *gatewayImpl) Latency() time.Duration {
	return g.lastHeartbeatReceived.Sub(g.lastHeartbeatSent)
}

func (g *gatewayImpl) Send(ctx context.Context, op Opcode, d GatewayMessageData) error {
	data, err := json.Marshal(GatewayMessage{
		Op: op,
		D:  d,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal voice gateway message: %w", err)
	}
	return g.send(ctx, websocket.TextMessage, data)
}

func (g *gatewayImpl) send(ctx context.Context, messageType int, data []byte) error {
	g.connMu.Lock()
	defer g.connMu.Unlock()
	if g.conn == nil {
		return ErrGatewayNotConnected
	}

	g.config.Logger.Debug("sending message to voice gateway", slog.String("data", string(data)))
	deadline, ok := ctx.Deadline()
	if ok {
		if err := g.conn.SetWriteDeadline(deadline); err != nil {
			return err
		}
	}
	if err := g.conn.WriteMessage(messageType, data); err != nil {
		return fmt.Errorf("failed to send message to voice gateway: %w", err)
	}
	return nil
}

func (g *gatewayImpl) reconnectTry(ctx context.Context, try int) error {
	delay := time.Duration(try) * 2 * time.Second
	if delay > 30*time.Second {
		delay = 30 * time.Second
	}

	timer := time.NewTimer(time.Duration(try) * delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		timer.Stop()
		return ctx.Err()
	case <-timer.C:
	}

	g.config.Logger.Debug("reconnecting voice gateway")
	if err := g.Open(ctx, g.state); err != nil {
		if errors.Is(err, discord.ErrGatewayAlreadyConnected) {
			return err
		}
		g.config.Logger.Error("failed to reconnect voice gateway", slog.Any("err", err))
		g.status = StatusDisconnected
		return g.reconnectTry(ctx, try+1)
	}
	return nil
}

func (g *gatewayImpl) reconnect() {
	if err := g.reconnectTry(context.Background(), 0); err != nil {
		g.config.Logger.Error("failed to reopen voice gateway", slog.Any("err", err))
	}
}

func (g *gatewayImpl) parseMessage(r io.Reader) (GatewayMessage, error) {
	buff := &bytes.Buffer{}
	data, _ := io.ReadAll(io.TeeReader(r, buff))
	g.config.Logger.Debug("received message from voice gateway", slog.String("data", string(data)))

	var message GatewayMessage
	if err := json.NewDecoder(buff).Decode(&message); err != nil {
		return GatewayMessage{}, err
	}
	return message, nil
}
