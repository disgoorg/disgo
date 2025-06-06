package gateway

import (
	"bytes"
	"compress/zlib"
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
	"github.com/gorilla/websocket"

	"github.com/disgoorg/disgo/discord"
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

// String returns the string representation of the Status.
func (s Status) String() string {
	switch s {
	case StatusUnconnected:
		return "Unconnected"
	case StatusConnecting:
		return "Connecting"
	case StatusWaitingForHello:
		return "WaitingForHello"
	case StatusIdentifying:
		return "Identifying"
	case StatusResuming:
		return "Resuming"
	case StatusWaitingForReady:
		return "WaitingForReady"
	case StatusReady:
		return "Ready"
	case StatusDisconnected:
		return "Disconnected"
	default:
		return "Unknown"
	}
}

// Indicates how far along the client is too connecting.
const (
	// StatusUnconnected is the initial state when a new Gateway is created.
	StatusUnconnected Status = iota

	// StatusConnecting is the state when the client is connecting to the Discord gateway.
	StatusConnecting

	// StatusWaitingForHello is the state when the Gateway is waiting for the first OpcodeHello packet.
	StatusWaitingForHello

	// StatusIdentifying is the state when the Gateway received its first OpcodeHello packet and now sends a OpcodeIdentify packet.
	StatusIdentifying

	// StatusResuming is the state when the Gateway received its first OpcodeHello packet and now sends a OpcodeResume packet.
	StatusResuming

	// StatusWaitingForReady is the state when the Gateway received sent a OpcodeIdentify or OpcodeResume packet and now waits for a OpcodeDispatch with EventTypeReady packet.
	StatusWaitingForReady

	// StatusReady is the state when the Gateway received a OpcodeDispatch with EventTypeReady packet.
	StatusReady

	// StatusDisconnected is the state when the Gateway is disconnected.
	// Either due to an error or because the Gateway was closed gracefully.
	StatusDisconnected
)

type (
	// EventHandlerFunc is a function that is called when an event is received.
	EventHandlerFunc func(gatewayEventType EventType, sequenceNumber int, shardID int, event EventData)

	// CreateFunc is a type that is used to create a new Gateway(s).
	CreateFunc func(token string, eventHandlerFunc EventHandlerFunc, closeHandlerFUnc CloseHandlerFunc, opts ...ConfigOpt) Gateway

	// CloseHandlerFunc is a function that is called when the Gateway is closed.
	CloseHandlerFunc func(gateway Gateway, err error)
)

// Gateway is what is used to connect to discord.
type Gateway interface {
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

	// Intents returns the Intents that are used by this Gateway.
	Intents() Intents

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
	Send(ctx context.Context, op Opcode, data MessageData) error

	// Latency returns the latency of the Gateway.
	// This is calculated by the time it takes to send a heartbeat and receive a heartbeat ack by discord.
	Latency() time.Duration

	// Presence returns the current presence of the Gateway.
	Presence() *MessageDataPresenceUpdate
}

var _ Gateway = (*gatewayImpl)(nil)

// New creates a new Gateway instance with the provided token, eventHandlerFunc, closeHandlerFunc and ConfigOpt(s).
func New(token string, eventHandlerFunc EventHandlerFunc, closeHandlerFunc CloseHandlerFunc, opts ...ConfigOpt) Gateway {
	cfg := defaultConfig()
	cfg.apply(opts)

	return &gatewayImpl{
		config:           cfg,
		eventHandlerFunc: eventHandlerFunc,
		closeHandlerFunc: closeHandlerFunc,
		token:            token,
		status:           StatusUnconnected,
	}
}

type gatewayImpl struct {
	config           config
	eventHandlerFunc EventHandlerFunc
	closeHandlerFunc CloseHandlerFunc
	token            string

	conn            *websocket.Conn
	connMu          sync.Mutex
	heartbeatCancel context.CancelFunc
	status          Status

	heartbeatInterval     time.Duration
	lastHeartbeatSent     time.Time
	lastHeartbeatReceived time.Time
}

func (g *gatewayImpl) ShardID() int {
	return g.config.ShardID
}

func (g *gatewayImpl) ShardCount() int {
	return g.config.ShardCount
}

func (g *gatewayImpl) SessionID() *string {
	return g.config.SessionID
}

func (g *gatewayImpl) LastSequenceReceived() *int {
	return g.config.LastSequenceReceived
}

func (g *gatewayImpl) Intents() Intents {
	return g.config.Intents
}

func (g *gatewayImpl) Open(ctx context.Context) error {
	return g.reconnectTry(ctx, 0)
}

func (g *gatewayImpl) open(ctx context.Context) error {
	g.config.Logger.Debug("opening gateway connection")

	g.connMu.Lock()
	if g.conn != nil {
		g.connMu.Unlock()
		return discord.ErrGatewayAlreadyConnected
	}
	g.status = StatusConnecting

	wsURL := g.config.URL
	if g.config.ResumeURL != nil && g.config.EnableResumeURL {
		wsURL = *g.config.ResumeURL
	}
	gatewayURL := fmt.Sprintf("%s?v=%d&encoding=json", wsURL, Version)
	g.lastHeartbeatSent = time.Now().UTC()
	conn, rs, err := g.config.Dialer.DialContext(ctx, gatewayURL, nil)
	if err != nil {
		body := ""
		if rs != nil && rs.Body != nil {
			defer func() {
				_ = rs.Body.Close()
			}()
			rawBody, bErr := io.ReadAll(rs.Body)
			if bErr != nil {
				g.config.Logger.Error("error while reading response body", slog.Any("err", bErr))
			}
			body = string(rawBody)
		}

		g.config.Logger.Error("error connecting to the gateway", slog.Any("err", err), slog.String("url", gatewayURL), slog.String("body", body))
		g.connMu.Unlock()
		return err
	}

	conn.SetCloseHandler(func(code int, text string) error {
		return nil
	})

	g.conn = conn
	g.connMu.Unlock()

	// reset rate limiter when connecting
	g.config.RateLimiter.Reset()

	g.status = StatusWaitingForHello

	readyChan := make(chan error)
	go g.listen(conn, readyChan)

	select {
	case <-ctx.Done():
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		g.Close(closeCtx)
		return ctx.Err()
	case err = <-readyChan:
		if err != nil {
			closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			g.Close(closeCtx)
			return fmt.Errorf("failed to open gateway connection: %w", err)
		}
	}

	return nil
}

func (g *gatewayImpl) Close(ctx context.Context) {
	g.CloseWithCode(ctx, websocket.CloseNormalClosure, "Shutting down")
}

func (g *gatewayImpl) CloseWithCode(ctx context.Context, code int, message string) {
	if g.heartbeatCancel != nil {
		g.config.Logger.Debug("closing heartbeat goroutines...")
		g.heartbeatCancel()
	}

	g.connMu.Lock()
	defer g.connMu.Unlock()
	if g.conn != nil {
		g.config.RateLimiter.Close(ctx)
		g.config.Logger.Debug("closing gateway connection", slog.Int("code", code), slog.String("message", message))
		if err := g.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(code, message)); err != nil && !errors.Is(err, websocket.ErrCloseSent) {
			g.config.Logger.Debug("error writing close code", slog.Any("err", err))
		}
		_ = g.conn.Close()
		g.conn = nil

		// clear resume data as we closed gracefully
		if code == websocket.CloseNormalClosure || code == websocket.CloseGoingAway {
			g.config.SessionID = nil
			g.config.ResumeURL = nil
			g.config.LastSequenceReceived = nil
		}
	}
	g.status = StatusDisconnected
}

