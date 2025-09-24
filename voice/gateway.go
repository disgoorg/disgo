package voice

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"sync"
	"syscall"
	"time"

	"github.com/disgoorg/json/v2"
	"github.com/disgoorg/snowflake/v2"
	"github.com/gorilla/websocket"

	"github.com/disgoorg/disgo/discord"
)

// GatewayVersion is the version of the voice gateway we are using.
const GatewayVersion = 8

const maximumConnectDelay = 10 * time.Second

var (
	// ErrGatewayNotConnected is returned when the gateway is not connected and a message is attempted to be sent.
	ErrGatewayNotConnected = fmt.Errorf("voice gateway not connected")

	// ErrGatewayAlreadyConnected is returned when the gateway is already connected and a connection is attempted to be opened.
	ErrGatewayAlreadyConnected = fmt.Errorf("voice gateway already connected")
)

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
	EventHandlerFunc func(gateway Gateway, opCode Opcode, sequenceNumber int, data GatewayMessageData)

	// GatewayCreateFunc is a function that creates a new voice gateway.
	GatewayCreateFunc func(eventHandlerFunc EventHandlerFunc, closeHandlerFunc CloseHandlerFunc, opts ...GatewayConfigOpt) Gateway

	// CloseHandlerFunc is a function that handles a voice gateway close.
	CloseHandlerFunc func(gateway Gateway, err error, reconnect bool)

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

	// Open opens a new websocket connection to the voice gateway.
	Open(ctx context.Context, state State) error

	// Close closes the websocket connection to the voice gateway.
	Close()

	// CloseWithCode closes the websocket connection to the voice gateway with a specific close code.
	CloseWithCode(code int, message string)

	// Status returns the Status of the Gateway.
	Status() Status

	// Send sends a message to the voice gateway.
	Send(ctx context.Context, opCode Opcode, data GatewayMessageData) error

	// Latency returns the current latency of the voice gateway connection.
	Latency() time.Duration
}

// NewGateway creates a new voice Gateway.
func NewGateway(eventHandlerFunc EventHandlerFunc, closeHandlerFunc CloseHandlerFunc, opts ...GatewayConfigOpt) Gateway {
	cfg := defaultGatewayConfig()
	cfg.apply(opts)

	return &gatewayImpl{
		config:           cfg,
		eventHandlerFunc: eventHandlerFunc,
		closeHandlerFunc: closeHandlerFunc,
	}
}

type gatewayImpl struct {
	config           gatewayConfig
	eventHandlerFunc EventHandlerFunc
	closeHandlerFunc CloseHandlerFunc

	ssrc  uint32
	state State
	seq   int

	conn            *websocket.Conn
	connMu          sync.Mutex
	heartbeatCancel context.CancelFunc
	status          Status
	statusMu        sync.Mutex

	heartbeatInterval     time.Duration
	lastHeartbeatSent     time.Time
	lastHeartbeatReceived time.Time
	lastNonce             int64
}

func (g *gatewayImpl) SSRC() uint32 {
	return g.ssrc
}

func (g *gatewayImpl) Open(ctx context.Context, state State) error {
	return g.doReconnect(ctx, state)
}

func (g *gatewayImpl) open(ctx context.Context, state State) error {
	g.config.Logger.Debug("opening voice gateway connection")

	g.connMu.Lock()
	if g.conn != nil {
		g.connMu.Unlock()
		return ErrGatewayAlreadyConnected
	}

	g.state = state
	g.statusMu.Lock()
	g.status = StatusConnecting
	g.statusMu.Unlock()

	gatewayURL := fmt.Sprintf("wss://%s?v=%d", state.Endpoint, GatewayVersion)
	g.lastHeartbeatSent = time.Now().UTC()
	conn, rs, err := g.config.Dialer.DialContext(ctx, gatewayURL, nil)
	if err != nil {
		var body []byte
		if rs != nil {
			defer func() {
				_ = rs.Body.Close()
			}()
			body, err = io.ReadAll(rs.Body)
			if err != nil {
				g.config.Logger.ErrorContext(ctx, "error while reading response body", slog.Any("err", err))
			}
		}

		g.config.Logger.ErrorContext(ctx, "error connecting to the voice gateway",
			slog.Any("err", err),
			slog.String("url", gatewayURL),
			slog.String("body", string(body)),
		)
		g.connMu.Unlock()
		return err
	}

	conn.SetCloseHandler(func(code int, text string) error {
		return nil
	})

	g.conn = conn
	g.connMu.Unlock()

	g.statusMu.Lock()
	g.status = StatusWaitingForHello
	g.statusMu.Unlock()

	var readyOnce sync.Once
	readyChan := make(chan error)
	go g.listen(g.conn, func(err error) {
		readyOnce.Do(func() {
			readyChan <- err
			close(readyChan)
		})
	})

	select {
	case <-ctx.Done():
		g.Close()
		return ctx.Err()
	case err = <-readyChan:
		if err != nil {
			g.Close()
			return fmt.Errorf("failed to open voice gateway connection: %w", err)
		}
	}

	return nil
}

