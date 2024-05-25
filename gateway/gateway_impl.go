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

	"github.com/disgoorg/json"
	"github.com/gorilla/websocket"

	"github.com/disgoorg/disgo/discord"
)

var _ Gateway = (*gatewayImpl)(nil)

// New creates a new Gateway instance with the provided token, eventHandlerFunc, closeHandlerFunc and ConfigOpt(s).
func New(token string, eventHandlerFunc EventHandlerFunc, closeHandlerFunc CloseHandlerFunc, opts ...ConfigOpt) Gateway {
	config := DefaultConfig()
	config.Apply(opts)
	config.Logger = config.Logger.With(slog.String("name", "gateway"), slog.Int("shard_id", config.ShardID), slog.Int("shard_count", config.ShardCount))

	return &gatewayImpl{
		config:           *config,
		eventHandlerFunc: eventHandlerFunc,
		closeHandlerFunc: closeHandlerFunc,
		token:            token,
		status:           StatusUnconnected,
	}
}

type gatewayImpl struct {
	config           Config
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
	defer g.connMu.Unlock()
	if g.conn != nil {
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
		return err
	}

	conn.SetCloseHandler(func(code int, text string) error {
		return nil
	})

	g.conn = conn

	// reset rate limiter when connecting
	g.config.RateLimiter.Reset()

	g.status = StatusWaitingForHello

	go g.listen(conn)

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

	ctx, cancel := context.WithTimeout(context.Background(), g.heartbeatInterval)
	defer cancel()
	if err := g.Send(ctx, OpcodeHeartbeat, MessageDataHeartbeat(*g.config.LastSequenceReceived)); err != nil {
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

func (g *gatewayImpl) identify() {
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
		g.config.Logger.Error("error sending Identify command", slog.Any("err", err))
	}
	g.status = StatusWaitingForReady
}

func (g *gatewayImpl) resume() {
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
		g.config.Logger.Error("error sending resume command", slog.Any("err", err))
	}
}

func (g *gatewayImpl) listen(conn *websocket.Conn) {
	defer g.config.Logger.Debug("exiting listen goroutine")
loop:
	for {
		mt, r, err := conn.NextReader()
		if err != nil {
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
				g.identify()
			} else {
				g.resume()
			}

		case OpcodeDispatch:
			// set last sequence received
			g.config.LastSequenceReceived = &message.S

			eventData, ok := message.D.(EventData)
			if !ok && message.D != nil {
				g.config.Logger.Error("invalid message data received", slog.String("data", fmt.Sprintf("%T", message.D)))
				continue
			}

			// get session id here
			if readyEvent, ok := eventData.(EventReady); ok {
				g.config.SessionID = &readyEvent.SessionID
				g.config.ResumeURL = &readyEvent.ResumeGatewayURL
				g.status = StatusReady
				g.config.Logger.Debug("ready message received")
			}

			if unknownEvent, ok := eventData.(EventUnknown); ok {
				g.config.Logger.Debug("unknown event received", slog.String("event", string(message.T)), slog.String("data", string(unknownEvent)))
				continue
			}

			// push message to the command manager
			if g.config.EnableRawEvents {
				g.eventHandlerFunc(EventTypeRaw, message.S, g.config.ShardID, EventRaw{
					EventType: message.T,
					Payload:   bytes.NewReader(message.RawD),
				})
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
