package gateway

import (
	"bytes"
	"compress/zlib"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/internal/tokenhelper"
	"github.com/disgoorg/disgo/json"
	"github.com/disgoorg/log"

	"github.com/gorilla/websocket"
)

var _ Gateway = (*gatewayImpl)(nil)

func New(token string, eventHandlerFunc EventHandlerFunc, opts ...ConfigOpt) Gateway {
	config := DefaultConfig()
	config.Apply(opts)

	return &gatewayImpl{
		config:           *config,
		eventHandlerFunc: eventHandlerFunc,
		token:            token,
		status:           StatusUnconnected,
	}
}

type gatewayImpl struct {
	config           Config
	eventHandlerFunc EventHandlerFunc
	token            string

	conn            *websocket.Conn
	heartbeatTicker *time.Ticker
	status          Status

	heartbeatInterval     time.Duration
	lastHeartbeatSent     time.Time
	lastHeartbeatReceived time.Time
}

func (g *gatewayImpl) Logger() log.Logger {
	return g.config.Logger
}

func (g *gatewayImpl) ShardID() int {
	return g.config.ShardID
}

func (g *gatewayImpl) ShardCount() int {
	return g.config.ShardCount
}

func (g *gatewayImpl) GatewayIntents() discord.GatewayIntents {
	return g.config.GatewayIntents
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
	g.Logger().Debug(g.formatLogs("opening gateway connection"))
	if g.conn != nil {
		return discord.ErrGatewayAlreadyConnected
	}
	g.status = StatusConnecting

	gatewayURL := g.config.GatewayURL + "?v=" + APIVersion + "&encoding=json"
	var rs *http.Response
	var err error
	g.conn, rs, err = websocket.DefaultDialer.DialContext(ctx, gatewayURL, nil)
	if err != nil {
		g.Close(ctx)
		body := "null"
		if rs != nil && rs.Body != nil {
			defer rs.Body.Close()
			rawBody, bErr := ioutil.ReadAll(rs.Body)
			if bErr != nil {
				g.Logger().Error(g.formatLogs("error while reading response body: ", err))
			}
			body = string(rawBody)
		}

		g.Logger().Error(g.formatLogsf("error connecting to the gateway. url: %s, error: %s, body: %s", gatewayURL, err, body))
		return err
	}

	g.status = StatusWaitingForHello

	go g.listen()

	return nil
}

func (g *gatewayImpl) Close(ctx context.Context) {
	g.CloseWithCode(ctx, websocket.CloseNormalClosure, "Shutting down")
}

func (g *gatewayImpl) CloseWithCode(ctx context.Context, code int, message string) {
	g.Logger().Debug(g.formatLogsf("closing gateway connection with code: %d, message: %s", code, message))
	_ = g.config.RateLimiter.Close(ctx)

	if g.heartbeatTicker != nil {
		g.Logger().Debug(g.formatLogs("closing heartbeat goroutines..."))
		g.heartbeatTicker.Stop()
		g.heartbeatTicker = nil
	}

	if g.conn != nil {
		err := g.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(code, message))
		if err != nil && err != websocket.ErrCloseSent {
			g.Logger().Error(g.formatLogs("error writing close code. error: ", err))
		}

		err = g.conn.Close()
		if err != nil {
			g.Logger().Error(g.formatLogs("error closing conn: ", err))
		}
		g.conn = nil
	}
}

func (g *gatewayImpl) Status() Status {
	return g.status
}

func (g *gatewayImpl) Send(ctx context.Context, op discord.GatewayOpcode, d discord.GatewayMessageData) error {
	if g.conn == nil {
		return discord.ErrShardNotConnected
	}

	if err := g.config.RateLimiter.Wait(ctx); err != nil {
		return err
	}

	defer g.config.RateLimiter.Unlock()
	data, err := json.Marshal(discord.GatewayMessage{
		Op: op,
		D:  d,
	})
	if err != nil {
		return err
	}
	g.Logger().Trace(g.formatLogs("sending gateway command: ", string(data)))
	return g.conn.WriteMessage(websocket.TextMessage, data)
}

func (g *gatewayImpl) Latency() time.Duration {
	return g.lastHeartbeatReceived.Sub(g.lastHeartbeatSent)
}

func (g *gatewayImpl) ReOpen(ctx context.Context, delay time.Duration) error {
	return g.reOpen(ctx, 0, delay)
}

