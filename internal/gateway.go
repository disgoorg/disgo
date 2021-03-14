package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"

	"github.com/DiscoOrg/disgo/api"
	"github.com/DiscoOrg/disgo/api/endpoints"
)

func newGatewayImpl(disgo api.Disgo) api.Gateway {
	return &GatewayImpl{
		disgo: disgo,
	}
}

// GatewayImpl is what is used to connect to discord
type GatewayImpl struct {
	disgo                 api.Disgo
	conn                  *websocket.Conn
	connectionStatus      api.ConnectionStatus
	heartbeatInterval     int
	lastHeartbeatSent     time.Time
	lastHeartbeatReceived time.Time
	sessionID             string
	lastSequenceReceived  *int
	url                   *string
}

// Disgo returns the gateway's disgo client
func (g GatewayImpl) Disgo() api.Disgo {
	return g.disgo
}

// Open initializes the client and connection to discord
func (g GatewayImpl) Open() error {
	g.connectionStatus = api.Connecting
	log.Info("starting ws...")

	gatewayBase := "wss://gateway.discord.gg"
	g.url = &gatewayBase

	if g.url == nil {
		log.Println("GetGateway url is nil, fetching...")
		gatewayRs := api.GatewayRs{}
		if err := g.Disgo().RestClient().Request(endpoints.GetGateway, nil, &gatewayRs); err != nil {
			return err
		}
		g.url = &gatewayRs.URL
	}

	gatewayUrl := *g.url + "?v=" + endpoints.APIVersion + "&encoding=json"
	wsConn, _, err := websocket.DefaultDialer.Dial(gatewayUrl, nil)
	if err != nil {
		log.Errorf("error connecting to gateway. url: %s, error: %s", gatewayUrl, err.Error())
		return err
	}
	wsConn.SetCloseHandler(func(code int, error string) error {
		log.Errorf("connection to websocket closed with code: %d, error: %s", code, error)
		return nil
	})

	g.conn = wsConn
	g.connectionStatus = api.WaitingForHello

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

	var eventData api.HelloEvent
	if err = json.Unmarshal(event.D, &eventData); err != nil {
		return err
	}

	g.connectionStatus = api.Identifying
	g.heartbeatInterval = eventData.HeartbeatInterval

	if err = wsConn.WriteJSON(api.IdentifyCommand{
		GatewayCommand: api.GatewayCommand{
			Op: api.OpIdentify,
		},
		D: api.IdentifyCommandData{
			Token: g.Disgo().Token(),
			Properties: api.IdentifyCommandDataProperties{
				OS:      api.GetOS(),
				Browser: "disgo",
				Device:  "disgo",
			},
			Compress:       false,
			LargeThreshold: 50,
			Intents:        g.Disgo().Intents(),
		},
	}); err != nil {
		return err
	}

	g.connectionStatus = api.WaitingForReady

	go g.heartbeat()
	go g.listen()

	return nil
}

func (g GatewayImpl) heartbeat() {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("recovered heartbeat goroutine error: %s", r)
			g.heartbeat()
			return
		}
		log.Info("shutting down heartbeat goroutine...")
	}()

	for {
		time.Sleep(time.Duration(g.heartbeatInterval) * time.Millisecond)
		g.sendHeartbeat()
	}
}

// Status returns the gateway connection status
func (g GatewayImpl) Status() api.ConnectionStatus {
	return g.connectionStatus
}

// Close cleans up the gateway internals
func (g GatewayImpl) Close() {
	log.Info("Implement closing smh...")
}

func (g GatewayImpl) sendHeartbeat() {
	log.Info("sending heartbeat...")

	err := g.conn.WriteJSON(api.HeartbeatCommand{
		GatewayCommand: api.GatewayCommand{
			Op: api.OpHeartbeat,
		},
		D: g.lastSequenceReceived,
	})
	if err != nil {
		log.Errorf("failed to send heartbeat with error: %s", err)
		_ = g.conn.Close()
		// Todo: reconnect
	}
	g.lastHeartbeatSent = time.Now().UTC()
}

func (g GatewayImpl) listen() {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("recovered listen goroutine error: %s", r)
			g.listen()
			return
		}
		log.Info("shutting down ws goroutine...")
	}()
	for {
		mt, data, err := g.conn.ReadMessage()
		if err != nil {
			log.Errorf("error while reading from ws. error: %s", err)
		}

		event, err := parseGatewayEvent(mt, data)
		if err != nil {
			log.Errorf("error while unpacking gateway event. error: %s", err)
		}

		switch op := event.Op; op {

		case api.OpDispatch:
			//log.Infof("received: OpDispatch")
			if event.S != nil {
				g.lastSequenceReceived = event.S
			}

			if event.T != nil && *event.T == "READY" {
				var readyEvent api.ReadyEventData
				if err := parseEventToStruct(event, &readyEvent); err != nil {
					return
				}
				g.sessionID = readyEvent.SessionID
				g.Disgo().SetSelfUser(readyEvent.User)
				log.Info("Client Ready")
			}

			if event.T == nil {
				log.Errorf("received event without T. playload: %s", string(data))
				continue
			}
			d := g.Disgo()
			e := d.EventManager()
			e.Handle(*event.T, event.D)

		case api.OpHeartbeat:
			log.Infof("received: OpHeartbeat")
			g.sendHeartbeat()

		case api.OpReconnect:
			log.Infof("received: OpReconnect")

		case api.OpInvalidSession:
			log.Infof("received: OpInvalidSession")

		case api.OpHeartbeatACK:
			log.Infof("received: OpHeartbeatACK")
			g.lastHeartbeatReceived = time.Now().UTC()
		}
	}
}

func parseEventToStruct(event *api.RawGatewayCommand, v interface{}) error {
	if err := json.Unmarshal(event.D, v); err != nil {
		log.Errorf("error while unmarshaling event. error: %s", err)
		return err
	}
	return nil
}

func parseGatewayEvent(mt int, data []byte) (*api.RawGatewayCommand, error) {

	var reader io.Reader = bytes.NewBuffer(data)

	if mt == websocket.BinaryMessage {
		return nil, errors.New("we don't handle compressed yet")
	}
	if mt != websocket.TextMessage {
		return nil, fmt.Errorf("recieved unexpected message_events type: %d", mt)
	}
	var event api.RawGatewayCommand

	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&event); err != nil {
		log.Errorf("error decoding websocket message_events, %s", err)
		return nil, err
	}
	return &event, nil
}
