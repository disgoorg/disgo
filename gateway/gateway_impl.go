package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/disgo/rest/route"
	"github.com/DisgoOrg/log"

	"github.com/gorilla/websocket"
)

func New(logger log.Logger, restServices rest.Services, token string, config Config, eventHandlerFunc EventHandlerFunc) Gateway {
	return &GatewayImpl{
		logger:           logger,
		restServices:     restServices,
		config:           config,
		token:            token,
		eventHandlerFunc: eventHandlerFunc,
		status:           StatusUnconnected,
	}
}

// GatewayImpl is what is used to connect to discord
//goland:noinspection GoNameStartsWithPackageName
type GatewayImpl struct {
	logger                log.Logger
	restServices          rest.Services
	config                Config
	token                 string
	eventHandlerFunc      EventHandlerFunc
	conn                  *websocket.Conn
	quit                  chan struct{}
	status                Status
	heartbeatInterval     time.Duration
	lastHeartbeatSent     time.Time
	lastHeartbeatReceived time.Time
	sessionID             *string
	lastSequenceReceived  *int
	url                   *string
}

func (g *GatewayImpl) Logger() log.Logger {
	return g.logger
}

func (g *GatewayImpl) Config() Config {
	return g.config
}

// Open initializes the client and connection to discord
func (g *GatewayImpl) Open() error {
	if g.lastSequenceReceived == nil || g.sessionID == nil {
		g.status = StatusConnecting
	} else {
		g.status = StatusReconnecting
	}

	g.Logger().Info("starting ws...")

	if g.url == nil {
		g.Logger().Debug("gateway url empty, fetching...")
		gatewayRs, err := g.restServices.GatewayService().GetGateway(context.TODO())
		if err != nil {
			return err
		}
		g.url = &gatewayRs.URL
	}

	gatewayURL := *g.url + "?v=" + route.APIVersion + "&encoding=json"
	var rs *http.Response
	var err error
	g.conn, rs, err = websocket.DefaultDialer.Dial(gatewayURL, nil)
	if err != nil {
		g.Close()
		var body []byte
		if rs != nil && rs.Body != nil {
			body, err = ioutil.ReadAll(rs.Body)
			if err != nil {
				g.Logger().Errorf("error while reading response body: %s", err)
				g.url = nil
				return err
			}
		} else {
			body = []byte("null")
		}

		g.Logger().Errorf("error connecting to gateway. url: %s, error: %s, body: %s", gatewayURL, err, string(body))
		return err
	}

	g.conn.SetCloseHandler(func(code int, error string) error {
		g.Logger().Infof("connection to websocket closedwith code: %d, error: %s", code, error)
		return nil
	})

	g.status = StatusWaitingForHello

	mt, reader, err := g.conn.NextReader()
	if err != nil {
		return err
	}

	event, err := g.parseGatewayEvent(mt, reader)
	if err != nil {
		return err
	}

	if event.Op != OpHello {
		return discord.ErrUnexpectedGatewayOp(int(OpHello), int(event.Op))
	}

	g.lastHeartbeatReceived = time.Now().UTC()

	var eventData HelloGatewayEventData
	if err = json.Unmarshal(event.D, &eventData); err != nil {
		return err
	}

	g.heartbeatInterval = eventData.HeartbeatInterval * time.Millisecond

	if g.lastSequenceReceived == nil || g.sessionID == nil {
		g.status = StatusIdentifying
		g.Logger().Infof("sending StatusIdentifying command...")
		if err = g.Send(
			NewGatewayCommand(OpIdentify, IdentifyCommand{
				Token: g.token,
				Properties: IdentifyCommandDataProperties{
					OS:      g.config.OS,
					Browser: g.config.Browser,
					Device:  g.config.Device,
				},
				Compress:       false,
				LargeThreshold: g.config.LargeThreshold,
				GatewayIntents: g.config.GatewayIntents,
			}),
		); err != nil {
			return err
		}
		g.status = StatusWaitingForReady
	} else {
		g.status = StatusResuming
		cmd := NewGatewayCommand(OpResume, ResumeCommand{
			Token:     g.token,
			SessionID: *g.sessionID,
			Seq:       *g.lastSequenceReceived,
		})
		g.Logger().Infof("sending StatusResuming command...")

		if err = g.Send(cmd); err != nil {
			return err
		}
	}

	g.quit = make(chan struct{})

	go g.heartbeat()
	go g.listen()

	return nil
}

// Close cleans up the gateway internals
func (g *GatewayImpl) Close() {
	g.closeWithCode(websocket.CloseNormalClosure)
}

// Status returns the gateway connection status
func (g *GatewayImpl) Status() Status {
	return g.status
}

