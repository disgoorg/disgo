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
	"github.com/DisgoOrg/disgo/json"
	"github.com/DisgoOrg/disgo/rest/route"
	"github.com/DisgoOrg/log"
	"github.com/pkg/errors"

	"github.com/gorilla/websocket"
)

func New(token string, url string, shardID int, shardCount int, eventHandlerFunc EventHandlerFunc, config *Config) Gateway {
	if config == nil {
		config = &DefaultConfig
	}
	if config.Logger == nil {
		config.Logger = log.Default()
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

// gatewayImpl is what is used to connect to discord
//goland:noinspection GoNameStartsWithPackageName
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

// Open initializes the client and connection to discord
func (g *gatewayImpl) Open() error {
	return g.OpenContext(context.Background())
}

// OpenContext initializes the client and connection to discord and takes a context.Context
func (g *gatewayImpl) OpenContext(ctx context.Context) error {
	if g.lastSequenceReceived == nil || g.sessionID == nil {
		g.status = StatusConnecting
	} else {
		g.status = StatusReconnecting
	}

	g.Logger().Info("starting ws...")

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
				g.Logger().Errorf("error while reading response body: %s", err)
				return err
			}
		} else {
			body = []byte("null")
		}

		g.Logger().Errorf("error connecting to gateway. url: %s, error: %s, body: %s", gatewayURL, err, string(body))
		return err
	}

	g.conn.SetCloseHandler(func(code int, error string) error {
		g.Logger().Infof("connection to websocket closed with code: %d, error: %s", code, error)
		return nil
	})

	g.status = StatusWaitingForHello

	go g.listen()

	return nil
}

// Close cleans up the gateway internals
func (g *gatewayImpl) Close() {
	g.closeWithCode(websocket.CloseNormalClosure)
}

// Status returns the gateway connection status
func (g *gatewayImpl) Status() Status {
	return g.status
}

// Send sends a payload to the Gateway
func (g *gatewayImpl) Send(command discord.GatewayCommand) error {
	return g.conn.WriteJSON(command)
}

// Latency returns the core.Gateway latency
func (g *gatewayImpl) Latency() time.Duration {
	return g.lastHeartbeatReceived.Sub(g.lastHeartbeatSent)
}

func (g *gatewayImpl) reconnect(delay time.Duration) {
	go func() {
		time.Sleep(delay)

		if g.Status() == StatusConnecting || g.Status() == StatusReconnecting {
			g.Logger().Error("tried to reconnect gateway while connecting/reconnecting")
			return
		}
		g.Logger().Info("reconnecting gateway...")
		if err := g.Open(); err != nil {
			g.Logger().Errorf("failed to reconnect gateway: %s", err)
			g.status = StatusDisconnected
			g.reconnect(delay * 2)
		}
	}()
}

func (g *gatewayImpl) closeWithCode(code int) {
	if g.heartbeatChan != nil {
		g.Logger().Info("closing gateway goroutines...")
		close(g.heartbeatChan)
		g.heartbeatChan = nil
		g.Logger().Info("closed gateway goroutines")
	}
	if g.conn != nil {
		err := g.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(code, ""))
		if err != nil {
			g.Logger().Errorf("error writing close code: %s", err)
		}

		// TODO: Wait for Discord to actually close the connection.
		time.Sleep(1 * time.Second)

		err = g.conn.Close()
		if err != nil {
			g.Logger().Errorf("error closing conn: %s", err)
		}
		g.conn = nil
	}
}

func (g *gatewayImpl) heartbeat() {
	defer g.Logger().Debug("exiting heartbeat goroutine...")
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
	g.Logger().Debug("sending heartbeat...")
	// TODO: check this
	/*
		heartbeatEvent := &events.HeartbeatEvent{
			GenericEvent: events.NewGenericEvent(g.Bot(), 0),
			OldPing:      g.Latency(),
		}*/

	if err := g.Send(discord.NewGatewayCommand(discord.OpHeartbeat, g.lastSequenceReceived)); err != nil {
		g.Logger().Errorf("failed to send heartbeat with error: %s", err)
		g.closeWithCode(websocket.CloseServiceRestart)
		g.reconnect(1 * time.Second)
	}
	g.lastHeartbeatSent = time.Now().UTC()

	//heartbeatEvent.NewPing = g.Latency()
	//g.Bot().EventManager().Dispatch(heartbeatEvent)*/
}

