package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/DisgoOrg/disgo/api/events"
	"github.com/gorilla/websocket"

	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/restclient"
)

func newGatewayImpl(disgo api.Disgo) api.Gateway {
	return &GatewayImpl{
		disgo:  disgo,
		status: api.GatewayStatusUnconnected,
	}
}

// GatewayImpl is what is used to connect to discord
type GatewayImpl struct {
	disgo                 api.Disgo
	conn                  *websocket.Conn
	quit                  chan interface{}
	status                api.GatewayStatus
	heartbeatInterval     time.Duration
	lastHeartbeatSent     time.Time
	lastHeartbeatReceived time.Time
	sessionID             *string
	lastSequenceReceived  *int
	url                   *string
}

// Disgo returns the gateway's disgo client
func (g *GatewayImpl) Disgo() api.Disgo {
	return g.disgo
}

func (g *GatewayImpl) reconnect(delay time.Duration) {
	go func() {
		time.Sleep(delay)

		if g.Status() == api.GatewayStatusConnecting || g.Status() == api.GatewayStatusReconnecting {
			g.Disgo().Logger().Error("tried to reconnect gateway while connecting/reconnecting")
			return
		}
		g.Disgo().Logger().Info("reconnecting gateway...")
		if err := g.Open(); err != nil {
			g.Disgo().Logger().Errorf("failed to reconnect gateway: %s", err)
			g.status = api.GatewayStatusDisconnected
			g.reconnect(delay * 2)
		}
	}()
}

// Open initializes the client and connection to discord
func (g *GatewayImpl) Open() error {
	if g.lastSequenceReceived == nil || g.sessionID == nil {
		g.status = api.GatewayStatusConnecting
	} else {
		g.status = api.GatewayStatusReconnecting
	}

	g.Disgo().Logger().Info("starting ws...")

	if g.url == nil {
		g.Disgo().Logger().Debug("gateway url empty, fetching...")
		gatewayRs := api.GatewayRs{}
		compiledRoute, err := restclient.GetGateway.Compile(nil)
		if err != nil {
			return err
		}
		if err = g.Disgo().RestClient().Do(compiledRoute, nil, &gatewayRs); err != nil {
			return err
		}
		g.url = &gatewayRs.URL
	}

	gatewayURL := *g.url + "?v=" + restclient.APIVersion + "&encoding=json"
	var rs *http.Response
	var err error
	g.conn, rs, err = websocket.DefaultDialer.Dial(gatewayURL, nil)
	if err != nil {
		g.Close()
		var body []byte
		if rs != nil && rs.Body != nil {
			body, err = ioutil.ReadAll(rs.Body)
			if err != nil {
				g.Disgo().Logger().Errorf("error while reading response body: %s", err)
				g.url = nil
				return err
			}
		} else {
			body = []byte("null")
		}

		g.Disgo().Logger().Errorf("error connecting to gateway. url: %s, error: %s, body: %s", gatewayURL, err, string(body))
		return err
	}

	g.conn.SetCloseHandler(func(code int, error string) error {
		g.Disgo().Logger().Infof("connection to websocket closed with code: %d, error: %s", code, error)
		return nil
	})

	g.status = api.GatewayStatusWaitingForHello

	mt, reader, err := g.conn.NextReader()
	if err != nil {
		return err
	}

	event, err := g.parseGatewayEvent(mt, reader)
	if err != nil {
		return err
	}

	if event.Op != api.OpHello {
		return fmt.Errorf("expected op: hello type: 10, received: %d", mt)
	}

	g.lastHeartbeatReceived = time.Now().UTC()

	var eventData api.HelloGatewayEventData
	if err = json.Unmarshal(event.D, &eventData); err != nil {
		return err
	}

	g.heartbeatInterval = eventData.HeartbeatInterval * time.Millisecond

	if g.lastSequenceReceived == nil || g.sessionID == nil {
		g.status = api.GatewayStatusIdentifying
		g.Disgo().Logger().Infof("sending GatewayStatusIdentifying command...")
		if err = g.Send(
			api.NewGatewayCommand(api.OpIdentify, api.IdentifyCommand{
				Token: g.Disgo().Token(),
				Properties: api.IdentifyCommandDataProperties{
					OS:      api.GetOS(),
					Browser: "disgo",
					Device:  "disgo",
				},
				Compress:       false,
				LargeThreshold: g.Disgo().LargeThreshold(),
				GatewayIntents: g.Disgo().GatewayIntents(),
			}),
		); err != nil {
			return err
		}
		g.status = api.GatewayStatusWaitingForReady
	} else {
		g.status = api.GatewayStatusResuming
		cmd := api.NewGatewayCommand(api.OpResume, api.ResumeCommand{
			Token:     g.Disgo().Token(),
			SessionID: *g.sessionID,
			Seq:       *g.lastSequenceReceived,
		})
		g.Disgo().Logger().Infof("sending GatewayStatusResuming command...")

		if err = g.Send(cmd); err != nil {
			return err
		}
	}

	g.quit = make(chan interface{})

	go g.heartbeat()
	go g.listen()

	return nil
}

