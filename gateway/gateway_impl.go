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
	"github.com/DisgoOrg/disgo/internal/tokenhelper"
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

func (g *gatewayImpl) formatLogsf(format string, a ...interface{}) string {
	return fmt.Sprintf("[%d/%d] %s", g.shardID, g.shardCount, fmt.Sprintf(format, a...))
}

func (g *gatewayImpl) formatLogs(a ...interface{}) string {
	return fmt.Sprintf("[%d/%d] %s", g.shardID, g.shardCount, fmt.Sprint(a...))
}

func (g *gatewayImpl) Open() error {
	return g.OpenCtx(context.Background())
}

func (g *gatewayImpl) OpenCtx(ctx context.Context) error {
	g.Logger().Info(g.formatLogs("opening gateway connection"))
	if g.lastSequenceReceived == nil || g.sessionID == nil {
		g.status = StatusConnecting
	} else {
		g.status = StatusReconnecting
	}

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
				g.Logger().Error(g.formatLogs("error while reading response body: ", err))
				return err
			}
		} else {
			body = []byte("null")
		}

		g.Logger().Error(g.formatLogsf("error connecting to gateway. url: %s, error: %s, body: %s", gatewayURL, err, string(body)))
		return err
	}

	//g.conn.SetCloseHandler(nil)

	g.status = StatusWaitingForHello

	go g.listen()

	return nil
}

func (g *gatewayImpl) Close() {
	g.CloseWithCode(websocket.CloseNormalClosure)
}

func (g *gatewayImpl) CloseWithCode(code int) {
	g.Logger().Info(g.formatLogs("closing gateway connection with code: ", code))
	if g.heartbeatChan != nil {
		g.Logger().Debug(g.formatLogs("closing heartbeat goroutines..."))
		close(g.heartbeatChan)
		g.heartbeatChan = nil
	}
	if g.conn != nil {
		err := g.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(code, ""))
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
	g.Logger().Debugf(g.formatLogs("sending gateway command: ", string(data)))
	return g.conn.WriteMessage(websocket.TextMessage, data)
}

func (g *gatewayImpl) Latency() time.Duration {
	return g.lastHeartbeatReceived.Sub(g.lastHeartbeatSent)
}

