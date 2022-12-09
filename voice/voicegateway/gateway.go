package voicegateway

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"syscall"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/json"
	"github.com/disgoorg/snowflake/v2"
	"github.com/gorilla/websocket"
)

var (
	ErrGatewayNotConnected     = fmt.Errorf("voice gateway not connected")
	ErrGatewayAlreadyConnected = fmt.Errorf("voice gateway already connected")
)

var Version = 4

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
	EventHandlerFunc func(opCode Opcode, data MessageData)
	CloseHandlerFunc func(gateway Gateway, err error)
	CreateFunc       func(eventHandlerFunc EventHandlerFunc, closeHandlerFunc CloseHandlerFunc, opts ...ConfigOpt) Gateway
)

type State struct {
	GuildID snowflake.ID
	UserID  snowflake.ID

	ChannelID snowflake.ID
	SessionID string
	Token     string
	Endpoint  string
}

type Gateway interface {
	SSRC() uint32
	Latency() time.Duration

	Open(ctx context.Context, state State) error
	Close()
	CloseWithCode(code int, message string)

	Send(ctx context.Context, opCode Opcode, data MessageData) error
}

func New(eventHandlerFunc EventHandlerFunc, closeHandlerFunc CloseHandlerFunc, opts ...ConfigOpt) Gateway {
	config := DefaultConfig()
	config.Apply(opts)

	return &gatewayImpl{
		config:           config,
		eventHandlerFunc: eventHandlerFunc,
		closeHandlerFunc: closeHandlerFunc,
	}
}

type gatewayImpl struct {
	config           Config
	eventHandlerFunc EventHandlerFunc
	closeHandlerFunc CloseHandlerFunc

	ssrc  uint32
	state State

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
	g.state = state

	g.connMu.Lock()
	defer g.connMu.Unlock()
	if g.conn != nil {
		return ErrGatewayAlreadyConnected
	}
	g.status = StatusConnecting

	gatewayURL := fmt.Sprintf("wss://%s?v=%d", state.Endpoint, Version)
	g.config.Logger.Debugf("connecting to voice gateway at: %s", gatewayURL)
	g.lastHeartbeatSent = time.Now().UTC()
	conn, rs, err := g.config.Dialer.DialContext(ctx, gatewayURL, nil)
	if err != nil {
		g.Close()
		defer rs.Body.Close()
		return fmt.Errorf("error connecting to voice gateway. err: %w", err)
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
		g.config.Logger.Debug("closing heartbeat goroutines...")
		g.heartbeatTicker.Stop()
		g.heartbeatTicker = nil
	}

	g.connMu.Lock()
	defer g.connMu.Unlock()
	if g.conn != nil {
		g.config.Logger.Debugf("closing voice gateway connection with code: %d, message: %s", code, message)
		if err := g.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(code, message)); err != nil && err != websocket.ErrCloseSent {
			g.config.Logger.Debug("error writing close code. error: ", err)
		}
		_ = g.conn.Close()
		g.conn = nil

		// clear resume data as we closed gracefully
		if code == websocket.CloseNormalClosure || code == websocket.CloseGoingAway {
			g.ssrc = 0
		}
	}
}

func (g *gatewayImpl) heartbeat() {
	g.heartbeatTicker = time.NewTicker(g.heartbeatInterval)
	defer g.heartbeatTicker.Stop()
	defer g.config.Logger.Debug("exiting voice heartbeat goroutine...")

	for range g.heartbeatTicker.C {
		g.sendHeartbeat()
	}
}

func (g *gatewayImpl) sendHeartbeat() {
	g.lastNonce = time.Now().UnixMilli()
	ctx, cancel := context.WithTimeout(context.Background(), g.heartbeatInterval)
	defer cancel()

	if err := g.Send(ctx, OpcodeHeartbeat, MessageDataHeartbeat(g.lastNonce)); err != nil {
		if err != ErrGatewayNotConnected || errors.Is(err, syscall.EPIPE) {
			return
		}
		g.config.Logger.Error("failed to send heartbeat. error: ", err)
		g.CloseWithCode(websocket.CloseServiceRestart, "heartbeat timeout")
		go g.reconnect()
		return
	}
	g.lastHeartbeatSent = time.Now().UTC()
}

func (g *gatewayImpl) listen(conn *websocket.Conn) {
	defer g.config.Logger.Debug("exiting listen goroutine...")
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
			if closeError, ok := err.(*websocket.CloseError); ok {
				closeCode := CloseEventCodeByCode(closeError.Code)
				reconnect = closeCode.Reconnect
			} else if errors.Is(err, net.ErrClosed) {
				// we closed the connection ourselves. Don't try to reconnect here
				reconnect = false
			} else {
				g.config.Logger.Debug("failed to read next message from voice gateway. error: ", err)
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
			g.config.Logger.Error("error while parsing voice gateway event. error: ", err)
			continue
		}

		switch d := message.D.(type) {
		case MessageDataHello:
			g.status = StatusWaitingForReady
			g.lastHeartbeatReceived = time.Now().UTC()
			g.heartbeatInterval = time.Duration(d.HeartbeatInterval) * time.Millisecond
			go g.heartbeat()

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

			if g.ssrc == 0 {
				g.status = StatusIdentifying
				err = g.Send(ctx, OpcodeIdentify, MessageDataIdentify{
					GuildID:   g.state.GuildID,
					UserID:    g.state.UserID,
					SessionID: g.state.SessionID,
					Token:     g.state.Token,
				})
			} else {
				g.status = StatusResuming
				err = g.Send(ctx, OpcodeIdentify, MessageDataResume{
					GuildID:   g.state.GuildID,
					SessionID: g.state.SessionID,
					Token:     g.state.Token,
				})
			}
			cancel()
			if err != nil {
				g.CloseWithCode(websocket.CloseServiceRestart, "failed to send identify or resume")
				go g.reconnect()
				return
			}

		case MessageDataReady:
			g.status = StatusReady
			g.ssrc = d.SSRC

		case MessageDataHeartbeatACK:
			if int64(d) != g.lastNonce {
				g.config.Logger.Errorf("received heartbeat ack with nonce: %d, expected nonce: %d", int64(d), g.lastNonce)
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

func (g *gatewayImpl) Send(ctx context.Context, op Opcode, d MessageData) error {
	data, err := json.Marshal(Message{
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

	g.config.Logger.Trace("sending message to voice gateway. data: ", string(data))
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

	g.config.Logger.Debug("reconnecting voice gateway...")
	if err := g.Open(ctx, g.state); err != nil {
		if err == discord.ErrGatewayAlreadyConnected {
			return err
		}
		g.config.Logger.Error("failed to reconnect voice gateway. error: ", err)
		g.status = StatusDisconnected
		return g.reconnectTry(ctx, try+1)
	}
	return nil
}

func (g *gatewayImpl) reconnect() {
	if err := g.reconnectTry(context.Background(), 0); err != nil {
		g.config.Logger.Error("failed to reopen voice gateway", err)
	}
}

func (g *gatewayImpl) parseMessage(r io.Reader) (Message, error) {
	buff := &bytes.Buffer{}
	data, _ := io.ReadAll(io.TeeReader(r, buff))
	g.config.Logger.Tracef("received message from voice gateway. data: %s", string(data))

	var message Message
	if err := json.NewDecoder(buff).Decode(&message); err != nil {
		return Message{}, err
	}
	return message, nil
}