// Status returns the gateway connection status
func (g *GatewayImpl) Status() api.GatewayStatus {
	return g.status
}

// Send sends a payload to the Gateway
func (g *GatewayImpl) Send(command api.GatewayCommand) error {
	return g.conn.WriteJSON(command)
}

// Latency returns the api.Gateway latency
func (g *GatewayImpl) Latency() time.Duration {
	return g.lastHeartbeatReceived.Sub(g.lastHeartbeatSent)
}

// Close cleans up the gateway internals
func (g *GatewayImpl) Close() {
	g.closeWithCode(websocket.CloseNormalClosure)
}

func (g *GatewayImpl) closeWithCode(code int) {
	if g.quit != nil {
		g.Disgo().Logger().Info("closing gateway goroutines...")
		close(g.quit)
		g.quit = nil
		g.Disgo().Logger().Info("closed gateway goroutines")
	}
	if g.conn != nil {
		err := g.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(code, ""))
		if err != nil {
			g.Disgo().Logger().Errorf("error writing close code: %s", err)
		}

		// TODO: Wait for Discord to actually close the connection.
		time.Sleep(1 * time.Second)

		err = g.conn.Close()
		if err != nil {
			g.Disgo().Logger().Errorf("error closing conn: %s", err)
		}
		g.conn = nil
	}
}

// Conn returns the underlying websocket.Conn of this api.Gateway
func (g *GatewayImpl) Conn() *websocket.Conn {
	return g.conn
}