func (g *gatewayImpl) Close() {
	g.CloseWithCode(websocket.CloseNormalClosure, "Shutting down")
}

func (g *gatewayImpl) CloseWithCode(code int, message string) {
	if g.heartbeatCancel != nil {
		g.config.Logger.Debug("closing heartbeat goroutine")
		g.heartbeatCancel()
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
	g.statusMu.Lock()
	g.status = StatusDisconnected
	g.statusMu.Unlock()
}

func (g *gatewayImpl) Status() Status {
	g.statusMu.Lock()
	defer g.statusMu.Unlock()
	return g.status
}

func (g *gatewayImpl) Send(ctx context.Context, op Opcode, d GatewayMessageData) error {
	g.statusMu.Lock()
	defer g.statusMu.Unlock()
	if g.status != StatusReady {
		return discord.ErrShardNotReady
	}

	return g.sendInternal(ctx, op, d)
}

func (g *gatewayImpl) sendInternal(ctx context.Context, op Opcode, d GatewayMessageData) error {
	data, err := json.Marshal(GatewayMessage{
		Op: op,
		D:  d,
	})
	if err != nil {
		return err
	}
	return g.send(ctx, websocket.TextMessage, data)
}

func (g *gatewayImpl) send(ctx context.Context, messageType int, data []byte) error {
	g.connMu.Lock()
	defer g.connMu.Unlock()
	if g.conn == nil {
		return ErrGatewayNotConnected
	}

	g.config.Logger.DebugContext(ctx, "sending voice gateway command", slog.String("data", string(data)))
	return g.conn.WriteMessage(messageType, data)
}

func (g *gatewayImpl) Latency() time.Duration {
	return g.lastHeartbeatReceived.Sub(g.lastHeartbeatSent)
}

func (g *gatewayImpl) doReconnect(ctx context.Context, state State) error {
	var (
		try              int
		backoffIncrement int
	)

	for {
		// Exponentially backoff up to a limit of 10s
		delay := time.Duration(1<<backoffIncrement) * time.Second
		if delay > maximumConnectDelay {
			delay = maximumConnectDelay
		} else {
			backoffIncrement++
		}

		timer := time.NewTimer(time.Duration(try) * delay)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
		}

		err := g.open(ctx, state)
		if err == nil {
			// Successfully connected, our job here is done
			return nil
		}

		if errors.Is(err, discord.ErrGatewayAlreadyConnected) {
			return err
		}
		g.config.Logger.Error("failed to reconnect voice gateway", slog.Any("err", err), slog.Int("try", try), slog.Duration("delay", delay))
		g.statusMu.Lock()
		g.status = StatusDisconnected
		g.statusMu.Unlock()

		try++
	}
}

func (g *gatewayImpl) reconnect() {
	if err := g.doReconnect(context.Background(), g.state); err != nil {
		g.config.Logger.Error("failed to reopen voice gateway", slog.Any("err", err))

		if g.closeHandlerFunc != nil {
			g.closeHandlerFunc(g, err, false)
		}
	}
}

func (g *gatewayImpl) heartbeat() {
	defer g.config.Logger.Debug("exiting voice heartbeat goroutine")

	ctx, cancel := context.WithCancel(context.Background())
	g.heartbeatCancel = cancel

	g.sendHeartbeat()

	heartbeatTicker := time.NewTicker(g.heartbeatInterval)
	for {
		select {
		case <-ctx.Done():
			return

		case <-heartbeatTicker.C:
			if g.lastHeartbeatSent.After(g.lastHeartbeatReceived) {
				g.config.Logger.Warn("ACK of last heartbeat not received, connection went zombie")
				g.CloseWithCode(websocket.CloseServiceRestart, "heartbeat ACK not received")
				go g.reconnect()
				return
			}

			g.sendHeartbeat()
		}
	}
}

