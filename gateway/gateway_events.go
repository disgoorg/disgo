package gateway

import (
	"encoding/json"
	"time"

	"github.com/DisgoOrg/disgo/discord"
)

// Packet raw GatewayEvent type
type Packet struct {
	Op Op         `json:"op"`
	S  *int       `json:"s,omitempty"`
	T  *EventType `json:"t,omitempty"`
}

// RawEvent specifies the data for the GatewayCommand payload that is being sent
type RawEvent struct {
	Packet
	D json.RawMessage `json:"d"`
}

// ReadyGatewayEvent is the event sent by discord when you successfully Identify
type ReadyGatewayEvent struct {
	Version   int              `json:"v"`
	SelfUser  discord.SelfUser `json:"user"`
	Guilds    []discord.Guild  `json:"guilds"`
	SessionID string           `json:"session_id"`
	Shard     *[2]int          `json:"shard,omitempty"`
}

// HelloGatewayEventData is sent when we connect to the gateway
type HelloGatewayEventData struct {
	HeartbeatInterval time.Duration `json:"heartbeat_interval"`
}