func (g *GatewayImpl) heartbeat() {
	defer func() {
		if r := recover(); r != nil {
			g.Disgo().Logger().Panicf("recovered heartbeat goroutine error: %s", r)
			debug.PrintStack()
			g.heartbeat()
			return
		}
		g.Disgo().Logger().Info("shut down heartbeat goroutine")
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
	g.Disgo().Logger().Debug("sending heartbeat...")

	heartbeatEvent := &events.HeartbeatEvent{
		GenericEvent: events.NewGenericEvent(g.Disgo(), 0),
		OldPing:      g.Latency(),
	}

	if err := g.Send(api.NewGatewayCommand(api.OpHeartbeat, g.lastSequenceReceived)); err != nil {
		g.Disgo().Logger().Errorf("failed to send heartbeat with error: %s", err)
		g.closeWithCode(websocket.CloseServiceRestart)
		g.reconnect(1 * time.Second)
	}
	g.lastHeartbeatSent = time.Now().UTC()

	heartbeatEvent.NewPing = g.Latency()
	g.Disgo().EventManager().Dispatch(heartbeatEvent)
}

func (g *GatewayImpl) listen() {
	defer func() {
		if r := recover(); r != nil {
			g.Disgo().Logger().Panicf("recovered listen goroutine error: %s", r)
			debug.PrintStack()
			g.listen()
			return
		}
		g.Disgo().Logger().Info("shut down listen goroutine")
	}()

	for {
		select {
		case <-g.quit:
			g.Disgo().Logger().Infof("existed listen routine")
			return
		default:
			if g.conn == nil {
				return
			}
			mt, reader, err := g.conn.NextReader()
			if err != nil {
				g.Disgo().Logger().Errorf("error while reading from ws. error: %s", err)
				g.closeWithCode(websocket.CloseServiceRestart)
				g.reconnect(1 * time.Second)
				return
			}

			event, err := g.parseGatewayEvent(mt, reader)
			if err != nil {
				g.Disgo().Logger().Errorf("error while unpacking gateway event. error: %s", err)
			}

			switch op := event.Op; op {

			case api.OpDispatch:
				g.Disgo().Logger().Debugf("received: OpDispatch")
				if event.S != nil {
					g.lastSequenceReceived = event.S
				}
				if event.T == nil {
					g.Disgo().Logger().Errorf("received event without T. payload: %+v", event)
					continue
				}

				g.Disgo().Logger().Debugf("received: %s", *event.T)

				if *event.T == api.GatewayEventReady {
					var readyEvent api.ReadyGatewayEvent
					if err = g.parseEventToStruct(event, &readyEvent); err != nil {
						g.Disgo().Logger().Errorf("Error parsing ready event: %s", err)
						continue
					}
					g.sessionID = &readyEvent.SessionID
					g.status = api.GatewayStatusWaitingForGuilds
					g.Disgo().Logger().Info("ready event received")
				}

				if g.Disgo().RawGatewayEventsEnabled() {
					var payload map[string]interface{}
					if err = g.parseEventToStruct(event, &payload); err != nil {
						g.Disgo().Logger().Errorf("Error parsing raw gateway event: %s", err)
					}
					g.Disgo().EventManager().Dispatch(&events.RawGatewayEvent{
						GenericEvent: events.NewGenericEvent(g.Disgo(), *event.S),
						Type:         *event.T,
						RawPayload:   event.D,
						Payload:      payload,
					})
				}

				d := g.Disgo()
				e := d.EventManager()
				e.Handle(*event.T, nil, *event.S, event.D)

			case api.OpHeartbeat:
				g.Disgo().Logger().Debugf("received: OpHeartbeat")
				g.sendHeartbeat()

			case api.OpReconnect:
				g.Disgo().Logger().Debugf("received: OpReconnect")
				g.closeWithCode(websocket.CloseServiceRestart)
				g.reconnect(1 * time.Second)

			case api.OpInvalidSession:
				var canResume bool
				if err = g.parseEventToStruct(event, &canResume); err != nil {
					g.Disgo().Logger().Errorf("Error parsing invalid session data: %s", err)
				}
				g.Disgo().Logger().Debugf("received: OpInvalidSession, canResume: %b", canResume)
				if canResume {
					g.closeWithCode(websocket.CloseServiceRestart)
				} else {
					g.Close()
					// clear reconnect info
					g.sessionID = nil
					g.lastSequenceReceived = nil
				}
				g.reconnect(5 * time.Second)

			case api.OpHeartbeatACK:
				g.Disgo().Logger().Debugf("received: OpHeartbeatACK")
				g.lastHeartbeatReceived = time.Now().UTC()
			}

		}

	}
}

func (g *GatewayImpl) parseEventToStruct(event *api.RawGatewayEvent, v interface{}) error {
	if err := json.Unmarshal(event.D, v); err != nil {
		g.Disgo().Logger().Errorf("error while unmarshalling event. error: %s", err)
		return err
	}
	return nil
}

func (g *GatewayImpl) parseGatewayEvent(mt int, reader io.Reader) (*api.RawGatewayEvent, error) {
	if mt == websocket.BinaryMessage {
		return nil, errors.New("we don't handle compressed yet")
	}

	if mt != websocket.TextMessage {
		return nil, fmt.Errorf("recieved unexpected message type: %d", mt)
	}

	decoder := json.NewDecoder(reader)
	var event api.RawGatewayEvent
	if err := decoder.Decode(&event); err != nil {
		g.Disgo().Logger().Errorf("error decoding websocket message, %s", err)
		return nil, err
	}
	return &event, nil
}
