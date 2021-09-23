package discord

// NewGatewayCommand returns a new GatewayCommand struct with the given payload
func NewGatewayCommand(op GatewayOpcode, d interface{}) GatewayCommand {
	return GatewayCommand{
		GatewayPayload: GatewayPayload{
			Op: op,
		},
		D: d,
	}
}

// GatewayCommand object is used when sending data to discord's websocket, it's recommended that you don't use these
//goland:noinspection GoNameStartsWithPackageName
type GatewayCommand struct {
	GatewayPayload
	D interface{} `json:"d"`
}

// IdentifyCommand is the data used in IdentifyCommand
type IdentifyCommand struct {
	Token          string                        `json:"token"`
	Properties     IdentifyCommandDataProperties `json:"properties"`
	Compress       bool                          `json:"compress,omitempty"`
	LargeThreshold int                           `json:"large_threshold,omitempty"`
	Shard          []int                         `json:"shard,omitempty"`
	GatewayIntents GatewayIntents                `json:"intents"`
	// Todo: Add presence property here, need presence methods/struct
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
	Token     string `json:"token"`
	SessionID string `json:"session_id"`
	Seq       int    `json:"seq"`
}

// HeartbeatCommand is used to ensure the websocket connection remains open, and disconnect if not.
type HeartbeatCommand struct {
	D *int `json:"d"`
}

// RequestGuildMembersCommand is used for fetching all of the members of a guild_events. It is recommended you have a strict
// member caching policy when using this.
type RequestGuildMembersCommand struct {
	GuildID   Snowflake   `json:"guild_id"`
	Query     string      `json:"query"` //If specified, user_ids must not be entered
	Limit     int         `json:"limit"` //Must be >=1 if query/user_ids is used, otherwise 0
	Presences bool        `json:"presences,omitempty"`
	UserIDs   []Snowflake `json:"user_ids"`        //If specified, query must not be entered
	Nonce     string      `json:"nonce,omitempty"` //All responses are hashed with this nonce, optional
}

// UpdateVoiceStateCommand is used for updating the bots voice state in a guild_events
type UpdateVoiceStateCommand struct {
	GuildID   Snowflake  `json:"guild_id"`
	ChannelID *Snowflake `json:"channel_id"`
	SelfMute  bool       `json:"self_mute"`
	SelfDeaf  bool       `json:"self_deaf"`
}

// UpdateStatusCommand is used for updating Bot's presence
type UpdateStatusCommand struct {
	Since      *int       `json:"since"`
	Activities []Activity `json:"activities"`
	Status     bool       `json:"status"`
	AFK        bool       `json:"afk"`
}
