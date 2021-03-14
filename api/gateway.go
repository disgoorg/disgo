package api

// ConnectionStatus is the state that the client is currently in
type ConnectionStatus int

// Indicates how far along the client is to connecting
const (
	Ready ConnectionStatus = iota
	Connecting
	Reconnecting
	WaitingForHello
	WaitingForReady
	Disconnected
	WaitingForGuilds
	Identifying
	Resuming
)

// Gateway is what is used to connect to discord
type Gateway interface {
	Disgo() Disgo
	Open() error
	Status() ConnectionStatus
	Close()
}


// GatewayOp are opcodes used by discord
type GatewayOp int

// Constants for the gateway opcodes
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

// Constants for the gateway events
const (
	ChannelCreateGatewayEvent            = "CHANNEL_CREATE"
	ChannelDeleteGatewayEvent            = "CHANNEL_DELETE"
	ChannelPinsUpdateGatewayEvent        = "CHANNEL_PINS_UPDATE"
	ChannelUpdateGatewayEvent            = "CHANNEL_UPDATE"
	GuildBanAddGatewayEvent              = "GUILD_BAN_ADD"
	GuildBanRemoveGatewayEvent           = "GUILD_BAN_REMOVE"
	GuildCreateGatewayEvent              = "GUILD_CREATE"
	GuildDeleteGatewayEvent              = "GUILD_DELETE"
	GuildEmojisUpdateGatewayEvent        = "GUILD_EMOJIS_UPDATE"
	GuildIntegrationsUpdateGatewayEvent  = "GUILD_INTEGRATIONS_UPDATE"
	GuildMemberAddGatewayEvent           = "GUILD_MEMBER_ADD"
	GuildMemberRemoveGatewayEvent        = "GUILD_MEMBER_REMOVE"
	GuildMemberUpdateGatewayEvent        = "GUILD_MEMBER_UPDATE"
	GuildMembersChunkGatewayEvent        = "GUILD_MEMBERS_CHUNK"
	GuildRoleCreateGatewayEvent          = "GUILD_ROLE_CREATE"
	GuildRoleDeleteGatewayEvent          = "GUILD_ROLE_DELETE"
	GuildRoleUpdateGatewayEvent          = "GUILD_ROLE_UPDATE"
	GuildUpdateGatewayEvent              = "GUILD_UPDATE"
	InteractionCreateGatewayEvent        = "INTERACTION_CREATE"
	MessageAckGatewayEvent               = "MESSAGE_ACK"
	MessageCreateGatewayEvent            = "MESSAGE_CREATE"
	MessageDeleteGatewayEvent            = "MESSAGE_DELETE"
	MessageDeleteBulkGatewayEvent        = "MESSAGE_DELETE_BULK"
	MessageReactionAddGatewayEvent       = "MESSAGE_REACTION_ADD"
	MessageReactionRemoveGatewayEvent    = "MESSAGE_REACTION_REMOVE"
	MessageReactionRemoveAllGatewayEvent = "MESSAGE_REACTION_REMOVE_ALL"
	MessageUpdateGatewayEvent            = "MESSAGE_UPDATE"
	PresenceUpdateGatewayEvent           = "PRESENCE_UPDATE"
	PresencesReplaceGatewayEvent         = "PRESENCES_REPLACE"
	ReadyGatewayEvent                    = "READY"
	ResumedGatewayEvent                  = "RESUMED"
	TypingStartGatewayEvent              = "TYPING_START"
	UserGuildSettingsUpdateGatewayEvent  = "USER_GUILD_SETTINGS_UPDATE"
	UserNoteUpdateGatewayEvent           = "USER_NOTE_UPDATE"
	UserSettingsUpdateGatewayEvent       = "USER_SETTINGS_UPDATE"
	UserUpdateGatewayEvent               = "USER_UPDATE"
	VoiceServerUpdateGatewayEvent        = "VOICE_SERVER_UPDATE"
	VoiceStateUpdateGatewayEvent         = "VOICE_STATE_UPDATE"
	WebhooksUpdateGatewayEvent           = "WEBHOOKS_UPDATE"
)


// GatewayRs contains the response for GET /gateway
type GatewayRs struct {
	URL string `json:"url"`
}

// GatewayBotRs contains the response for GET /gateway/bot
type GatewayBotRs struct {
	URL               string `json:"url"`
	Shards            int    `json:"shards"`
	SessionStartLimit struct {
		Total          int `json:"total"`
		Remaining      int `json:"remaining"`
		ResetAfter     int `json:"reset_after"`
		MaxConcurrency int `json:"max_concurrency"`
	} `json:"session_start_limit"`
}
