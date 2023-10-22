package gateway

import (
	"bytes"
	"compress/zlib"
	"context"
	"errors"
	"fmt"
	"io"
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

	conn          *websocket.Conn
	connMu        sync.Mutex
	heartbeatChan chan struct{}
	status        Status

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

func (g *gatewayImpl) formatLogsf(format string, a ...any) string {
	if g.config.ShardCount > 1 {
		return fmt.Sprintf("[%d/%d] %s", g.config.ShardID, g.config.ShardCount, fmt.Sprintf(format, a...))
	}
	return fmt.Sprintf(format, a...)
}

func (g *gatewayImpl) formatLogs(a ...any) string {
	if g.config.ShardCount > 1 {
		return fmt.Sprintf("[%d/%d] %s", g.config.ShardID, g.config.ShardCount, fmt.Sprint(a...))
	}
	return fmt.Sprint(a...)
}

func (g *gatewayImpl) Open(ctx context.Context) error {
	return g.reconnectTry(ctx, 0)
}

func (g *gatewayImpl) open(ctx context.Context) error {
	g.config.Logger.Debug(g.formatLogs("opening gateway connection"))

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
		g.Close(ctx)
		body := "empty"
		if rs != nil && rs.Body != nil {
			defer func() {
				_ = rs.Body.Close()
			}()
			rawBody, bErr := io.ReadAll(rs.Body)
			if bErr != nil {
				g.config.Logger.Error(g.formatLogs("error while reading response body: ", err))
			}
			body = string(rawBody)
		}

		g.config.Logger.Error(g.formatLogsf("error connecting to the gateway. url: %s, error: %s, body: %s", gatewayURL, err, body))
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
	if g.heartbeatChan != nil {
		g.config.Logger.Debug(g.formatLogs("closing heartbeat goroutines..."))
		g.heartbeatChan <- struct{}{}
	}

	g.connMu.Lock()
	defer g.connMu.Unlock()
	if g.conn != nil {
		g.config.RateLimiter.Close(ctx)
		g.config.Logger.Debug(g.formatLogsf("closing gateway connection with code: %d, message: %s", code, message))
		if err := g.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(code, message)); err != nil && !errors.Is(err, websocket.ErrCloseSent) {
			g.config.Logger.Debug(g.formatLogs("error writing close code. error: ", err))
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
	g.config.Logger.Trace(g.formatLogs("sending gateway command: ", string(data)))
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
		g.config.Logger.Error(g.formatLogs("failed to reconnect gateway. error: ", err))
		g.status = StatusDisconnected
		return g.reconnectTry(ctx, try+1)
	}
	return nil
}

func (g *gatewayImpl) reconnect() {
	err := g.reconnectTry(context.Background(), 0)
	if err != nil {
		g.config.Logger.Error(g.formatLogs("failed to reopen gateway. error: ", err))
	}
}

func (g *gatewayImpl) heartbeat() {
	if g.heartbeatChan == nil {
		g.heartbeatChan = make(chan struct{})
	}
	heartbeatTicker := time.NewTicker(g.heartbeatInterval)
	defer heartbeatTicker.Stop()
	defer g.config.Logger.Debug(g.formatLogs("exiting heartbeat goroutine..."))

	for {
		select {
		case <-g.heartbeatChan:
			return

		case <-heartbeatTicker.C:
			g.sendHeartbeat()
		}
	}
}

func (g *gatewayImpl) sendHeartbeat() {
	g.config.Logger.Debug(g.formatLogs("sending heartbeat..."))

	ctx, cancel := context.WithTimeout(context.Background(), g.heartbeatInterval)
	defer cancel()
	if err := g.Send(ctx, OpcodeHeartbeat, MessageDataHeartbeat(*g.config.LastSequenceReceived)); err != nil {
		if errors.Is(err, discord.ErrShardNotConnected) || errors.Is(err, syscall.EPIPE) {
			return
		}
		g.config.Logger.Error(g.formatLogs("failed to send heartbeat. error: ", err))
		g.CloseWithCode(context.TODO(), websocket.CloseServiceRestart, "heartbeat timeout")
		go g.reconnect()
		return
	}
	g.lastHeartbeatSent = time.Now().UTC()
}

func (g *gatewayImpl) identify() {
	g.status = StatusIdentifying
	g.config.Logger.Debug(g.formatLogs("sending Identify command..."))

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

	if err := g.Send(context.TODO(), OpcodeIdentify, identify); err != nil {
		g.config.Logger.Error(g.formatLogs("error sending Identify command err: ", err))
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

	g.config.Logger.Debug(g.formatLogs("sending Resume command..."))
	if err := g.Send(context.TODO(), OpcodeResume, resume); err != nil {
		g.config.Logger.Error(g.formatLogs("error sending resume command err: ", err))
	}
}

func (g *gatewayImpl) listen(conn *websocket.Conn) {
	defer g.config.Logger.Debug(g.formatLogs("exiting listen goroutine..."))
loop:
	for {
		mt, data, err := conn.ReadMessage()
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
				message := g.formatLogsf("gateway close received, reconnect: %t, code: %d, error: %s", g.config.AutoReconnect && reconnect, closeError.Code, closeError.Text)
				if reconnect {
					g.config.Logger.Debug(message)
				} else {
					g.config.Logger.Error(message)
				}
			} else if errors.Is(err, net.ErrClosed) {
				// we closed the connection ourselves. Don't try to reconnect here
				reconnect = false
			} else {
				g.config.Logger.Debug(g.formatLogs("failed to read next message from gateway. error: ", err))
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

		message, err := g.parseMessage(mt, data)
		if err != nil {
			g.config.Logger.Error(g.formatLogs("error while parsing gateway message. error: ", err))
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
				g.config.Logger.Error(g.formatLogsf("invalid message data of type %T received", message.D))
				continue
			}

			// get session id here
			if readyEvent, ok := eventData.(EventReady); ok {
				g.config.SessionID = &readyEvent.SessionID
				g.config.ResumeURL = &readyEvent.ResumeGatewayURL
				g.status = StatusReady
				g.config.Logger.Debug(g.formatLogs("ready message received"))
			}

			if unknownEvent, ok := eventData.(EventUnknown); ok {
				g.config.Logger.Debug(g.formatLogsf("unknown event received: %s, data: %s", message.T, unknownEvent))
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
			g.config.Logger.Debug(g.formatLogsf("unknown opcode received: %d, data: %s", message.Op, message.D))
		}
	}
}

func (g *gatewayImpl) parseMessage(mt int, data []byte) (Message, error) {
	var finalData []byte
	if mt == websocket.BinaryMessage {
		g.config.Logger.Trace(g.formatLogs("binary message received. decompressing..."))

		reader, err := zlib.NewReader(bytes.NewReader(data))
		if err != nil {
			return Message{}, fmt.Errorf("failed to decompress zlib: %w", err)
		}
		defer reader.Close()
		finalData, err = io.ReadAll(reader)
		if err != nil {
			return Message{}, fmt.Errorf("failed to read decompressed data: %w", err)
		}
	} else {
		finalData = data
	}

	g.config.Logger.Trace(g.formatLogs("received gateway message: ", string(finalData)))

	var message Message
	return message, json.Unmarshal(finalData, &message)
}