func (g *gatewayImpl) reconnect(try int, delay time.Duration) {
	if try >= g.config.MaxReconnectTries {
		g.Logger().Error(g.formatLogs("failed to reconnect. exceeded max reconnect tries of: ", g.config.MaxReconnectTries))
		return
	}
	time.Sleep(delay)

	if g.Status() == StatusConnecting || g.Status() == StatusReconnecting {
		g.Logger().Error(g.formatLogs("tried to reconnect gateway while connecting/reconnecting"))
		return
	}
	g.Logger().Info(g.formatLogs("reconnecting gateway..."))
	if err := g.Open(); err != nil {
		g.Logger().Error(g.formatLogs("failed to reconnect gateway. error: ", err))
		g.status = StatusDisconnected
		g.reconnect(try+1, delay*2)
	}
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
	g.Logger().Debug(g.formatLogs("sending heartbeat..."))

	if err := g.Send(discord.NewGatewayCommand(discord.GatewayOpcodeHeartbeat, g.lastSequenceReceived)); err != nil && err != discord.ErrShardNotConnected {
		g.Logger().Error(g.formatLogs("failed to send heartbeat. error: ", err))
		g.CloseWithCode(websocket.CloseServiceRestart)
		go g.reconnect(0, 1*time.Second)
		return
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
			reconnect := true
			if closeError, ok := err.(*websocket.CloseError); ok {
				switch discord.GatewayCloseEventCode(closeError.Code) {
				case websocket.CloseNormalClosure:
					g.Logger().Debug(g.formatLogs("gracefully closed gateway"))
					reconnect = false

				case discord.GatewayCloseEventCodeUnknownError:
					g.Logger().Error(g.formatLogsf("unknown gateway error tying ro reconnect. code: %d error: %s", closeError.Code, closeError.Text))

				case discord.GatewayCloseEventCodeRateLimited:
					g.Logger().Error(g.formatLogs("sent too much gateway commands. reconnecting..."))

				case discord.GatewayCloseEventCodeInvalidShard:
					g.Logger().Error(g.formatLogs("invalid sharding config supplied"))
					reconnect = false

				case discord.GatewayCloseEventCodeInvalidIntents:
					g.Logger().Error(g.formatLogs("invalid gateway intents supplied. intents: ", g.config.GatewayIntents))
					reconnect = false

				case discord.GatewayCloseEventCodeDisallowedIntents:
					var intentsURL string
					if id, err := tokenhelper.IDFromToken(g.token); err == nil {
						intentsURL = fmt.Sprintf("https://discord.com/developers/applications/%s/bot", *id)
					} else {
						intentsURL = "https://discord.com/developers/applications"
					}
					g.Logger().Error(g.formatLogsf("disallowed gateway intents supplied. go to %s and enable the privileged intent for your application. intents: %d", intentsURL, g.config.GatewayIntents))
					reconnect = false

				case discord.GatewayCloseEventCodeInvalidSeq:
					g.Logger().Error(g.formatLogs("invalid sequence provided. reconnecting..."))
					g.lastSequenceReceived = nil

				default:
					g.Logger().Error(g.formatLogsf("unknown close code trying to reconnect. code: %d error: %s", closeError.Code, closeError.Text))
				}
			}

			if reconnect {
				go g.reconnect(0, 1*time.Second)
			} else {
				g.Close()
			}
			return
		}

		event, err := g.parseGatewayEvent(mt, reader)
		if err != nil {
			g.Logger().Error(g.formatLogs("error while parsing gateway event. error: ", err))
			continue
		}

		switch event.Op {
		case discord.GatewayOpcodeHello:
			g.lastHeartbeatReceived = time.Now().UTC()

			var eventData discord.GatewayEventHello
			if err = json.Unmarshal(event.D, &eventData); err != nil {
				g.Logger().Error(g.formatLogs("error parsing op hello payload data: ", err))
				continue
			}

			g.heartbeatInterval = eventData.HeartbeatInterval * time.Millisecond

			if g.lastSequenceReceived == nil || g.sessionID == nil {
				g.status = StatusIdentifying
				g.Logger().Info(g.formatLogs("sending Identify command..."))

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
					g.Logger().Error(g.formatLogs("error sending Identify command err: ", err))
				}
				g.status = StatusWaitingForReady
			} else {
				g.status = StatusResuming
				resume := discord.ResumeCommand{
					Token:     g.token,
					SessionID: *g.sessionID,
					Seq:       *g.lastSequenceReceived,
				}

				g.Logger().Info(g.formatLogs("sending Resume command..."))
				if err = g.Send(discord.NewGatewayCommand(discord.GatewayOpcodeResume, resume)); err != nil {
					g.Logger().Error(g.formatLogs("error sending resume command err: ", err))
				}
			}
			g.heartbeatChan = make(chan struct{})
			go g.heartbeat()

		case discord.GatewayOpcodeDispatch:
			g.Logger().Debugf(g.formatLogsf("received: OpcodeDispatch %s, data: %s", event.T, string(event.D)))

			// set last sequence received
			g.lastSequenceReceived = &event.S

			// get session id here
			if event.T == discord.GatewayEventTypeReady {
				var readyEvent discord.GatewayEventReady
				if err = json.Unmarshal(event.D, &readyEvent); err != nil {
					g.Logger().Error(g.formatLogs("Error parsing ready event. error: ", err))
					continue
				}
				g.sessionID = &readyEvent.SessionID
				g.status = StatusReady
				g.Logger().Debug(g.formatLogs("ready event received"))
			}

			// push event to the command manager
			g.config.EventHandlerFunc(event.T, event.S, bytes.NewBuffer(event.D))

		case discord.GatewayOpcodeHeartbeat:
			g.Logger().Debug(g.formatLogs("received: OpcodeHeartbeat"))
			g.sendHeartbeat()

		case discord.GatewayOpcodeReconnect:
			g.Logger().Debug(g.formatLogs("received: OpcodeReconnect"))
			g.CloseWithCode(websocket.CloseServiceRestart)
			go g.reconnect(0, 1 * time.Second)

		case discord.GatewayOpcodeInvalidSession:
			var canResume bool
			if err = json.Unmarshal(event.D, &canResume); err != nil {
				g.Logger().Error(g.formatLogs("error parsing invalid session data. error: ", err))
			}
			g.Logger().Debug(g.formatLogs("received: OpcodeInvalidSession, canResume: ", canResume))
			if canResume {
				g.CloseWithCode(websocket.CloseServiceRestart)
			} else {
				g.Close()
				// clear reconnect info
				g.sessionID = nil
				g.lastSequenceReceived = nil
			}
			go g.reconnect(0, 5*time.Second)

		case discord.GatewayOpcodeHeartbeatACK:
			g.Logger().Debug(g.formatLogs("received: OpcodeHeartbeatACK"))
			g.lastHeartbeatReceived = time.Now().UTC()
		}
	}
}

func (g *gatewayImpl) parseGatewayEvent(mt int, reader io.Reader) (*discord.GatewayPayload, error) {
	if mt == websocket.BinaryMessage {
		g.Logger().Debug(g.formatLogs("binary message received. decompressing..."))
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
		g.Logger().Error(g.formatLogs("error decoding websocket message: ", err))
		return nil, err
	}
	return &event, nil
}
