package gateway

import (
	"bytes"
	"compress/zlib"
	"context"
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

func (g *gatewayImpl) Open() error {
	return g.OpenCtx(context.Background())
}

func (g *gatewayImpl) OpenCtx(ctx context.Context) error {
	if g.lastSequenceReceived == nil || g.sessionID == nil {
		g.status = StatusConnecting
	} else {
		g.status = StatusReconnecting
	}

	g.Logger().Infof("[%d/%d] starting ws...", g.shardID, g.shardCount)

	gatewayURL := g.url + "?v=" + route.APIVersion + "&encoding=json"
	var rs *http.Response
	var err error
	g.conn, rs, err = websocket.DefaultDialer.DialContext(ctx, gatewayURL, nil)
	if err != nil {
		g.Close()
		var body []byte
		if rs != nil && rs.Body != nil {
			body, err = ioutil.ReadAll(rs.Body)
			if err != nil {
				g.Logger().Errorf("[%d/%d] error while reading response body: ", g.shardID, g.shardCount, err)
				return err
			}
		} else {
			body = []byte("null")
		}

		g.Logger().Errorf("[%d/%d] error connecting to gateway. url: %s, error: %s, body: %s", g.shardID, g.shardCount, gatewayURL, err, string(body))
		return err
	}

	g.conn.SetCloseHandler(func(code int, error string) error {
		g.Logger().Infof("[%d/%d] connection to websocket closed with code: %d, error: %s", g.shardID, g.shardCount, code, error)
		return nil
	})

	g.status = StatusWaitingForHello

	go g.listen()

	return nil
}

func (g *gatewayImpl) Close() {
	g.closeWithCode(websocket.CloseNormalClosure)
}

func (g *gatewayImpl) Status() Status {
	return g.status
}

func (g *gatewayImpl) Send(command discord.GatewayCommand) error {
	return g.SendCtx(context.Background(), command)
}

func (g *gatewayImpl) SendCtx(ctx context.Context, command discord.GatewayCommand) error {
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
	g.Logger().Debugf("[%d/%d] sending gateway command: ", g.shardID, g.shardCount, string(data))
	return g.conn.WriteMessage(websocket.TextMessage, data)
}

func (g *gatewayImpl) Latency() time.Duration {
	return g.lastHeartbeatReceived.Sub(g.lastHeartbeatSent)
}

func (g *gatewayImpl) reconnect(delay time.Duration) {
	go func() {
		time.Sleep(delay)

		if g.Status() == StatusConnecting || g.Status() == StatusReconnecting {
			g.Logger().Errorf("[%d/%d] tried to reconnect gateway while connecting/reconnecting", g.shardID, g.shardCount)
			return
		}
		g.Logger().Info("[%d/%d] reconnecting gateway...", g.shardID, g.shardCount)
		if err := g.Open(); err != nil {
			g.Logger().Errorf("[%d/%d] failed to reconnect gateway: ", g.shardID, g.shardCount, err)
			g.status = StatusDisconnected
			g.reconnect(delay * 2)
		}
	}()
}

func (g *gatewayImpl) closeWithCode(code int) {
	if g.heartbeatChan != nil {
		g.Logger().Infof("[%d/%d] closing gateway goroutines...", g.shardID, g.shardCount)
		close(g.heartbeatChan)
		g.heartbeatChan = nil
		g.Logger().Infof("[%d/%d] closed gateway goroutines", g.shardID, g.shardCount)
	}
	if g.conn != nil {
		err := g.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(code, ""))
		if err != nil {
			g.Logger().Errorf("[%d/%d] error writing close code: ", g.shardID, g.shardCount, err)
		}

		// TODO: Wait for Discord to actually close the connection.
		time.Sleep(1 * time.Second)

		err = g.conn.Close()
		if err != nil {
			g.Logger().Errorf("[%d/%d] error closing conn: ", g.shardID, g.shardCount, err)
		}
		g.conn = nil
	}
}

func (g *gatewayImpl) heartbeat() {
	defer g.Logger().Debugf("[%d/%d] exiting heartbeat goroutine...", g.shardID, g.shardCount)
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
	g.Logger().Debugf("[%d/%d] sending heartbeat...", g.shardID, g.shardCount)
	// TODO: check this
	/*
		heartbeatEvent := &events.HeartbeatEvent{
			GenericEvent: events.NewGenericEvent(g.Bot(), 0),
			OldPing:      g.Latency(),
		}*/

	if err := g.Send(discord.NewGatewayCommand(discord.GatewayOpcodeHeartbeat, g.lastSequenceReceived)); err != nil {
		g.Logger().Errorf("[%d/%d] failed to send heartbeat with error: ", g.shardID, g.shardCount, err)
		g.closeWithCode(websocket.CloseServiceRestart)
		g.reconnect(1 * time.Second)
	}
	g.lastHeartbeatSent = time.Now().UTC()

	//heartbeatEvent.NewPing = g.Latency()
	//g.Bot().EventManager().Dispatch(heartbeatEvent)*/
}

