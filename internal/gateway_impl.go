package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"runtime/debug"
	"time"

	"github.com/DisgoOrg/disgo/api/events"
	"github.com/gorilla/websocket"

	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/endpoints"
)

func newGatewayImpl(disgo api.Disgo) api.Gateway {
	return &GatewayImpl{
		disgo:  disgo,
		status: api.Unconnected,
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

		if g.Status() == api.Connecting || g.Status() == api.Reconnecting {
			g.Disgo().Logger().Error("tried to reconnect gateway while connecting/reconnecting")
			return
		}
		g.Disgo().Logger().Info("reconnecting gateway...")
		if err := g.Open(); err != nil {
			g.Disgo().Logger().Errorf("failed to reconnect gateway: %s", err)
			g.status = api.Disconnected
			g.reconnect(delay * 2)
		}
	}()
}

// Open initializes the client and connection to discord
func (g *GatewayImpl) Open() error {
	if g.lastSequenceReceived == nil || g.sessionID == nil {
		g.status = api.Connecting
	} else {
		g.status = api.Reconnecting
	}

	g.Disgo().Logger().Info("starting ws...")

	if g.url == nil {
		g.Disgo().Logger().Debug("gateway url empty, fetching...")
		gatewayRs := api.GatewayRs{}
		compiledRoute, err := endpoints.GetGateway.Compile()
		if err != nil {
			return err
		}
		if err = g.Disgo().RestClient().Request(*compiledRoute, nil, &gatewayRs); err != nil {
			return err
		}
		g.url = &gatewayRs.URL
	}

	gatewayURL := *g.url + "?v=" + endpoints.APIVersion + "&encoding=json"
	wsConn, rs, err := websocket.DefaultDialer.Dial(gatewayURL, nil)
	if err != nil {
		g.Close(false)
		var body string
		if rs != nil && rs.Body != nil {
			rawBody, err := ioutil.ReadAll(rs.Body)
			if err != nil {
				g.Disgo().Logger().Errorf("error while reading response body: %s", err)
				g.url = nil
				return err
			}
			body = string(rawBody)
		} else {
			body = "null"
		}

		g.Disgo().Logger().Errorf("error connecting to gateway. url: %s, error: %s, body: %s", gatewayURL, err.Error(), body)
		return err
	}
	wsConn.SetCloseHandler(func(code int, error string) error {
		g.Disgo().Logger().Infof("connection to websocket closed with code: %d, error: %s", code, error)
		return nil
	})

	g.conn = wsConn
	g.status = api.WaitingForHello

	mt, data, err := g.conn.ReadMessage()
	if err != nil {
		return err
	}
	event, err := g.parseGatewayEvent(mt, data)
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
		g.status = api.Identifying
		g.Disgo().Logger().Infof("sending Identifying command...")
		if err = wsConn.WriteJSON(
			api.NewGatewayCommand(api.OpIdentify, api.IdentifyCommand{
				Token: g.Disgo().Token(),
				Properties: api.IdentifyCommandDataProperties{
					OS:      api.GetOS(),
					Browser: "disgo",
					Device:  "disgo",
				},
				Compress:       false,
				LargeThreshold: g.Disgo().LargeThreshold(),
				Intents:        g.Disgo().Intents(),
			}),
		); err != nil {
			return err
		}
		g.status = api.WaitingForReady
	} else {
		g.Disgo().Logger().Infof("sending Resuming command...")
		g.status = api.Resuming
		if err = wsConn.WriteJSON(
			api.NewGatewayCommand(api.OpResume, api.ResumeCommand{
				Token:     g.Disgo().Token(),
				SessionID: *g.sessionID,
				Seq:       *g.lastSequenceReceived,
			}),
		); err != nil {
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

// Latency returns the api.Gateway latency
func (g *GatewayImpl) Latency() time.Duration {
	return g.lastHeartbeatReceived.Sub(g.lastHeartbeatSent)
}

// Close cleans up the gateway internals
func (g *GatewayImpl) Close(toReconnect bool) {
	g.status = api.Disconnected
	if g.quit != nil {
		g.Disgo().Logger().Info("closing gateway goroutines...")
		close(g.quit)
		g.quit = nil
		g.Disgo().Logger().Info("closed gateway goroutines")
	}
	if g.conn != nil {
		var err error
		if toReconnect {
			err = g.conn.Close()
			g.conn = nil
		} else {
			err = g.closeWithCode(websocket.CloseNormalClosure)
		}
		if err != nil {
			g.Disgo().Logger().Errorf("error while closing wsconn: %s", err)
		}
	}
}

func (g *GatewayImpl) closeWithCode(code int) error {
	if g.conn != nil {

		err := g.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(code, ""))
		if err != nil {
			return err
		}

		// TODO: Wait for Discord to actually close the connection.
		time.Sleep(1 * time.Second)

		err = g.conn.Close()
		g.conn = nil
		if err != nil {
			return err
		}

	}
	return nil
}

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
	//g.Disgo().Logger().Debug("sending heartbeat...")

	heartbeatEvent := events.HeartbeatEvent{
		GenericEvent: events.NewEvent(g.Disgo(), 0),
		OldPing:      g.Latency(),
	}

	if err := g.conn.WriteJSON(api.NewGatewayCommand(api.OpHeartbeat, g.lastSequenceReceived)); err != nil {
		g.Disgo().Logger().Errorf("failed to send heartbeat with error: %s", err)
		g.Close(true)
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
			mt, data, err := g.conn.ReadMessage()
			if err != nil {
				g.Disgo().Logger().Errorf("error while reading from ws. error: %s", err)
				g.Close(true)
				g.reconnect(1 * time.Second)
				return
			}

			event, err := g.parseGatewayEvent(mt, data)
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
					g.Disgo().Logger().Errorf("received event without T. playload: %s", string(data))
					continue
				}

				g.Disgo().Logger().Debugf("received: %s", *event.T)

				if *event.T == api.GatewayEventReady {
					var readyEvent api.ReadyGatewayEvent
					if err := g.parseEventToStruct(event, &readyEvent); err != nil {
						g.Disgo().Logger().Errorf("Error parsing ready event: %s", err)
						continue
					}
					g.sessionID = &readyEvent.D.SessionID

					g.Disgo().Logger().Info("ready event received")
				}

				// TODO: add setting to enable raw gateway events?
				var payload map[string]interface{}
				if err = g.parseEventToStruct(event, &payload); err != nil {
					g.Disgo().Logger().Errorf("Error parsing event: %s", err)
					continue
				}
				g.Disgo().EventManager().Dispatch(events.RawGatewayEvent{
					GenericEvent: events.NewEvent(g.Disgo(), *event.S),
					Type:         *event.T,
					RawPayload:   event.D,
					Payload:      payload,
				})

				d := g.Disgo()
				e := d.EventManager()
				e.Handle(*event.T, nil, *event.S, event.D)

			case api.OpHeartbeat:
				//g.Disgo().Logger().Debugf("received: OpHeartbeat")
				g.sendHeartbeat()

			case api.OpReconnect:
				g.Disgo().Logger().Debugf("received: OpReconnect")
				g.Close(true)
				g.reconnect(1 * time.Second)

			case api.OpInvalidSession:
				g.Disgo().Logger().Debugf("received: OpInvalidSession")
				g.Close(false)
				// clear reconnect info
				g.sessionID = nil
				g.lastSequenceReceived = nil
				g.reconnect(5 * time.Second)

			case api.OpHeartbeatACK:
				//g.Disgo().Logger().Debugf("received: OpHeartbeatACK")
				g.lastHeartbeatReceived = time.Now().UTC()
			}

		}

	}
}

func (g *GatewayImpl) parseEventToStruct(event *api.RawGatewayEvent, v interface{}) error {
	if err := json.Unmarshal(event.D, v); err != nil {
		g.Disgo().Logger().Errorf("error while unmarshaling event. error: %s", err)
		return err
	}
	return nil
}

func (g *GatewayImpl) parseGatewayEvent(mt int, data []byte) (*api.RawGatewayEvent, error) {

	var reader io.Reader = bytes.NewBuffer(data)

	if mt == websocket.BinaryMessage {
		return nil, errors.New("we don't handle compressed yet")
	}
	if mt != websocket.TextMessage {
		return nil, fmt.Errorf("recieved unexpected message_events type: %d", mt)
	}
	var event api.RawGatewayEvent

	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&event); err != nil {
		g.Disgo().Logger().Errorf("error decoding websocket message_events, %s", err)
		return nil, err
	}
	return &event, nil
}
