package src

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"runtime"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"

	"github.com/DiscoOrg/disgo/src/endpoints"
	"github.com/DiscoOrg/disgo/src/models"
)

// Gateway is what is used to connect to discord
type Gateway struct {
	Disgo                 *Disgo
	conn                  *websocket.Conn
	connectionStatus      ConnectionStatus
	heartbeatInterval     int
	lastHeartbeatSent     time.Time
	lastHeartbeatReceived time.Time
	sessionID             string
	lastSequenceReceived  *int
	url                   *string
	messageChannel        chan interface{}
}

func (g Gateway) Open() error {
	g.connectionStatus = Initializing
	log.Info("starting ws...")

	gatewayBase := "wss://gateway.discord.gg"
	g.url = &gatewayBase

	if g.url == nil {
		log.Println("Gateway url is nil, fetching...")
		gatewayRs := models.GatewayRs{}
		if err := g.Disgo.RestClient.Request(endpoints.Gateway, nil, &gatewayRs); err != nil {
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

	mt, data, err := g.conn.ReadMessage()
	if err != nil {
		return err
	}
	event, err := parseGatewayEvent(mt, data)
	if err != nil {
		return err
	}
	if event.Op != OpHello {
		return fmt.Errorf("expected op: hello type: 10, received: %d", mt)
	}

	g.lastHeartbeatReceived = time.Now().UTC()

	var eventData HelloEvent
	if err = json.Unmarshal(event.D, &eventData); err != nil {
		return err
	}

	g.heartbeatInterval = eventData.HeartbeatInterval

	if err = wsConn.WriteJSON(IdentifyEvent{
		UnresolvedGatewayEvent: UnresolvedGatewayEvent{
			Op: OpIdentify,
		},
		D: IdentifyEventData{
			Token: g.Disgo.Token,
			Properties: OpIdentifyDataProperties{
				OS:      getOS(),
				Browser: "disgo",
				Device:  "disgo",
			},
			Compress:       false,
			LargeThreshold: 50,
			Intents:        int64(g.Disgo.Intents),
		},
	}); err != nil {
		return err
	}

	g.messageChannel = make(chan interface{})

	go g.heartbeat()
	go g.listen()

	return nil
}

func (g Gateway) heartbeat() {
	for {
		time.Sleep(time.Duration(g.heartbeatInterval) * time.Millisecond)
		g.sendHeartbeat()
	}
}

func (g Gateway) sendHeartbeat() {
	log.Info("sending heartbeat...")

	if err := g.conn.WriteJSON(HeartbeatEvent{
		UnresolvedGatewayEvent: UnresolvedGatewayEvent{
			Op: OpHeartbeat,
		},
		D: g.lastSequenceReceived,
	}); err != nil {
		log.Errorf("failed to send heartbeat with error: %s", err)
	}
	g.lastHeartbeatSent = time.Now().UTC()
}

func (g Gateway) listen() {
	for {
		mt, data, err := g.conn.ReadMessage()
		if err != nil {
			log.Errorf("error while reading from ws. error: %s", err)
		}

		log.Info(string(data))

		event, err := parseGatewayEvent(mt, data)
		if err != nil {
			log.Errorf("error while unpacking gateway event. error: %s", err)
		}

		switch op := event.Op; op {

		case OpDispatch:
			if event.S != nil {
				g.lastSequenceReceived = event.S
			}
			// Todo: handle the thingy here!

		case OpHeartbeat:
			g.sendHeartbeat()

		case OpHeartbeatACK:
			g.lastHeartbeatReceived = time.Now().UTC()
		}
	}
}

func parseGatewayEvent(mt int, data []byte) (*GatewayEvent, error) {

	var reader io.Reader = bytes.NewBuffer(data)

	if mt == websocket.BinaryMessage {
		return nil, errors.New("we don't handle compressed yet")
	}
	if mt != websocket.TextMessage {
		return nil, fmt.Errorf("recieved unexpected message type: %d", mt)
	}
	var event GatewayEvent

	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&event); err != nil {
		log.Errorf("error decoding websocket message, %s", err)
		return nil, err
	}
	return &event, nil
}

type GatewayOp int

const (
	OpDispatch GatewayOp = iota
	OpHeartbeat
	OpIdentify
	OpPresenceUpdate
	OpVoiceStateUpdate
	_
	OpResume
	OpReconnect
	OpRequestGuildMembers
	OpInvalidSession
	OpHello
	OpHeartbeatACK
)

type ConnectionStatus int

const (
	Initializing ConnectionStatus = iota
	Initialized
	LoggingIn
	ConnectingToWebsocket
	IdentifyingSession
	AwaitingLoginInformation
	LoadingSubsystems
	Connected
	Disconnected
	ReconnectQueued
	WaitingToReconnect
	AttemptingToReconnect
	ShuttingDown
	Shutdown
	FailedToLogin
)

type UnresolvedGatewayEvent struct {
	Op GatewayOp `json:"op"`
	S  *int      `json:"s,omitempty"`
	T  *string   `json:"t,omitempty"`
}

type GatewayEvent struct {
	UnresolvedGatewayEvent
	D json.RawMessage `json:"d"`
}

type HelloEvent struct {
	HeartbeatInterval int `json:"heartbeat_interval"`
}

type IdentifyEvent struct {
	UnresolvedGatewayEvent
	D IdentifyEventData `json:"d"`
}

type IdentifyEventData struct {
	Token              string                   `json:"token"`
	Properties         OpIdentifyDataProperties `json:"properties"`
	Compress           bool                     `json:"compress,omitempty"`
	LargeThreshold     int                      `json:"large_threshold,omitempty"`
	GuildSubscriptions bool                     `json:"guild_subscriptions,omitempty"` // Deprecated, should not be specified when using intents
	Intents            int64                    `json:"intents"`
	// Todo: Add presence property here, need presence methods/struct
	// Todo: Add shard property here, need to discuss
}

type OpIdentifyDataProperties struct {
	OS      string `json:"$os"`      // user OS
	Browser string `json:"$browser"` // library name
	Device  string `json:"$device"`  // library name
}

type HeartbeatEvent struct {
	UnresolvedGatewayEvent
	D *int `json:"d"`
}

type requestMembersPayload struct {
	GuildID   models.Snowflake   `json:"guild_id"`
	Query     string             `json:"query"` //If specified, user_ids must not be entered
	Limit     int                `json:"limit"` //Must be >=1 if query/user_ids is used, otherwise 0
	Presences bool               `json:"presences,omitempty"`
	UserIDs   []models.Snowflake `json:"user_ids"`        //If specified, query must not be entered
	Nonce     string             `json:"nonce,omitempty"` //All responses are hashed with this nonce, optional
}

type voiceStateUpdatePayload struct {
	GuildID   models.Snowflake `json:"guild_id"`
	ChannelID models.Snowflake `json:"channel_id"`
	SelfMute  bool             `json:"self_mute"`
	SelfDeaf  bool             `json:"self_deaf"`
}

func getOS() string {
	OS := runtime.GOOS
	if strings.HasPrefix(OS, "windows") {
		return "windows"
	}
	if strings.HasPrefix(OS, "darwin") {
		return "darwin"
	}
	return "linux"
}