func (g *gatewayImpl) Status() Status {
	g.connMu.Lock()
	defer g.connMu.Unlock()
	return g.status
}

func (g *gatewayImpl) Send(ctx context.Context, op Opcode, d MessageData) error {
	data, err := json.Marshal(Message{
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
		return discord.ErrShardNotConnected
	}

	if err := g.config.RateLimiter.Wait(ctx); err != nil {
		return err
	}

	defer g.config.RateLimiter.Unlock()
	if g.config.Logger.Enabled(ctx, slog.LevelDebug) {
		g.config.Logger.Debug("sending gateway command", slog.String("data", string(data)))
	}
	return g.conn.WriteMessage(messageType, data)
}

func (g *gatewayImpl) Latency() time.Duration {
	return g.lastHeartbeatReceived.Sub(g.lastHeartbeatSent)
}

func (g *gatewayImpl) Presence() *MessageDataPresenceUpdate {
	return g.config.Presence
}

func (g *gatewayImpl) reconnectTry(ctx context.Context, try int) error {
	delay := time.Duration(try) * 2 * time.Second
	if delay > 30*time.Second {
		delay = 30 * time.Second
	}

	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		timer.Stop()
		return ctx.Err()
	case <-timer.C:
	}

	if err := g.open(ctx); err != nil {
		var closeError *websocket.CloseError
		if errors.As(err, &closeError) {
			closeCode := CloseEventCodeByCode(closeError.Code)
			if !closeCode.Reconnect {
				return err
			}
		}
		if errors.Is(err, discord.ErrGatewayAlreadyConnected) {
			return err
		}
		g.config.Logger.Error("failed to reconnect gateway", slog.Any("err", err))
		g.status = StatusDisconnected
		return g.reconnectTry(ctx, try+1)
	}
	return nil
}

