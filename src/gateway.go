package src

import (
	"time"

	"github.com/DiscoOrg/disgo/src/models"
	"github.com/gorilla/websocket"
)

// Gateway is what is used to connect to discord
type Gateway struct {
	wsConnection          *websocket.Conn
	heartbeatInterval     int
	lastHeartbeatSent     time.Time
	lastHeartbeatReceived time.Time
	sessionID             string
	lastSequenceReceived  int
}

type gatewayPayload struct {
	Op int         `json:"op"`
	D  interface{} `json:"d"`
	S  int         `json:"s,omitempty"`
	T  string      `json:"t,omitempty"`
}

type identifyPayload struct {
	Op int              `json:"op"`
	D  identifyDPayload `json:"d"`
}

type identifyDPayload struct {
	Token          string               `json:"token"`
	Properties     identifyPropsPayload `json:"properties"`
	Compress       bool                 `json:"compress,omitempty"`
	LargeThreshold int                  `json:"large_threshold,omitempty"`
	// Todo: Add shard property here, need to discuss
	// Todo: Add presence property here, need presence methods/struct
	GuildSubscriptions bool  `json:"guild_subscriptions,omitempty"`
	Intents            int64 `json:"intents"`
}

type identifyPropsPayload struct {
	OS      string `json:"$os"`      // user OS
	Browser string `json:"$browser"` // library name
	Device  string `json:"$device"`  // library name
}

type requestMembersPayload struct {
	GuildID   models.Snowflake   `json:"guild_id"`
	Query     string             `json:"query"` //If specified, user_ids must not be entered
	Limit     int                `json:"limit"` //Must be >=1 if query/user_ids is used, otherwise 0
	Presences bool               `json:"presences,omitempty"`
	UserIDs   []models.Snowflake `json:"user_ids"` //If specified, query must not be entered
	Nonce     string             `json:"nonce,omitempty"` //All responses are hashed with this nonce, optional
}

type voiceStateUpdatePayload struct {
	GuildID   models.Snowflake `json:"guild_id"`
	ChannelID models.Snowflake `json:"channel_id"`
	SelfMute  bool             `json:"self_mute"`
	SelfDeaf  bool             `json:"self_deaf"`
}