func (g *gatewayImpl) reOpen(ctx context.Context, try int, delay time.Duration) error {
	if try >= g.config.MaxReconnectTries-1 {
		return fmt.Errorf("failed to reconnect. exceeded max reconnect tries of %d reached", g.config.MaxReconnectTries)
	}
	timer := time.NewTimer(time.Duration(try) * delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		timer.Stop()
		return ctx.Err()
	case <-timer.C:
	}

	g.Close(ctx)

	g.Logger().Debug(g.formatLogs("reconnecting gateway..."))
	if err := g.Open(ctx); err != nil {
		if err == discord.ErrGatewayAlreadyConnected {
			return err
		}
		g.Logger().Error(g.formatLogs("failed to reconnect gateway. error: ", err))
		g.status = StatusDisconnected
		return g.reOpen(ctx, try+1, delay)
	}
	return nil
}

func (g *gatewayImpl) reconnect(ctx context.Context) {
	err := g.ReOpen(ctx, time.Second)
	if ctx.Err() != nil {
		g.Logger().Error(g.formatLogs("failed to reopen gateway", err))
	}
}

func (g *gatewayImpl) heartbeat() {
	g.heartbeatTicker = time.NewTicker(g.heartbeatInterval)
	defer g.heartbeatTicker.Stop()
	defer g.Logger().Debug(g.formatLogs("exiting heartbeat goroutine..."))

	for range g.heartbeatTicker.C {
		g.sendHeartbeat()
	}
}

func (g *gatewayImpl) sendHeartbeat() {
	g.Logger().Debug(g.formatLogs("sending heartbeat..."))

	ctx, cancel := context.WithTimeout(context.Background(), g.heartbeatInterval)
	defer cancel()
	if err := g.Send(ctx, discord.GatewayOpcodeHeartbeat, (*discord.GatewayMessageDataHeartbeat)(g.config.LastSequenceReceived)); err != nil && err != discord.ErrShardNotConnected {
		g.Logger().Error(g.formatLogs("failed to send heartbeat. error: ", err))
		g.CloseWithCode(context.TODO(), websocket.CloseServiceRestart, "heartbeat timeout")
		go g.reconnect(context.TODO())
		return
	}
	g.lastHeartbeatSent = time.Now().UTC()
}

func (g *gatewayImpl) connect() {
	g.status = StatusIdentifying
	g.Logger().Debug(g.formatLogs("sending Identify command..."))

	identify := discord.GatewayMessageDataIdentify{
		Token: g.token,
		Properties: discord.IdentifyCommandDataProperties{
			OS:      g.config.OS,
			Browser: g.config.Browser,
			Device:  g.config.Device,
		},
		Compress:       g.config.Compress,
		LargeThreshold: g.config.LargeThreshold,
		GatewayIntents: g.config.GatewayIntents,
		Presence:       g.config.Presence,
	}
	if g.ShardCount() > 1 {
		identify.Shard = [2]int{g.ShardID(), g.ShardCount()}
	}

	if err := g.Send(context.TODO(), discord.GatewayOpcodeIdentify, identify); err != nil {
		g.Logger().Error(g.formatLogs("error sending Identify command err: ", err))
	}
	g.status = StatusWaitingForReady
}

func (g *gatewayImpl) resume() {
	g.status = StatusResuming
	resume := discord.GatewayMessageDataResume{
		Token:     g.token,
		SessionID: *g.config.SessionID,
		Seq:       *g.config.LastSequenceReceived,
	}

	g.Logger().Info(g.formatLogs("sending Resume command..."))
	if err := g.Send(context.TODO(), discord.GatewayOpcodeResume, resume); err != nil {
		g.Logger().Error(g.formatLogs("error sending resume command err: ", err))
	}
}