// Send sends a payload to the Gateway
func (g *GatewayImpl) Send(command GatewayCommand) error {
	return g.conn.WriteJSON(command)
}

// Latency returns the api.Gateway latency
func (g *GatewayImpl) Latency() time.Duration {
	return g.lastHeartbeatReceived.Sub(g.lastHeartbeatSent)
}

func (g *GatewayImpl) reconnect(delay time.Duration) {
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

func (g *GatewayImpl) closeWithCode(code int) {
	if g.quit != nil {
		g.Logger().Info("closing gateway goroutines...")
		close(g.quit)
		g.quit = nil
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

func (g *GatewayImpl) heartbeat() {
	defer func() {
		if r := recover(); r != nil {
			g.Logger().Panicf("recovered heartbeat goroutine error: %s", r)
			debug.PrintStack()
			g.heartbeat()
			return
		}
		g.Logger().Info("shut down heartbeat goroutine")
	}()

	ticker := time.NewTicker(g.heartbeatInterval)
	for {
		select {
		case <-ticker.C:
			g.sendHeartbeat()
		case <-g.quit:
			ticker.Stop()
			return
		}
	}
}

func (g *GatewayImpl) sendHeartbeat() {
	g.Logger().Debug("sending heartbeat...")
	/*
		heartbeatEvent := &events.HeartbeatEvent{
			GenericEvent: events.NewGenericEvent(g.Disgo(), 0),
			OldPing:      g.Latency(),
		}

		if err := g.Send(NewGatewayCommand(OpHeartbeat, g.lastSequenceReceived)); err != nil {
			g.Logger().Errorf("failed to send heartbeat with error: %s", err)
			g.closeWithCode(websocket.CloseServiceRestart)
			g.reconnect(1 * time.Second)
		}
		g.lastHeartbeatSent = time.Now().UTC()

		heartbeatEvent.NewPing = g.Latency()
		g.Disgo().EventManager().Dispatch(heartbeatEvent)*/
}

func (g *GatewayImpl) listen() {
	defer func() {
		if r := recover(); r != nil {
			g.Logger().Panicf("recovered listen goroutine error: %s", r)
			debug.PrintStack()
			g.listen()
			return
		}
		g.Logger().Info("shut down listen goroutine")
	}()

	for {
		select {
		case <-g.quit:
			g.Logger().Infof("existed listen routine")
			return
		default:
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
				g.Logger().Errorf("error while unpacking gateway event. error: %s", err)
			}

			switch op := event.Op; op {

			case OpDispatch:
				g.Logger().Debugf("received: OpDispatch")
				if event.S != nil {
					g.lastSequenceReceived = event.S
				}
				if event.T == nil {
					g.Logger().Errorf("received event without T. payload: %+v", event)
					continue
				}

				g.Logger().Debugf("received: %s", *event.T)

				if *event.T == EventTypeReady {
					var readyEvent ReadyGatewayEvent
					if err = g.parseEventToStruct(event, &readyEvent); err != nil {
						g.Logger().Errorf("Error parsing ready event: %s", err)
						continue
					}
					g.sessionID = &readyEvent.SessionID
					g.status = StatusWaitingForGuilds
					g.Logger().Info("ready event received")
				}

				g.eventHandlerFunc(*event.T, *event.S, bytes.NewBuffer(event.D))

			case OpHeartbeat:
				g.Logger().Debugf("received: OpHeartbeat")
				g.sendHeartbeat()

			case OpReconnect:
				g.Logger().Debugf("received: OpReconnect")
				g.closeWithCode(websocket.CloseServiceRestart)
				g.reconnect(1 * time.Second)

			case OpInvalidSession:
				var canResume bool
				if err = g.parseEventToStruct(event, &canResume); err != nil {
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

			case OpHeartbeatACK:
				g.Logger().Debugf("received: OpHeartbeatACK")
				g.lastHeartbeatReceived = time.Now().UTC()
			}

		}

	}
}

func (g *GatewayImpl) parseEventToStruct(event *RawEvent, v interface{}) error {
	if err := json.Unmarshal(event.D, v); err != nil {
		g.Logger().Errorf("error while unmarshalling event. error: %s", err)
		return err
	}
	return nil
}

func (g *GatewayImpl) parseGatewayEvent(mt int, reader io.Reader) (*RawEvent, error) {
	if mt == websocket.BinaryMessage {
		return nil, discord.ErrGatewayCompressedData
	}

	if mt != websocket.TextMessage {
		return nil, discord.ErrUnexpectedMessageType(mt)
	}

	decoder := json.NewDecoder(reader)
	var event RawEvent
	if err := decoder.Decode(&event); err != nil {
		g.Logger().Errorf("error decoding websocket message, %s", err)
		return nil, err
	}
	return &event, nil
}
