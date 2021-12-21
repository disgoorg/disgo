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

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway/grate"
	"github.com/DisgoOrg/disgo/json"
	"github.com/DisgoOrg/disgo/rest/route"
	"github.com/DisgoOrg/log"
	"github.com/pkg/errors"

	"github.com/gorilla/websocket"
)

var _ Gateway = (*gatewayImpl)(nil)

func New(token string, url string, shardID int, shardCount int, eventHandlerFunc EventHandlerFunc, config *Config) Gateway {
	if config == nil {
		config = &DefaultConfig
	}
	if config.Logger == nil {
		config.Logger = log.Default()
	}
	if config.RateLimiterConfig == nil {
		config.RateLimiterConfig = &grate.DefaultConfig
	}
	if config.RateLimiterConfig.Logger == nil {
		config.RateLimiterConfig.Logger = config.Logger
	}
	if config.RateLimiter == nil {
		config.RateLimiter = grate.NewLimiter(config.RateLimiterConfig)
	}
	config.EventHandlerFunc = eventHandlerFunc

	return &gatewayImpl{
		config:     *config,
		token:      token,
		url:        url,
		shardID:    shardID,
		shardCount: shardCount,
		status:     StatusUnconnected,
	}
}

type gatewayImpl struct {
	config                Config
	token                 string
	url                   string
	shardID               int
	shardCount            int
	conn                  *websocket.Conn
	heartbeatChan         chan struct{}
	status                Status
	heartbeatInterval     time.Duration
	lastHeartbeatSent     time.Time
	lastHeartbeatReceived time.Time
	sessionID             *string
	lastSequenceReceived  *int
}

func (g *gatewayImpl) Logger() log.Logger {
	return g.config.Logger
}

func (g *gatewayImpl) Config() Config {
	return g.config
}

func (g *gatewayImpl) ShardID() int {
	return g.shardID
}

func (g *gatewayImpl) ShardCount() int {
	return g.shardCount
}

func (g *gatewayImpl) formatLogsf(format string, a ...interface{}) string {
	if g.shardCount > 1 {
		return fmt.Sprintf("[%d/%d] %s", g.shardID, g.shardCount, fmt.Sprintf(format, a...))
	}
	return fmt.Sprintf(format, a...)
}

func (g *gatewayImpl) formatLogs(a ...interface{}) string {
	if g.shardCount > 1 {
		return fmt.Sprintf("[%d/%d] %s", g.shardID, g.shardCount, fmt.Sprint(a...))
	}
	return fmt.Sprint(a...)
}

func (g *gatewayImpl) Open(ctx context.Context) error {
	if g.lastSequenceReceived == nil || g.sessionID == nil {
		g.status = StatusConnecting
	} else {
		g.status = StatusReconnecting
	}

	g.Logger().Infof(g.formatLogs("starting ws..."))

	gatewayURL := g.url + "?v=" + route.APIVersion + "&encoding=json"
	var rs *http.Response
	var err error
	g.conn, rs, err = websocket.DefaultDialer.DialContext(ctx, gatewayURL, nil)
	if err != nil {
		if err = g.Close(ctx); err != nil {
			return err
		}
		var body []byte
		if rs != nil && rs.Body != nil {
			body, err = ioutil.ReadAll(rs.Body)
			if err != nil {
				g.Logger().Error(g.formatLogs("error while reading response body: ", err))
				return err
			}
		} else {
			body = []byte("null")
		}

		g.Logger().Error(g.formatLogsf("error connecting to gateway. url: %s, error: %s, body: %s", gatewayURL, err, string(body)))
		return err
	}

	g.conn.SetCloseHandler(func(code int, error string) error {
		g.Logger().Info(g.formatLogsf("connection to websocket closed with code: %d, error: %s", code, error))
		return nil
	})

	g.status = StatusWaitingForHello

	go g.listen()

	return nil
}

func (g *gatewayImpl) Close(ctx context.Context) error {
	return g.closeWithCode(ctx, websocket.CloseNormalClosure)
}

func (g *gatewayImpl) Status() Status {
	return g.status
}

