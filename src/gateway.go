package src

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"

	"github.com/DiscoOrg/disgo"
	"github.com/DiscoOrg/disgo/src/endpoints"
	"github.com/DiscoOrg/disgo/src/models"
)

// Gateway is what is used to connect to discord
type Gateway struct {
	disgo                 disgo.Disgo
	wsConnection          *websocket.Conn
	heartbeatInterval     int
	lastHeartbeatSent     time.Time
	lastHeartbeatReceived time.Time
	sessionID             string
	lastSequenceReceived  int
	gateway               *string
}

func (g Gateway) Open() error {
	log.Info("starting ws...")
	if g.gateway == nil {
		gatewayRs := models.GatewayRs{}
		if err := g.disgo.RestClient.Request(endpoints.Gateway, nil, gatewayRs); err != nil {
			return err
		}
		g.gateway = &gatewayRs.URL
	}

	gatewayUrl := *g.gateway + "?v=" + endpoints.APIVersion + "&encoding=json"
	wsConn, _, err := websocket.DefaultDialer.Dial(gatewayUrl, nil)
	if err != nil {
		log.Error("error connecting to gateway. url: %s, error: %s", gatewayUrl, err)
		return err
	}
	wsConn.SetCloseHandler(func(code int, error string) error {
		log.Error("connection to websocket closed with code: %d, error: %s", code, error)
		return nil
	})

	g.wsConnection = wsConn

	mt, data, err := g.wsConnection.ReadMessage()
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

	var eventData OpHelloData
	if err = json.Unmarshal(event.D, &eventData); err != nil {
		return err
	}

	g.heartbeatInterval = eventData.HeartbeatInterval

	wsConn.WriteJSON(GatewayEvent{
		Op: OpIdentify,
		D: ,
	})

	return nil
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
	OpDispatch            GatewayOp = 0
	OpHeartbeat           GatewayOp = 1
	OpIdentify            GatewayOp = 2
	OpPresenceUpdate      GatewayOp = 3
	OpVoiceStateUpdate    GatewayOp = 4
	OpResume              GatewayOp = 6
	OpReconnect           GatewayOp = 7
	OpRequestGuildMembers GatewayOp = 8
	OpInvalidSession      GatewayOp = 9
	OpHello               GatewayOp = 10
	OpHeartbeatACK        GatewayOp = 11
)

type GatewayEvent struct {
	Op GatewayOp       `json:"op"`
	D  json.RawMessage `json:"d"`
	S  int             `json:"s,omitempty"`
	T  string          `json:"t,omitempty"`
}

type OpHelloData struct {
	HeartbeatInterval int `json:"heartbeat_interval"`
}

type OpIdentifyData struct {
	Token              string                   `json:"token"`
	Properties         OpIdentifyDataProperties `json:"properties"`
	Compress           bool                     `json:"compress,omitempty"`
	LargeThreshold     int                      `json:"large_threshold,omitempty"`
	GuildSubscriptions bool                     `json:"guild_subscriptions,omitempty"`
	Intents            int64                    `json:"intents"`
	// Todo: Add presence property here, need presence methods/struct
	// Todo: Add shard property here, need to discuss
}

type OpIdentifyDataProperties struct {
	OS      string `json:"$os"`      // user OS
	Browser string `json:"$browser"` // library name
	Device  string `json:"$device"`  // library name
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
