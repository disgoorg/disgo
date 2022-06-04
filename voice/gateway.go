package voice

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"sync"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/json"
	"github.com/disgoorg/log"
	"github.com/gorilla/websocket"
)

var (
	ErrGatewayNotConnected     = fmt.Errorf("voice gateway not connected")
	ErrGatewayAlreadyConnected = fmt.Errorf("voice gateway already connected")
)

var GatewayVersion = 4

type GatewayStatus int

const (
	GatewayStatusUnconnected GatewayStatus = iota
	GatewayStatusConnecting
	GatewayStatusWaitingForHello
	GatewayStatusIdentifying
	GatewayStatusResuming
	GatewayStatusWaitingForReady
	GatewayStatusReady
	GatewayStatusDisconnected
)

type (
	EventHandlerFunc  func(opCode GatewayOpcode, data GatewayMessageData)
	CloseHandlerFunc  func(gateway *Gateway, err error)
	GatewayCreateFunc func(state State, eventHandlerFunc EventHandlerFunc, closeHandlerFunc CloseHandlerFunc, opts ...GatewayConfigOpt) *Gateway
)

func NewGateway(state State, eventHandlerFunc EventHandlerFunc, closeHandlerFunc CloseHandlerFunc, opts ...GatewayConfigOpt) *Gateway {
	config := DefaultGatewayConfig()
	config.Apply(opts)

	return &Gateway{
		config:           *config,
		eventHandlerFunc: eventHandlerFunc,
		closeHandlerFunc: closeHandlerFunc,
		state:            state,
	}
}

type Gateway struct {
	config           GatewayConfig
	eventHandlerFunc EventHandlerFunc
	closeHandlerFunc CloseHandlerFunc
	state            State

	ssrc uint32

	conn   *websocket.Conn
	status GatewayStatus
	mu     sync.Mutex

	heartbeatTicker       *time.Ticker
	lastHeartbeatSent     time.Time
	lastHeartbeatReceived time.Time
	lastNonce             int64
}

func (g *Gateway) Logger() log.Logger {
	return g.config.Logger
}

func (g *Gateway) SSRC() uint32 {
	return g.ssrc
}

func (g *Gateway) Open(ctx context.Context) error {
	g.Logger().Debug("opening voice gateway connection")

	g.mu.Lock()
	defer g.mu.Unlock()
	if g.conn != nil {
		return ErrGatewayAlreadyConnected
	}
	g.status = GatewayStatusConnecting

	gatewayURL := fmt.Sprintf("wss://%s?v=%d", g.state.endpoint, GatewayVersion)
	g.Logger().Debugf("connecting to gateway at: %s", gatewayURL)
	g.lastHeartbeatSent = time.Now().UTC()
	conn, rs, err := g.config.Dialer.DialContext(ctx, gatewayURL, nil)
	if err != nil {
		g.Close()
		body := "null"
		if rs != nil && rs.Body != nil {
			defer rs.Body.Close()
			rawBody, bErr := ioutil.ReadAll(rs.Body)
			if bErr != nil {
				g.Logger().Error("error while reading voice response body: ", err)
			}
			body = string(rawBody)
		}

		g.Logger().Error("error connecting to the gateway. url: %s, error: %s, body: %s", gatewayURL, err, body)
		return err
	}

	conn.SetCloseHandler(func(code int, text string) error {
		return nil
	})

	g.conn = conn

	g.status = GatewayStatusWaitingForHello

	go g.listen(g.conn)
	return nil
}

func (g *Gateway) Close() {
	g.CloseWithCode(websocket.CloseNormalClosure, "Shutting down")
}

func (g *Gateway) CloseWithCode(code int, message string) {
	if g.heartbeatTicker != nil {
		g.Logger().Debug("closing heartbeat goroutines...")
		g.heartbeatTicker.Stop()
		g.heartbeatTicker = nil
	}

	g.mu.Lock()
	defer g.mu.Unlock()
	if g.conn != nil {
		g.Logger().Debugf("closing gateway connection with code: %d, message: %s", code, message)
		if err := g.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(code, message)); err != nil && err != websocket.ErrCloseSent {
			g.Logger().Debug("error writing close code. error: ", err)
		}
		_ = g.conn.Close()
		g.conn = nil

		// clear resume data as we closed gracefully
		if code == websocket.CloseNormalClosure || code == websocket.CloseGoingAway {
			g.ssrc = 0
		}
	}
}

func (g *Gateway) heartbeat(heartbeatInterval time.Duration) {
	g.heartbeatTicker = time.NewTicker(heartbeatInterval)
	defer g.heartbeatTicker.Stop()
	defer g.Logger().Debug("exiting heartbeat goroutine...")

	for range g.heartbeatTicker.C {
		g.sendHeartbeat()
	}
}

func (g *Gateway) sendHeartbeat() {
	g.Logger().Debug("sending heartbeat...")

	g.lastNonce = time.Now().UnixMilli()
	if err := g.Send(GatewayOpcodeHeartbeat, GatewayMessageDataHeartbeat(g.lastNonce)); err != nil && err != ErrGatewayNotConnected {
		g.Logger().Error("failed to send heartbeat. error: ", err)
		g.CloseWithCode(websocket.CloseServiceRestart, "heartbeat timeout")
		go g.reconnect(context.TODO())
		return
	}
	g.lastHeartbeatSent = time.Now().UTC()
}