func (g *gatewayImpl) Send(ctx context.Context, command discord.GatewayCommand) error {
	if g.conn == nil {
		return discord.ErrShardNotConnected
	}

	if err := g.config.RateLimiter.Wait(ctx); err != nil {
		return err
	}

	defer g.config.RateLimiter.Unlock()
	data, err := json.Marshal(command)
	if err != nil {
		return err
	}
	g.Logger().Debugf(g.formatLogs("sending gateway command: ", string(data)))
	return g.conn.WriteMessage(websocket.TextMessage, data)
}

func (g *gatewayImpl) Latency() time.Duration {
	return g.lastHeartbeatReceived.Sub(g.lastHeartbeatSent)
}

func (g *gatewayImpl) reconnect(ctx context.Context, delay time.Duration) {
	go func() {
		time.Sleep(delay)

		if g.Status() == StatusConnecting || g.Status() == StatusReconnecting {
			g.Logger().Error(g.formatLogs("tried to reconnect gateway while connecting/reconnecting"))
			return
		}
		g.Logger().Info(g.formatLogs("reconnecting gateway..."))
		if err := g.Open(ctx); err != nil {
			g.Logger().Error(g.formatLogs("failed to reconnect gateway: ", err))
			g.status = StatusDisconnected
			g.reconnect(ctx, delay*2)
		}
	}()
}

func (g *gatewayImpl) closeWithCode(ctx context.Context, code int) error {
	if err := g.config.RateLimiter.Wait(ctx); err != nil {
		return err
	}
	if g.heartbeatChan != nil {
		g.Logger().Info(g.formatLogs("closing gateway goroutines..."))
		close(g.heartbeatChan)
		g.heartbeatChan = nil
		g.Logger().Info(g.formatLogs("closed gateway goroutines"))
	}
	if g.conn != nil {
		err := g.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(code, ""))
		if err != nil {
			return errors.Wrap(err, "error writing close code")
		}

		// TODO: Wait for Discord to actually close the connection.
		time.Sleep(1 * time.Second)

		err = g.conn.Close()
		if err != nil {
			return errors.Wrap(err, "error closing websocket")
		}
		g.conn = nil
	}
	return nil
}

func (g *gatewayImpl) heartbeat() {
	defer g.Logger().Debug(g.formatLogs("exiting heartbeat goroutine..."))
	ticker := time.NewTicker(g.heartbeatInterval)
	for {
		select {
		case <-ticker.C:
			g.sendHeartbeat()
		case <-g.heartbeatChan:
			ticker.Stop()
			return
		}
	}
}

func (g *gatewayImpl) sendHeartbeat() {
	g.Logger().Debugf(g.formatLogs("sending heartbeat..."))

	if err := g.Send(context.TODO(), discord.NewGatewayCommand(discord.GatewayOpcodeHeartbeat, discord.HeartbeatCommandData(*g.lastSequenceReceived))); err != nil {
		g.Logger().Error(g.formatLogs("failed to send heartbeat with error: ", err))
		_ = g.closeWithCode(context.TODO(), websocket.CloseServiceRestart)
		g.reconnect(context.TODO(), 1*time.Second)
	}
	g.lastHeartbeatSent = time.Now().UTC()
}