func (g *gatewayImpl) listen() {
	defer g.Logger().Debugf("[%d/%d] exiting listen goroutine...", g.shardID, g.shardCount)
	for {
		if g.conn == nil {
			return
		}
		mt, reader, err := g.conn.NextReader()
		if err != nil {
			g.Logger().Errorf("[%d/%d] error while reading from ws. error: ", g.shardID, g.shardCount, err)
			g.closeWithCode(websocket.CloseServiceRestart)
			g.reconnect(1 * time.Second)
			return
		}

		event, err := g.parseGatewayEvent(mt, reader)
		if err != nil {
			g.Logger().Errorf("[%d/%d] error while unpacking gateway event. error: ", g.shardID, g.shardCount, err)
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
				g.Logger().Errorf("[%d/%d] error parsing op hello payload data: ", g.shardID, g.shardCount, err)
			}

			g.heartbeatInterval = eventData.HeartbeatInterval * time.Millisecond

			if g.lastSequenceReceived == nil || g.sessionID == nil {
				g.status = StatusIdentifying
				g.Logger().Infof("[%d/%d] sending StatusIdentifying command...", g.shardID, g.shardCount)

				identify := discord.IdentifyCommand{
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

				if err = g.Send(discord.NewGatewayCommand(discord.GatewayOpcodeIdentify, identify)); err != nil {
					g.Logger().Errorf("[%d/%d] error sending identify payload: ", g.shardID, g.shardCount, err)
				}
				g.status = StatusWaitingForReady
			} else {
				g.status = StatusResuming
				cmd := discord.NewGatewayCommand(discord.GatewayOpcodeResume, discord.ResumeCommand{
					Token:     g.token,
					SessionID: *g.sessionID,
					Seq:       *g.lastSequenceReceived,
				})
				g.Logger().Infof("[%d/%d] sending StatusResuming command...")

				if err = g.Send(cmd); err != nil {
					g.Logger().Errorf("[%d/%d] error sending resume payload: ", g.shardID, g.shardCount, err)
				}
			}
			g.heartbeatChan = make(chan struct{})
			go g.heartbeat()

		case discord.GatewayOpcodeDispatch:
			g.Logger().Debugf("[%d/%d] received: OpcodeDispatch", g.shardID, g.shardCount)
			if event.S != 0 {
				g.lastSequenceReceived = &event.S
			}
			if event.T == "" {
				g.Logger().Errorf("[%d/%d] received event without T. payload: ", g.shardID, g.shardCount, event)
				continue
			}

			g.Logger().Debugf("[%d/%d] received: '%s', data: %s", g.shardID, g.shardCount, event.T, string(event.D))

			if event.T == discord.GatewayEventTypeReady {
				var readyEvent discord.GatewayEventReady
				if err = json.Unmarshal(event.D, &readyEvent); err != nil {
					g.Logger().Errorf("[%d/%d] Error parsing ready event: ", g.shardID, g.shardCount, err)
					continue
				}
				g.sessionID = &readyEvent.SessionID
				g.status = StatusWaitingForGuilds
				g.Logger().Infof("[%d/%d] ready event received", g.shardID, g.shardCount)
			}

			go g.config.EventHandlerFunc(event.T, event.S, bytes.NewBuffer(event.D))

		case discord.GatewayOpcodeHeartbeat:
			g.Logger().Debugf("[%d/%d] received: OpcodeHeartbeat", g.shardID, g.shardCount)
			g.sendHeartbeat()

		case discord.GatewayOpcodeReconnect:
			g.Logger().Debugf("[%d/%d] received: OpcodeReconnect", g.shardID, g.shardCount)
			g.closeWithCode(websocket.CloseServiceRestart)
			g.reconnect(1 * time.Second)

		case discord.GatewayOpcodeInvalidSession:
			var canResume bool
			if err = json.Unmarshal(event.D, &canResume); err != nil {
				g.Logger().Errorf("[%d/%d] error parsing invalid session data: ", g.shardID, g.shardCount, err)
			}
			g.Logger().Debugf("[%d/%d] received: OpcodeInvalidSession, canResume: ", g.shardID, g.shardCount, canResume)
			if canResume {
				g.closeWithCode(websocket.CloseServiceRestart)
			} else {
				g.Close()
				// clear reconnect info
				g.sessionID = nil
				g.lastSequenceReceived = nil
			}
			g.reconnect(5 * time.Second)

		case discord.GatewayOpcodeHeartbeatACK:
			g.Logger().Debugf("[%d/%d] received: OpcodeHeartbeatACK", g.shardID, g.shardCount)
			g.lastHeartbeatReceived = time.Now().UTC()
		}

	}
}

func (g *gatewayImpl) parseGatewayEvent(mt int, reader io.Reader) (*discord.GatewayPayload, error) {
	if mt == websocket.BinaryMessage {
		g.Logger().Debugf("[%d/%d] binary message received. decompressing...", g.shardID, g.shardCount)
		readCloser, err := zlib.NewReader(reader)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decompress zlib")
		}
		defer readCloser.Close()
		reader = readCloser
	}

	decoder := json.NewDecoder(reader)
	var event discord.GatewayPayload
	if err := decoder.Decode(&event); err != nil {
		g.Logger().Errorf("[%d/%d] error decoding websocket message: ", g.shardID, g.shardCount, err)
		return nil, err
	}
	return &event, nil
}