func (g *gatewayImpl) reconnect() {
	err := g.reconnectTry(context.Background(), 0)
	if err != nil {
		g.config.Logger.Error("failed to reopen gateway", slog.Any("err", err))
	}
}

func (g *gatewayImpl) heartbeat() {
	ctx, cancel := context.WithCancel(context.Background())
	g.heartbeatCancel = cancel

	heartbeatTicker := time.NewTicker(g.heartbeatInterval)
	defer heartbeatTicker.Stop()
	defer g.config.Logger.Debug("exiting heartbeat goroutine")

	for {
		select {
		case <-ctx.Done():
			return

		case <-heartbeatTicker.C:
			g.sendHeartbeat()
		}
	}
}

func (g *gatewayImpl) sendHeartbeat() {
	g.config.Logger.Debug("sending heartbeat")

	sequence := 0
	if g.config.LastSequenceReceived != nil {
		sequence = *g.config.LastSequenceReceived
	}

	ctx, cancel := context.WithTimeout(context.Background(), g.heartbeatInterval)
	defer cancel()
	if err := g.Send(ctx, OpcodeHeartbeat, MessageDataHeartbeat(sequence)); err != nil {
		if errors.Is(err, discord.ErrShardNotConnected) || errors.Is(err, syscall.EPIPE) {
			return
		}
		g.config.Logger.Error("failed to send heartbeat", slog.Any("err", err))
		closeCtx, closeCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer closeCancel()
		g.CloseWithCode(closeCtx, websocket.CloseServiceRestart, "heartbeat timeout")
		go g.reconnect()
		return
	}
	g.lastHeartbeatSent = time.Now().UTC()
}