func (g *gatewayImpl) listen() {
	defer g.Logger().Debug(g.formatLogs("exiting listen goroutine..."))
	for {
		if g.conn == nil {
			return
		}
		mt, reader, err := g.conn.NextReader()
		if err != nil {
			g.Logger().Error(g.formatLogs("error while reading from ws. error: ", err))
			_ = g.closeWithCode(context.TODO(), websocket.CloseServiceRestart)
			g.reconnect(context.TODO(), 1*time.Second)
			return
		}

		event, err := g.parseGatewayEvent(mt, reader)
		if err != nil {
			g.Logger().Error(g.formatLogs("error while unpacking gateway event. error: ", err))
			continue
		}

		if event == nil {
			continue
		}

		switch event.Op {
		case discord.GatewayOpcodeHello:
			g.lastHeartbeatReceived = time.Now().UTC()

			var eventData discord.GatewayEventHello
			if err = json.Unmarshal(event.D, &eventData); err != nil {
				g.Logger().Error(g.formatLogs("error parsing op hello payload data. error: ", err))
			}

			g.heartbeatInterval = eventData.HeartbeatInterval * time.Millisecond

			if g.lastSequenceReceived == nil || g.sessionID == nil {
				g.status = StatusIdentifying
				g.Logger().Info(g.formatLogs("sending StatusIdentifying command..."))

				identify := discord.IdentifyCommandData{
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
				if g.shardCount > 1 {
					identify.Shard = []int{g.shardID, g.shardCount}
				}

				if err = g.Send(context.TODO(), discord.NewGatewayCommand(discord.GatewayOpcodeIdentify, identify)); err != nil {
					g.Logger().Error(g.formatLogs("error sending identify payload. error: ", err))
				}
				g.status = StatusWaitingForReady
			} else {
				g.status = StatusResuming
				cmd := discord.NewGatewayCommand(discord.GatewayOpcodeResume, discord.ResumeCommandData{
					Token:     g.token,
					SessionID: *g.sessionID,
					Seq:       *g.lastSequenceReceived,
				})
				g.Logger().Info(g.formatLogs("sending StatusResuming command..."))

				if err = g.Send(context.TODO(), cmd); err != nil {
					g.Logger().Error(g.formatLogs("error sending resume payload. error: ", err))
				}
			}
			g.heartbeatChan = make(chan struct{})
			go g.heartbeat()

		case discord.GatewayOpcodeDispatch:
			g.Logger().Debug(g.formatLogs("received: OpcodeDispatch"))
			if event.S != 0 {
				g.lastSequenceReceived = &event.S
			}
			if event.T == "" {
				g.Logger().Error(g.formatLogs("received event without T. payload: ", event))
				continue
			}

			g.Logger().Debug(g.formatLogsf("received: '%s', data: %s", event.T, string(event.D)))

			if event.T == discord.GatewayEventTypeReady {
				var readyEvent discord.GatewayEventReady
				if err = json.Unmarshal(event.D, &readyEvent); err != nil {
					g.Logger().Error(g.formatLogs("Error parsing ready event. error: ", err))
					continue
				}
				g.sessionID = &readyEvent.SessionID
				g.status = StatusWaitingForGuilds
				g.Logger().Info(g.formatLogs("ready event received"))
			}

			go g.config.EventHandlerFunc(event.T, event.S, bytes.NewBuffer(event.D))

		case discord.GatewayOpcodeHeartbeat:
			g.Logger().Debug(g.formatLogs("received: OpcodeHeartbeat"))
			g.sendHeartbeat()

		case discord.GatewayOpcodeReconnect:
			g.Logger().Debug(g.formatLogs("received: OpcodeReconnect"))
			_ = g.closeWithCode(context.TODO(), websocket.CloseServiceRestart)
			g.reconnect(context.TODO(), 1*time.Second)

		case discord.GatewayOpcodeInvalidSession:
			var canResume bool
			if err = json.Unmarshal(event.D, &canResume); err != nil {
				g.Logger().Error(g.formatLogs("error parsing invalid session data. error: ", err))
			}
			g.Logger().Debug(g.formatLogs("received: OpcodeInvalidSession, canResume: ", canResume))
			if canResume {
				_ = g.closeWithCode(context.TODO(), websocket.CloseServiceRestart)
			} else {
				_ = g.Close(context.TODO())
				// clear reconnect info
				g.sessionID = nil
				g.lastSequenceReceived = nil
			}
			g.reconnect(context.TODO(), 5*time.Second)

		case discord.GatewayOpcodeHeartbeatACK:
			g.Logger().Debug(g.formatLogs("received: OpcodeHeartbeatACK"))
			g.lastHeartbeatReceived = time.Now().UTC()
		}

	}
}

func (g *gatewayImpl) parseGatewayEvent(mt int, reader io.Reader) (*discord.GatewayPayload, error) {
	if mt == websocket.BinaryMessage {
		g.Logger().Debugf("binary message received. decompressing...")
		readCloser, err := zlib.NewReader(reader)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decompress zlib")
		}
		defer func() {
			_ = readCloser.Close()
		}()
		reader = readCloser
	}

	decoder := json.NewDecoder(reader)
	var event discord.GatewayPayload
	if err := decoder.Decode(&event); err != nil {
		g.Logger().Errorf("error decoding websocket message: ", err)
		return nil, err
	}
	return &event, nil
}
