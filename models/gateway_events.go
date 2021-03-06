package models

import (
	"encoding/json"

	"github.com/DiscoOrg/disgo/constants"
)

type UnresolvedGatewayEvent struct {
	Op constants.GatewayOp `json:"op"`
	S  *int                `json:"s,omitempty"`
	T  *string             `json:"t,omitempty"`
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
	Intents            Intent                   `json:"intents"`
	// Todo: Add presence property here, need presence methods/struct
	// Todo: Add shard property here, need to discuss
}

type OpIdentifyDataProperties struct {
	OS      string `json:"$os"`      // user OS
	Browser string `json:"$browser"` // library name
	Device  string `json:"$device"`  // library name
}

type ReadyEvent struct {
	UnresolvedGatewayEvent
	D ReadyEventData `json:"d"`
}

type ReadyEventData struct {
	User User `json:"user"`
	PrivateChannels []PrivateChannel `json:"channel"`
	Guilds []Guild `json:"guild"`
	SessionID string `json:"session_id"`
	Shard [2]int `json:"shard,omitempty"`
}

type HeartbeatEvent struct {
	UnresolvedGatewayEvent
	D *int `json:"d"`
}

type RequestMembersPayload struct {
	GuildID   Snowflake   `json:"guild_id"`
	Query     string      `json:"query"` //If specified, user_ids must not be entered
	Limit     int         `json:"limit"` //Must be >=1 if query/user_ids is used, otherwise 0
	Presences bool        `json:"presences,omitempty"`
	UserIDs   []Snowflake `json:"user_ids"`        //If specified, query must not be entered
	Nonce     string      `json:"nonce,omitempty"` //All responses are hashed with this nonce, optional
}

type VoiceStateUpdatePayload struct {
	GuildID   Snowflake `json:"guild_id"`
	ChannelID Snowflake `json:"channel_id"`
	SelfMute  bool      `json:"self_mute"`
	SelfDeaf  bool      `json:"self_deaf"`
}
