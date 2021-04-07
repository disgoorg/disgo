package api

import "time"

// ConnectionStatus is the state that the client is currently in
type ConnectionStatus int

// Indicates how far along the client is to connecting
const (
	Ready ConnectionStatus = iota
	Unconnected
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

// GatewayEventName wraps all GatewayEventName types
type GatewayEventName string

// Constants for the gateway events
const (
	GatewayEventApplicationCommandCreate GatewayEventName = "APPLICATION_COMMAND_CREATE"
	GatewayEventApplicationCommandUpdate GatewayEventName = "APPLICATION_COMMAND_UPDATE"
	GatewayEventApplicationCommandDelete GatewayEventName = "APPLICATION_COMMAND_DELETE"
	GatewayEventChannelCreate            GatewayEventName = "CHANNEL_CREATE"
	GatewayEventChannelDelete            GatewayEventName = "CHANNEL_DELETE"
	GatewayEventChannelPinsUpdate        GatewayEventName = "CHANNEL_PINS_UPDATE"
	GatewayEventChannelUpdate            GatewayEventName = "CHANNEL_UPDATE"
	GatewayEventGuildBanAdd              GatewayEventName = "GUILD_BAN_ADD"
	GatewayEventGuildBanRemove           GatewayEventName = "GUILD_BAN_REMOVE"
	GatewayEventGuildCreate              GatewayEventName = "GUILD_CREATE"
	GatewayEventGuildDelete              GatewayEventName = "GUILD_DELETE"
	GatewayEventGuildEmojisUpdate        GatewayEventName = "GUILD_EMOJIS_UPDATE"
	GatewayEventGuildIntegrationsUpdate  GatewayEventName = "GUILD_INTEGRATIONS_UPDATE"
	GatewayEventGuildMemberAdd           GatewayEventName = "GUILD_MEMBER_ADD"
	GatewayEventGuildMemberRemove        GatewayEventName = "GUILD_MEMBER_REMOVE"
	GatewayEventGuildMemberUpdate        GatewayEventName = "GUILD_MEMBER_UPDATE"
	GatewayEventGuildMembersChunk        GatewayEventName = "GUILD_MEMBERS_CHUNK"
	GatewayEventGuildRoleCreate          GatewayEventName = "GUILD_ROLE_CREATE"
	GatewayEventGuildRoleDelete          GatewayEventName = "GUILD_ROLE_DELETE"
	GatewayEventGuildRoleUpdate          GatewayEventName = "GUILD_ROLE_UPDATE"
	GatewayEventGuildUpdate              GatewayEventName = "GUILD_UPDATE"
	GatewayEventInteractionCreate        GatewayEventName = "INTERACTION_CREATE"
	WebhookEventInteractionCreate        GatewayEventName = "INTERACTION_WEBHOOK_CREATE"
	GatewayEventMessageAck               GatewayEventName = "MESSAGE_ACK"
	GatewayEventMessageCreate            GatewayEventName = "MESSAGE_CREATE"
	GatewayEventMessageDelete            GatewayEventName = "MESSAGE_DELETE"
	GatewayEventMessageDeleteBulk        GatewayEventName = "MESSAGE_DELETE_BULK"
	GatewayEventMessageReactionAdd       GatewayEventName = "MESSAGE_REACTION_ADD"
	GatewayEventMessageReactionRemove    GatewayEventName = "MESSAGE_REACTION_REMOVE"
	GatewayEventMessageReactionRemoveAll GatewayEventName = "MESSAGE_REACTION_REMOVE_ALL"
	GatewayEventMessageUpdate            GatewayEventName = "MESSAGE_UPDATE"
	GatewayEventPresenceUpdate           GatewayEventName = "PRESENCE_UPDATE"
	GatewayEventPresencesReplace         GatewayEventName = "PRESENCES_REPLACE"
	GatewayEventReady                    GatewayEventName = "READY"
	GatewayEventResumed                  GatewayEventName = "RESUMED"
	GatewayEventTypingStart              GatewayEventName = "TYPING_START"
	GatewayEventUserGuildSettingsUpdate  GatewayEventName = "USER_GUILD_SETTINGS_UPDATE"
	GatewayEventUserNoteUpdate           GatewayEventName = "USER_NOTE_UPDATE"
	GatewayEventUserSettingsUpdate       GatewayEventName = "USER_SETTINGS_UPDATE"
	GatewayEventUserUpdate               GatewayEventName = "USER_UPDATE"
	GatewayEventVoiceServerUpdate        GatewayEventName = "VOICE_SERVER_UPDATE"
	GatewayEventVoiceStateUpdate         GatewayEventName = "VOICE_STATE_UPDATE"
	GatewayEventWebhooksUpdate           GatewayEventName = "WEBHOOKS_UPDATE"
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
