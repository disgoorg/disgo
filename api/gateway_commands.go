package api

import (
	"encoding/json"
	"time"
)

// GatewayCommand object is used when sending data to discord's websocket, it's recommended that you don't use these
type GatewayCommand struct {
	Op GatewayOp `json:"op"`
	S  *int      `json:"s,omitempty"`
	T  *string   `json:"t,omitempty"`
}

// RawGatewayCommand specifies the data for the GatewayCommand payload that is being sent
type RawGatewayCommand struct {
	GatewayCommand
	D json.RawMessage `json:"d"`
}

// IdentifyCommand is used for Identifying to discord
type IdentifyCommand struct {
	GatewayCommand
	D IdentifyCommandData `json:"d"`
}

// IdentifyCommandData is the data used in IdentifyCommand
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

// IdentifyCommandDataProperties is used for specifying to discord which library and OS the bot is using, is
// automatically handled by the library and should rarely be used.
type IdentifyCommandDataProperties struct {
	OS      string `json:"$os"`      // user OS
	Browser string `json:"$browser"` // library name
	Device  string `json:"$device"`  // library name
}

// ResumeCommand is used to resume a connection to discord in the case that you are disconnected. Is automatically
// handled by the library and should rarely be used.
type ResumeCommand struct {
	GatewayCommand
	D struct {
		Token     string `json:"token"`
		SessionID string `json:"session_id"`
		Seq       int    `json:"seq"`
	} `json:"d"`
}

// HeartbeatCommand is used to ensure the websocket connection remains open, and disconnect if not.
type HeartbeatCommand struct {
	GatewayCommand
	D *int `json:"d"`
}

// RequestGuildMembersCommand is used for fetching all of the members of a guild_events. It is recommended you have a strict
// member caching policy when using this.
type RequestGuildMembersCommand struct {
	GatewayCommand
	D RequestGuildMembersCommandData `json:"d"`
}

// RequestGuildMembersCommandData is the RequestGuildMembersCommand.D payload
type RequestGuildMembersCommandData struct {
	GuildID   Snowflake   `json:"guild_id"`
	Query     string      `json:"query"` //If specified, user_ids must not be entered
	Limit     int         `json:"limit"` //Must be >=1 if query/user_ids is used, otherwise 0
	Presences bool        `json:"presences,omitempty"`
	UserIDs   []Snowflake `json:"user_ids"`        //If specified, query must not be entered
	Nonce     string      `json:"nonce,omitempty"` //All responses are hashed with this nonce, optional
}

// UpdateVoiceStateCommand is used for updating the bots voice state in a guild_events
type UpdateVoiceStateCommand struct {
	GatewayCommand
	D UpdateVoiceStateCommandData `json:"d"`
}

// UpdateVoiceStateCommandData is the UpdateVoiceStateCommand.D payload
type UpdateVoiceStateCommandData struct {
	GuildID   Snowflake `json:"guild_id"`
	ChannelID Snowflake `json:"channel_id"`
	SelfMute  bool      `json:"self_mute"`
	SelfDeaf  bool      `json:"self_deaf"`
}

// UpdateStatusCommand is used for updating Disgo's presence
type UpdateStatusCommand struct {
	GatewayCommand
	D UpdateStatusCommandData `json:"d"`
}

// UpdateStatusCommandData is the UpdateStatusCommand.D payload
type UpdateStatusCommandData struct {
	Since      *int       `json:"since"`
	Activities []Activity `json:"activities"`
	Status     bool       `json:"status"`
	AFK        bool       `json:"afk"`
}

// HelloEvent is used when
type HelloEvent struct {
	GatewayCommand
	HeartbeatInterval time.Duration `json:"heartbeat_interval"`
}
