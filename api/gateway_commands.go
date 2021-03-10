package api

import (
	"encoding/json"
)

type GatewayCommand struct {
	Op GatewayOp `json:"op"`
	S  *int      `json:"s,omitempty"`
	T  *string   `json:"t,omitempty"`
}

type RawGatewayCommand struct {
	GatewayCommand
	D json.RawMessage `json:"d"`
}

type IdentifyCommand struct {
	GatewayCommand
	D IdentifyCommandData `json:"d"`
}

type IdentifyCommandData struct {
	Token              string                        `json:"token"`
	Properties         IdentifyCommandDataProperties `json:"properties"`
	Compress           bool                          `json:"compress,omitempty"`
	LargeThreshold     int                           `json:"large_threshold,omitempty"`
	GuildSubscriptions bool                          `json:"guild_subscriptions,omitempty"` // Deprecated, should not be specified when using intents
	Intents            Intent                        `json:"intents"`
	// Todo: Add presence property here, need presence methods/struct
	// Todo: Add shard property here, need to discuss
}

type IdentifyCommandDataProperties struct {
	OS      string `json:"$os"`      // user OS
	Browser string `json:"$browser"` // library name
	Device  string `json:"$device"`  // library name
}

type ResumeCommand struct {
	GatewayCommand
	D struct {
		Token     string `json:"token"`
		SessionID string `json:"session_id"`
		Seq       int    `json:"seq"`
	} `json:"d"`
}

type HeartbeatCommand struct {
	GatewayCommand
	D *int `json:"d"`
}

type RequestGuildMembersCommand struct {
	GatewayCommand
	D RequestGuildMembersCommandData `json:"d"`
}

type RequestGuildMembersCommandData struct {
	GuildID   Snowflake   `json:"guild_id"`
	Query     string      `json:"query"` //If specified, user_ids must not be entered
	Limit     int         `json:"limit"` //Must be >=1 if query/user_ids is used, otherwise 0
	Presences bool        `json:"presences,omitempty"`
	UserIDs   []Snowflake `json:"user_ids"`        //If specified, query must not be entered
	Nonce     string      `json:"nonce,omitempty"` //All responses are hashed with this nonce, optional
}

type UpdateVoiceStateCommand struct {
	GatewayCommand
	D UpdateVoiceStateCommandData `json:"d"`
}

type UpdateVoiceStateCommandData struct {
	GuildID   Snowflake `json:"guild_id"`
	ChannelID Snowflake `json:"channel_id"`
	SelfMute  bool      `json:"self_mute"`
	SelfDeaf  bool      `json:"self_deaf"`
}

type UpdateStatusCommand struct {
	GatewayCommand
	D UpdateStatusCommandData `json:"d"`
}

type UpdateStatusCommandData struct {
	Since      *int       `json:"since"`
	Activities []Activity `json:"activities"`
	Status     bool       `json:"status"`
	AFK        bool       `json:"afk"`
}

type HelloCommand struct {
	GatewayCommand
	HeartbeatInterval int `json:"heartbeat_interval"`
}
