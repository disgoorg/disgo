package api

import "time"

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
	Latency() time.Duration
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

// GatewayEvent wraps all GatewayEvent types
type GatewayEvent string

// Constants for the gateway events
const (
	GatewayEventChannelCreate            GatewayEvent = "CHANNEL_CREATE"
	GatewayEventChannelDelete            GatewayEvent = "CHANNEL_DELETE"
	GatewayEventChannelPinsUpdate        GatewayEvent = "CHANNEL_PINS_UPDATE"
	GatewayEventChannelUpdate            GatewayEvent = "CHANNEL_UPDATE"
	GatewayEventGuildBanAdd              GatewayEvent = "GUILD_BAN_ADD"
	GatewayEventGuildBanRemove           GatewayEvent = "GUILD_BAN_REMOVE"
	GatewayEventGuildCreate              GatewayEvent = "GUILD_CREATE"
	GatewayEventGuildDelete              GatewayEvent = "GUILD_DELETE"
	GatewayEventGuildEmojisUpdate        GatewayEvent = "GUILD_EMOJIS_UPDATE"
	GatewayEventGuildIntegrationsUpdate  GatewayEvent = "GUILD_INTEGRATIONS_UPDATE"
	GatewayEventGuildMemberAdd           GatewayEvent = "GUILD_MEMBER_ADD"
	GatewayEventGuildMemberRemove        GatewayEvent = "GUILD_MEMBER_REMOVE"
	GatewayEventGuildMemberUpdate        GatewayEvent = "GUILD_MEMBER_UPDATE"
	GatewayEventGuildMembersChunk        GatewayEvent = "GUILD_MEMBERS_CHUNK"
	GatewayEventGuildRoleCreate          GatewayEvent = "GUILD_ROLE_CREATE"
	GatewayEventGuildRoleDelete          GatewayEvent = "GUILD_ROLE_DELETE"
	GatewayEventGuildRoleUpdate          GatewayEvent = "GUILD_ROLE_UPDATE"
	GatewayEventGuildUpdate              GatewayEvent = "GUILD_UPDATE"
	GatewayEventInteractionCreate        GatewayEvent = "INTERACTION_CREATE"
	WebhookEventInteractionCreate        GatewayEvent = "INTERACTION_WEBHOOK_CREATE"
	GatewayEventMessageAck               GatewayEvent = "MESSAGE_ACK"
	GatewayEventMessageCreate            GatewayEvent = "MESSAGE_CREATE"
	GatewayEventMessageDelete            GatewayEvent = "MESSAGE_DELETE"
	GatewayEventMessageDeleteBulk        GatewayEvent = "MESSAGE_DELETE_BULK"
	GatewayEventMessageReactionAdd       GatewayEvent = "MESSAGE_REACTION_ADD"
	GatewayEventMessageReactionRemove    GatewayEvent = "MESSAGE_REACTION_REMOVE"
	GatewayEventMessageReactionRemoveAll GatewayEvent = "MESSAGE_REACTION_REMOVE_ALL"
	GatewayEventMessageUpdate            GatewayEvent = "MESSAGE_UPDATE"
	GatewayEventPresenceUpdate           GatewayEvent = "PRESENCE_UPDATE"
	GatewayEventPresencesReplace         GatewayEvent = "PRESENCES_REPLACE"
	GatewayEventReady                    GatewayEvent = "READY"
	GatewayEventResumed                  GatewayEvent = "RESUMED"
	GatewayEventTypingStart              GatewayEvent = "TYPING_START"
	GatewayEventUserGuildSettingsUpdate  GatewayEvent = "USER_GUILD_SETTINGS_UPDATE"
	GatewayEventUserNoteUpdate           GatewayEvent = "USER_NOTE_UPDATE"
	GatewayEventUserSettingsUpdate       GatewayEvent = "USER_SETTINGS_UPDATE"
	GatewayEventUserUpdate               GatewayEvent = "USER_UPDATE"
	GatewayEventVoiceServerUpdate        GatewayEvent = "VOICE_SERVER_UPDATE"
	GatewayEventVoiceStateUpdate         GatewayEvent = "VOICE_STATE_UPDATE"
	GatewayEventWebhooksUpdate           GatewayEvent = "WEBHOOKS_UPDATE"
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
