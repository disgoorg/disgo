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

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"

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
	status                api.ConnectionStatus
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
			log.Error("tried to reconnect gateway while connecting/reconnecting")
			return
		}
		log.Info("reconnecting gateway...")
		if err := g.Open(); err != nil {
			log.Errorf("failed to reconnect gateway: %s", err)
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

	log.Info("starting ws...")

	if g.url == nil {
		log.Debug("gateway url empty, fetching...")
		gatewayRs := api.GatewayRs{}
		if err := g.Disgo().RestClient().Request(endpoints.GetGateway.Compile(), nil, &gatewayRs); err != nil {
			return err
		}
		g.url = &gatewayRs.URL
	}

	gatewayURL := *g.url + "?v=" + endpoints.APIVersion + "&encoding=json"
	wsConn, rs, err := websocket.DefaultDialer.Dial(gatewayURL, nil)
	if err != nil {
		g.Close()
		var body string
		if rs != nil && rs.Body != nil {
			rawBody, err := ioutil.ReadAll(rs.Body)
			if err != nil {
				log.Errorf("error while reading response body: %s", err)
				g.url = nil
				return err
			}
			body = string(rawBody)
		} else {
			body = "null"
		}

		log.Errorf("error connecting to gateway. url: %s, error: %s, body: %s", gatewayURL, err.Error(), body)
		return err
	}
	wsConn.SetCloseHandler(func(code int, error string) error {
		log.Errorf("connection to websocket closed with code: %d, error: %s", code, error)
		return nil
	})

	g.conn = wsConn
	g.status = api.WaitingForHello

	mt, data, err := g.conn.ReadMessage()
	if err != nil {
		return err
	}
	event, err := parseGatewayEvent(mt, data)
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
	} else {
		g.status = api.Resuming
		if err = wsConn.WriteJSON(
			api.NewGatewayCommand(api.OpIdentify, api.ResumeCommand{
				Token:     g.Disgo().Token(),
				SessionID: *g.sessionID,
				Seq:       *g.lastSequenceReceived,
			}),
		); err != nil {
			return err
		}
	}

	g.status = api.WaitingForReady
	g.quit = make(chan interface{})

	go g.heartbeat()
	go g.listen()

	return nil
}

// Status returns the gateway connection status
func (g *GatewayImpl) Status() api.ConnectionStatus {
	return g.status
}

// Latency returns the api.Gateway latency
func (g *GatewayImpl) Latency() time.Duration {
	return g.lastHeartbeatSent.Sub(g.lastHeartbeatReceived)
}

// Close cleans up the gateway internals
func (g *GatewayImpl) Close() {
	g.status = api.Disconnected
	if g.quit != nil {
		log.Info("closing gateway goroutines...")
		close(g.quit)
		g.quit = nil
		log.Info("closed gateway goroutines")
	}
	if g.conn != nil {
		if err := g.conn.Close(); err != nil {
			log.Errorf("error while closing wsconn: %s", err)
			g.conn = nil
		}
	}
}

func (g *GatewayImpl) heartbeat() {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("recovered heartbeat goroutine error: %s", r)
			debug.PrintStack()
			g.heartbeat()
			return
		}
		log.Info("shut down heartbeat goroutine")
	}()

	ticker := time.NewTicker(g.heartbeatInterval)
	for {
		select {
		case <-ticker.C:
			g.sendHeartbeat()
		case _, ok := <-g.quit:
			if !ok {
				ticker.Stop()
				return
			}
		}
	}
}

func (g *GatewayImpl) sendHeartbeat() {
	log.Debug("sending heartbeat...")

	if err := g.conn.WriteJSON(api.NewGatewayCommand(api.OpHeartbeat, g.lastSequenceReceived)); err != nil {
		log.Errorf("failed to send heartbeat with error: %s", err)
		g.Close()
		g.reconnect(1 * time.Second)
	}
	g.lastHeartbeatSent = time.Now().UTC()
}

func (g *GatewayImpl) listen() {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("recovered listen goroutine error: %s", r)
			debug.PrintStack()
			g.listen()
			return
		}
		log.Info("shut down listen goroutine")
	}()
	for {
		select {
		case _, ok := <-g.quit:
			if !ok {
				return
			}
		default:
			mt, data, err := g.conn.ReadMessage()
			if err != nil {
				log.Errorf("error while reading from ws. error: %s", err)
				g.Close()
				g.reconnect(1 * time.Second)
				return
			}

			event, err := parseGatewayEvent(mt, data)
			if err != nil {
				log.Errorf("error while unpacking gateway event. error: %s", err)
			}

			switch op := event.Op; op {

			case api.OpDispatch:
				log.Debugf("received: OpDispatch")
				if event.S != nil {
					g.lastSequenceReceived = event.S
				}

				log.Debugf("received: %s", *event.T)

				if event.T != nil && *event.T == api.GatewayEventReady {
					var readyEvent api.ReadyGatewayEvent
					if err := parseEventToStruct(event, &readyEvent); err != nil {
						log.Errorf("Error parsing ready event: %s", err)
						continue
					}
					g.sessionID = &readyEvent.D.SessionID

					log.Info("ready event received")
				}

				if event.T == nil {
					log.Errorf("received event without T. playload: %s", string(data))
					continue
				}
				d := g.Disgo()
				e := d.EventManager()
				e.Handle(*event.T, event.D, nil)

			case api.OpHeartbeat:
				log.Debugf("received: OpHeartbeat")
				g.sendHeartbeat()

			case api.OpReconnect:
				g.Close()
				g.reconnect(0)
				log.Debugf("received: OpReconnect")

			case api.OpInvalidSession:
				log.Debugf("received: OpInvalidSession")

			case api.OpHeartbeatACK:
				log.Debugf("received: OpHeartbeatACK")
				g.lastHeartbeatReceived = time.Now().UTC()
			}

		}

	}
}

func parseEventToStruct(event *api.RawGatewayEvent, v interface{}) error {
	if err := json.Unmarshal(event.D, v); err != nil {
		log.Errorf("error while unmarshaling event. error: %s", err)
		return err
	}
	return nil
}

func parseGatewayEvent(mt int, data []byte) (*api.RawGatewayEvent, error) {

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
		log.Errorf("error decoding websocket message_events, %s", err)
		return nil, err
	}
	return &event, nil
}