func (g *gatewayImpl) listen() {
	defer g.Logger().Debug(g.formatLogs("exiting listen goroutine..."))
	for {
		if g.conn == nil {
			return
		}
		mt, reader, err := g.conn.NextReader()
		if err != nil {
			reconnect := true
			if closeError, ok := err.(*websocket.CloseError); ok {
				closeCode := discord.GatewayCloseEventCode(closeError.Code)
				reconnect = closeCode.ShouldReconnect()

				if closeCode == discord.GatewayCloseEventCodeDisallowedIntents {
					var intentsURL string
					if id, err := tokenhelper.IDFromToken(g.token); err == nil {
						intentsURL = fmt.Sprintf("https://discord.com/developers/applications/%s/bot", *id)
					} else {
						intentsURL = "https://discord.com/developers/applications"
					}
					g.Logger().Error(g.formatLogsf("disallowed gateway intents supplied. go to %s and enable the privileged intent for your application. intents: %d", intentsURL, g.config.GatewayIntents))
				} else if closeCode == discord.GatewayCloseEventCodeInvalidSeq {
					g.Logger().Error(g.formatLogs("invalid sequence provided. reconnecting..."))
					g.config.LastSequenceReceived = nil
					g.config.SessionID = nil
				} else {
					g.Logger().Error(g.formatLogsf("gateway close received, reconnect: %t, code: %d, error: %s", reconnect && g.config.AutoReconnect, closeError.Code, closeError.Text))
				}
			} else {
				g.Logger().Error(g.formatLogs("failed to read next message from gateway. error: ", err))
			}

			if g.config.AutoReconnect && reconnect {
				go g.reconnect(context.TODO())
			} else {
				g.Close(context.TODO())
			}
			return
		}

		event, err := g.parseGatewayMessage(mt, reader)
		if err != nil {
			g.Logger().Error(g.formatLogs("error while parsing gateway event. error: ", err))
			continue
		}

		switch event.Op {
		case discord.GatewayOpcodeHello:
			g.lastHeartbeatReceived = time.Now().UTC()
			g.lastHeartbeatSent = time.Now().UTC()

			g.heartbeatInterval = time.Duration(event.D.(discord.GatewayMessageDataHello).HeartbeatInterval) * time.Millisecond

			if g.config.LastSequenceReceived == nil || g.config.SessionID == nil {
				g.connect()
			} else {
				g.resume()
			}
			go g.heartbeat()

		case discord.GatewayOpcodeDispatch:
			data := event.D.(discord.GatewayMessageDataDispatch)
			g.Logger().Trace(g.formatLogsf("received: OpcodeDispatch %s, data: %s", event.T, string(data)))

			// set last sequence received
			g.config.LastSequenceReceived = &event.S

			// get session id here
			if event.T == discord.GatewayEventTypeReady {
				var readyEvent discord.GatewayEventReady
				if err = json.Unmarshal(data, &readyEvent); err != nil {
					g.Logger().Error(g.formatLogs("Error parsing ready event. error: ", err))
					continue
				}
				g.config.SessionID = &readyEvent.SessionID
				g.status = StatusReady
				g.Logger().Debug(g.formatLogs("ready event received"))
			}

			// push event to the command manager
			g.eventHandlerFunc(event.T, event.S, bytes.NewBuffer(data))

		case discord.GatewayOpcodeHeartbeat:
			g.Logger().Debug(g.formatLogs("received: OpcodeHeartbeat"))
			g.sendHeartbeat()

		case discord.GatewayOpcodeReconnect:
			g.Logger().Debug(g.formatLogs("received: OpcodeReconnect"))
			g.CloseWithCode(context.TODO(), websocket.CloseServiceRestart, "received reconnect")
			go g.reconnect(context.TODO())

		case discord.GatewayOpcodeInvalidSession:
			canResume := event.D.(discord.GatewayMessageDataInvalidSession)
			g.Logger().Debug(g.formatLogs("received: OpcodeInvalidSession, canResume: ", canResume))
			if canResume {
				g.CloseWithCode(context.TODO(), websocket.CloseServiceRestart, "invalid session")
			} else {
				g.Close(context.TODO())
				// clear resume info
				g.config.SessionID = nil
				g.config.LastSequenceReceived = nil
			}
			go g.reconnect(context.TODO())

		case discord.GatewayOpcodeHeartbeatACK:
			g.Logger().Debug(g.formatLogs("received: OpcodeHeartbeatACK"))
			g.lastHeartbeatReceived = time.Now().UTC()
		}
	}
}

func (g *gatewayImpl) parseGatewayMessage(mt int, reader io.Reader) (*discord.GatewayMessage, error) {
	var finalReadCloser io.ReadCloser
	if mt == websocket.BinaryMessage {
		g.Logger().Trace(g.formatLogs("binary message received. decompressing..."))
		readCloser, err := zlib.NewReader(reader)
		if err != nil {
			return nil, fmt.Errorf("failed to decompress zlib: %w", err)
		}
		finalReadCloser = readCloser
	} else {
		finalReadCloser = io.NopCloser(reader)
	}
	defer finalReadCloser.Close()

	var message discord.GatewayMessage
	if err := json.NewDecoder(finalReadCloser).Decode(&message); err != nil {
		g.Logger().Error(g.formatLogs("error decoding websocket message: ", err))
		return nil, err
	}
	return &message, nil
}