func (g *gatewayImpl) listen() {
	defer g.Logger().Debug("exiting listen goroutine...")
	for {
		if g.conn == nil {
			return
		}
		mt, reader, err := g.conn.NextReader()
		if err != nil {
			g.Logger().Errorf("error while reading from ws. error: %s", err)
			g.closeWithCode(websocket.CloseServiceRestart)
			g.reconnect(1 * time.Second)
			return
		}

		event, err := g.parseGatewayEvent(mt, reader)
		if err != nil {
			g.Logger().Error("error while unpacking gateway event. error: ", err)
			continue
		}

		if event == nil {
			continue
		}

		switch event.Op {
		case discord.OpHello:
			g.lastHeartbeatReceived = time.Now().UTC()

			var eventData discord.GatewayEventHello
			if err = json.Unmarshal(event.D, &eventData); err != nil {
				g.Logger().Error("error parsing op hello payload data: ", err)
			}

			g.heartbeatInterval = eventData.HeartbeatInterval * time.Millisecond

			if g.lastSequenceReceived == nil || g.sessionID == nil {
				g.status = StatusIdentifying
				g.Logger().Infof("sending StatusIdentifying command...")

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
				}
				if g.shardCount > 1 {
					identify.Shard = [2]int{g.shardID, g.shardCount}
				}

				if err = g.Send(
					discord.NewGatewayCommand(discord.OpIdentify, identify),
				); err != nil {
					g.Logger().Error("error sending identify payload: ", err)
				}
				g.status = StatusWaitingForReady
			} else {
				g.status = StatusResuming
				cmd := discord.NewGatewayCommand(discord.OpResume, discord.ResumeCommand{
					Token:     g.token,
					SessionID: *g.sessionID,
					Seq:       *g.lastSequenceReceived,
				})
				g.Logger().Infof("sending StatusResuming command...")

				if err = g.Send(cmd); err != nil {
					g.Logger().Error("error sending resume payload: ", err)
				}
			}
			g.heartbeatChan = make(chan struct{})
			go g.heartbeat()

		case discord.OpDispatch:
			g.Logger().Debug("received: OpDispatch")
			if event.S != 0 {
				g.lastSequenceReceived = &event.S
			}
			if event.T == "" {
				g.Logger().Errorf("received event without T. payload: %+v", event)
				continue
			}

			g.Logger().Debugf("received: '%s', data: %s", event.T, string(event.D))

			if event.T == discord.GatewayEventTypeReady {
				var readyEvent discord.GatewayEventReady
				if err = json.Unmarshal(event.D, &readyEvent); err != nil {
					g.Logger().Errorf("Error parsing ready event: %s", err)
					continue
				}
				g.sessionID = &readyEvent.SessionID
				g.status = StatusWaitingForGuilds
				g.Logger().Info("ready event received")
			}

			g.config.EventHandlerFunc(event.T, event.S, bytes.NewBuffer(event.D))

		case discord.OpHeartbeat:
			g.Logger().Debugf("received: OpHeartbeat")
			g.sendHeartbeat()

		case discord.OpReconnect:
			g.Logger().Debugf("received: OpReconnect")
			g.closeWithCode(websocket.CloseServiceRestart)
			g.reconnect(1 * time.Second)

		case discord.OpInvalidSession:
			var canResume bool
			if err = json.Unmarshal(event.D, &canResume); err != nil {
				g.Logger().Errorf("Error parsing invalid session data: %s", err)
			}
			g.Logger().Debugf("received: OpInvalidSession, canResume: %b", canResume)
			if canResume {
				g.closeWithCode(websocket.CloseServiceRestart)
			} else {
				g.Close()
				// clear reconnect info
				g.sessionID = nil
				g.lastSequenceReceived = nil
			}
			g.reconnect(5 * time.Second)

		case discord.OpHeartbeatACK:
			g.Logger().Debugf("received: OpHeartbeatACK")
			g.lastHeartbeatReceived = time.Now().UTC()
		}

	}
}

func (g *gatewayImpl) parseGatewayEvent(mt int, reader io.Reader) (*discord.GatewayPayload, error) {
	if mt == websocket.BinaryMessage {
		g.Logger().Debug("binary message received. decompressing...")
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
		g.Logger().Errorf("error decoding websocket message, %s", err)
		return nil, err
	}
	return &event, nil
}