func (g *gatewayImpl) identify() error {
	g.status = StatusIdentifying
	g.config.Logger.Debug("sending Identify command")

	identify := MessageDataIdentify{
		Token: g.token,
		Properties: IdentifyCommandDataProperties{
			OS:      g.config.OS,
			Browser: g.config.Browser,
			Device:  g.config.Device,
		},
		Compress:       g.config.Compress,
		LargeThreshold: g.config.LargeThreshold,
		Intents:        g.config.Intents,
		Presence:       g.config.Presence,
		Shard:          &[2]int{g.ShardID(), g.ShardCount()},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := g.Send(ctx, OpcodeIdentify, identify); err != nil {
		return err
	}
	g.status = StatusWaitingForReady
	return nil
}

func (g *gatewayImpl) resume() error {
	g.status = StatusResuming
	resume := MessageDataResume{
		Token:     g.token,
		SessionID: *g.config.SessionID,
		Seq:       *g.config.LastSequenceReceived,
	}
	g.config.Logger.Debug("sending Resume command")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := g.Send(ctx, OpcodeResume, resume); err != nil {
		return err
	}
	return nil
}

func (g *gatewayImpl) listen(conn *websocket.Conn, readyChan chan<- error) {
	defer g.config.Logger.Debug("exiting listen goroutine")
loop:
	for {
		mt, r, err := conn.NextReader()
		if err != nil {
			if g.status != StatusReady {
				readyChan <- err
				close(readyChan)
				break loop
			}
			g.connMu.Lock()
			sameConnection := g.conn == conn
			g.connMu.Unlock()

			// if sameConnection is false, it means the connection has been closed by the user, and we can just exit
			if !sameConnection {
				return
			}

			reconnect := true
			var closeError *websocket.CloseError
			if errors.As(err, &closeError) {
				closeCode := CloseEventCodeByCode(closeError.Code)
				reconnect = closeCode.Reconnect

				if closeCode == CloseEventCodeInvalidSeq {
					g.config.LastSequenceReceived = nil
					g.config.SessionID = nil
					g.config.ResumeURL = nil
				}
				msg := "gateway close received"
				args := []any{
					slog.Bool("reconnect", reconnect),
					slog.Int("code", closeError.Code),
					slog.String("error", closeError.Text),
				}
				if reconnect {
					g.config.Logger.Debug(msg, args...)
				} else {
					g.config.Logger.Error(msg, args...)
				}
			} else if errors.Is(err, net.ErrClosed) {
				// we closed the connection ourselves. Don't try to reconnect here
				reconnect = false
			} else {
				g.config.Logger.Debug("failed to read next message from gateway", slog.Any("err", err))
			}

			// make sure the connection is properly closed
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			g.CloseWithCode(ctx, websocket.CloseServiceRestart, "reconnecting")
			cancel()
			if g.config.AutoReconnect && reconnect {
				go g.reconnect()
			} else if g.closeHandlerFunc != nil {
				go g.closeHandlerFunc(g, err)
			}

			break loop
		}

		message, err := g.parseMessage(mt, r)
		if err != nil {
			g.config.Logger.Error("error while parsing gateway message", slog.Any("err", err))
			continue
		}

		switch message.Op {
		case OpcodeHello:
			g.heartbeatInterval = time.Duration(message.D.(MessageDataHello).HeartbeatInterval) * time.Millisecond
			g.lastHeartbeatReceived = time.Now().UTC()
			go g.heartbeat()

			if g.config.LastSequenceReceived == nil || g.config.SessionID == nil {
				err = g.identify()
			} else {
				err = g.resume()
			}
			if err != nil {
				readyChan <- err
				close(readyChan)
				return
			}

		case OpcodeDispatch:
			// set last sequence received
			g.config.LastSequenceReceived = &message.S

			eventData, ok := message.D.(EventData)
			if !ok && message.D != nil {
				g.config.Logger.Error("invalid message data received", slog.String("data", fmt.Sprintf("%T", message.D)))
				continue
			}

			if readyEvent, ok := eventData.(EventReady); ok {
				g.config.SessionID = &readyEvent.SessionID
				g.config.ResumeURL = &readyEvent.ResumeGatewayURL
				g.config.Logger.Debug("ready message received")
				g.status = StatusReady
				readyChan <- nil
				close(readyChan)
			} else if _, ok = eventData.(EventResumed); ok {
				g.config.Logger.Debug("resume message received")
				g.status = StatusReady
				readyChan <- nil
				close(readyChan)
			}

			// push message to the command manager
			if g.config.EnableRawEvents {
				g.eventHandlerFunc(EventTypeRaw, message.S, g.config.ShardID, EventRaw{
					EventType: message.T,
					Payload:   bytes.NewReader(message.RawD),
				})
			}

			if unknownEvent, ok := eventData.(EventUnknown); ok {
				g.config.Logger.Debug("unknown event received", slog.String("event", string(message.T)), slog.String("data", string(unknownEvent)))
				continue
			}
			g.eventHandlerFunc(message.T, message.S, g.config.ShardID, eventData)

		case OpcodeHeartbeat:
			g.sendHeartbeat()

		case OpcodeReconnect:
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			g.CloseWithCode(ctx, websocket.CloseServiceRestart, "received reconnect")
			cancel()
			go g.reconnect()
			break loop

		case OpcodeInvalidSession:
			canResume := message.D.(MessageDataInvalidSession)

			code := websocket.CloseNormalClosure
			if canResume {
				code = websocket.CloseServiceRestart
			} else {
				// clear resume info
				g.config.SessionID = nil
				g.config.LastSequenceReceived = nil
				g.config.ResumeURL = nil
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			g.CloseWithCode(ctx, code, "invalid session")
			cancel()
			go g.reconnect()
			break loop

		case OpcodeHeartbeatACK:
			newHeartbeat := time.Now().UTC()
			g.eventHandlerFunc(EventTypeHeartbeatAck, message.S, g.config.ShardID, EventHeartbeatAck{
				LastHeartbeat: g.lastHeartbeatReceived,
				NewHeartbeat:  newHeartbeat,
			})
			g.lastHeartbeatReceived = newHeartbeat

		default:

			g.config.Logger.Debug("unknown opcode received", slog.Int("opcode", int(message.Op)), slog.String("data", fmt.Sprintf("%s", message.D)))
		}
	}
}

func (g *gatewayImpl) parseMessage(mt int, r io.Reader) (Message, error) {
	if mt == websocket.BinaryMessage {
		g.config.Logger.Debug("binary message received. decompressing")

		reader, err := zlib.NewReader(r)
		if err != nil {
			return Message{}, fmt.Errorf("failed to decompress zlib: %w", err)
		}
		defer reader.Close()
		r = reader
	}

	if g.config.Logger.Enabled(context.Background(), slog.LevelDebug) {
		buff := new(bytes.Buffer)
		tr := io.TeeReader(r, buff)
		data, err := io.ReadAll(tr)
		if err != nil {
			return Message{}, fmt.Errorf("failed to read message: %w", err)
		}
		g.config.Logger.Debug("received gateway message", slog.String("data", string(data)))
		r = buff
	}

	var message Message
	return message, json.NewDecoder(r).Decode(&message)
}