func (g *gatewayImpl) sendHeartbeat() {
	g.config.Logger.Debug("sending heartbeat")

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

func (g *gatewayImpl) identify() error {
	g.statusMu.Lock()
	g.status = StatusIdentifying
	g.statusMu.Unlock()
	g.config.Logger.Debug("sending Identify command")

	identify := GatewayMessageDataIdentify{
		GuildID:   g.state.GuildID,
		UserID:    g.state.UserID,
		SessionID: g.state.SessionID,
		Token:     g.state.Token,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := g.sendInternal(ctx, OpcodeIdentify, identify); err != nil {
		return err
	}

	g.statusMu.Lock()
	g.status = StatusWaitingForReady
	g.statusMu.Unlock()
	return nil
}

func (g *gatewayImpl) resume() error {
	g.statusMu.Lock()
	g.status = StatusResuming
	g.statusMu.Unlock()

	resume := GatewayMessageDataResume{
		GuildID:   g.state.GuildID,
		SessionID: g.state.SessionID,
		Token:     g.state.Token,
		SeqAck:    g.seq,
	}
	g.config.Logger.Debug("sending Resume command")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := g.sendInternal(ctx, OpcodeResume, resume); err != nil {
		return err
	}
	return nil
}

func (g *gatewayImpl) listen(conn *websocket.Conn, ready func(error)) {
	defer g.config.Logger.Debug("exiting listen goroutine")

	// Ensure that we never leave this function without calling ready
	defer ready(nil)

	for {
		mt, reader, err := conn.NextReader()
		if err != nil {
			g.statusMu.Lock()
			if g.status != StatusReady {
				g.statusMu.Unlock()
				ready(err)
				return
			}
			g.statusMu.Unlock()
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

				msg := "voice gateway close received"
				args := []any{
					slog.Bool("reconnect", reconnect),
					slog.Int("code", closeError.Code),
					slog.String("error", closeError.Text),
				}
				if reconnect {
					g.config.Logger.Warn(msg, args...)
				} else {
					g.config.Logger.Error(msg, args...)
				}
			} else if errors.Is(err, net.ErrClosed) {
				// we closed the connection ourselves. Don't try to reconnect here
				reconnect = false
			} else {
				g.config.Logger.Warn("failed to read next message from voice gateway", slog.Any("err", err))
			}

			// make sure the connection is properly closed
			g.CloseWithCode(websocket.CloseServiceRestart, "reconnecting")
			if g.config.AutoReconnect && reconnect {
				go g.reconnect()
			} else if g.closeHandlerFunc != nil {
				go g.closeHandlerFunc(g, err, reconnect)
			}

			return
		}

		message, err := g.parseMessage(mt, reader)
		if err != nil {
			g.config.Logger.Error("error while parsing voice gateway event", slog.Any("err", err))
			continue
		}

		if message.Seq > 0 {
			g.seq = message.Seq
		}

		switch message.Op {
		case OpcodeHello:
			d := message.D.(GatewayMessageDataHello)
			g.heartbeatInterval = time.Duration(d.HeartbeatInterval) * time.Millisecond
			g.lastHeartbeatReceived = time.Now().UTC()
			go g.heartbeat()

			if g.ssrc == 0 || g.seq == 0 {
				err = g.identify()
			} else {
				err = g.resume()
			}
			if err != nil {
				ready(err)
				return
			}

		case OpcodeReady:
			g.statusMu.Lock()
			g.status = StatusReady
			g.statusMu.Unlock()
			d := message.D.(GatewayMessageDataReady)
			g.ssrc = d.SSRC
			ready(nil)

		case OpcodeResumed:
			g.statusMu.Lock()
			g.status = StatusReady
			g.statusMu.Unlock()
			ready(nil)

		case OpcodeHeartbeatACK:
			d := message.D.(GatewayMessageDataHeartbeatACK)
			if d.T != g.lastNonce {
				g.config.Logger.Error("received heartbeat ack with nonce", slog.Int64("nonce", d.T), slog.Int64("last_nonce", g.lastNonce))
				go g.reconnect()
				return
			}
			g.lastHeartbeatReceived = time.Now().UTC()
		}

		g.eventHandlerFunc(g, message.Op, message.Seq, message.D)
	}
}

func (g *gatewayImpl) parseMessage(mt int, r io.Reader) (GatewayMessage, error) {
	if mt != websocket.TextMessage {
		return GatewayMessage{}, fmt.Errorf("unsupported message type: %d", mt)
	}

	if g.config.Logger.Enabled(context.Background(), slog.LevelDebug) {
		buff := new(bytes.Buffer)
		tr := io.TeeReader(r, buff)
		data, err := io.ReadAll(tr)
		if err != nil {
			return GatewayMessage{}, fmt.Errorf("failed to read message: %w", err)
		}
		g.config.Logger.Debug("received voice gateway message", slog.String("data", string(data)))
		r = buff
	}

	var message GatewayMessage
	return message, json.NewDecoder(r).Decode(&message)
}