func (g *Gateway) listen(conn *websocket.Conn) {
	defer g.Logger().Debug("exiting listen goroutine...")
loop:
	for {
		_, reader, err := conn.NextReader()
		if err != nil {
			g.mu.Lock()
			sameConnection := g.conn == conn
			g.mu.Unlock()

			// if sameConnection is false, it means the connection has been closed by the user, and we can just exit
			if !sameConnection {
				return
			}

			reconnect := true
			if closeError, ok := err.(*websocket.CloseError); ok {
				closeCode := GatewayCloseEventCode(closeError.Code)
				reconnect = closeCode.ShouldReconnect()
			} else if errors.Is(err, net.ErrClosed) {
				// we closed the connection ourselves. Don't try to reconnect here
				reconnect = false
			} else {
				g.Logger().Debug("failed to read next message from gateway. error: ", err)
			}
			if g.config.AutoReconnect && reconnect {
				go g.reconnect(context.TODO())
			} else {
				g.Close()
				if g.closeHandlerFunc != nil {
					go g.closeHandlerFunc(g, err)
				}
			}
			return
		}

		message, err := g.parseGatewayMessage(reader)
		if err != nil {
			g.Logger().Error("error while parsing gateway event. error: ", err)
			continue
		}

		switch d := message.D.(type) {
		case GatewayMessageDataHello:
			g.status = GatewayStatusWaitingForReady
			g.lastHeartbeatReceived = time.Now().UTC()
			go g.heartbeat(time.Duration(d.HeartbeatInterval) * time.Millisecond)

			if g.ssrc == 0 {
				g.status = GatewayStatusIdentifying
				err = g.Send(GatewayOpcodeIdentify, GatewayMessageDataIdentify{
					GuildID:   g.state.guildID,
					UserID:    g.state.userID,
					SessionID: g.state.sessionID,
					Token:     g.state.token,
				})
			} else {
				g.status = GatewayStatusResuming
				err = g.Send(GatewayOpcodeIdentify, GatewayMessageDataResume{
					GuildID:   g.state.guildID,
					SessionID: g.state.sessionID,
					Token:     g.state.token,
				})
			}
			if err != nil {
				g.CloseWithCode(websocket.CloseServiceRestart, "failed to send identify or resume")
				go g.reconnect(context.TODO())
				return
			}

		case GatewayMessageDataReady:
			g.status = GatewayStatusReady
			g.ssrc = d.SSRC

		case GatewayMessageDataHeartbeatACK:
			g.config.Logger.Debug("received heartbeat ack")
			if int64(d) != g.lastNonce {
				g.Logger().Errorf("received heartbeat ack with nonce: %d, expected nonce: %d", int64(d), g.lastNonce)
				go g.reconnect(context.TODO())
				break loop
			}
			g.lastHeartbeatReceived = time.Now().UTC()
		}
		g.eventHandlerFunc(message.Op, message.D)
	}
}

func (g *Gateway) Latency() time.Duration {
	return g.lastHeartbeatReceived.Sub(g.lastHeartbeatSent)
}

func (g *Gateway) Send(op GatewayOpcode, d GatewayMessageData) error {
	data, err := json.Marshal(GatewayMessage{
		Op: op,
		D:  d,
	})
	if err != nil {
		return err
	}
	return g.send(websocket.TextMessage, data)
}

func (g *Gateway) send(messageType int, data []byte) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.conn == nil {
		return ErrGatewayNotConnected
	}

	g.Logger().Infof("sending message to gateway. data: %s", string(data))
	return g.conn.WriteMessage(messageType, data)
}

func (g *Gateway) reconnectTry(ctx context.Context, try int, delay time.Duration) error {
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

	g.Logger().Debug("reconnecting gateway...")
	if err := g.Open(ctx); err != nil {
		if err == discord.ErrGatewayAlreadyConnected {
			return err
		}
		g.Logger().Error("failed to reconnect gateway. error: ", err)
		g.status = GatewayStatusDisconnected
		return g.reconnectTry(ctx, try+1, delay)
	}
	return nil
}

func (g *Gateway) reconnect(ctx context.Context) {
	err := g.reconnectTry(ctx, 0, time.Second)
	if err != nil {
		g.Logger().Error("failed to reopen gateway", err)
	}
}

func (g *Gateway) parseGatewayMessage(reader io.Reader) (GatewayMessage, error) {
	buff := &bytes.Buffer{}
	data, _ := io.ReadAll(io.TeeReader(reader, buff))
	g.Logger().Infof("received message from gateway. data: %s", string(data))

	var message GatewayMessage
	if err := json.NewDecoder(buff).Decode(&message); err != nil {
		g.Logger().Error("error decoding voice websocket message: ", err)
		return GatewayMessage{}, err
	}
	return message, nil
}
