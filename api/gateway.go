package api

import (
	"time"

	"github.com/gorilla/websocket"
)

// GatewayStatus is the state that the client is currently in
type GatewayStatus int

// Indicates how far along the client is to connecting
const (
	Ready GatewayStatus = iota
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
	Status() GatewayStatus
	Close()
	Conn() *websocket.Conn
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

// GatewayEventType wraps all GatewayEventType types
type GatewayEventType string

// Constants for the gateway events
const (
	GatewayEventApplicationCommandCreate GatewayEventType = "APPLICATION_COMMAND_CREATE"
	GatewayEventApplicationCommandUpdate GatewayEventType = "APPLICATION_COMMAND_UPDATE"
	GatewayEventApplicationCommandDelete GatewayEventType = "APPLICATION_COMMAND_DELETE"
	GatewayEventChannelCreate            GatewayEventType = "CHANNEL_CREATE"
	GatewayEventChannelDelete            GatewayEventType = "CHANNEL_DELETE"
	GatewayEventChannelPinsUpdate        GatewayEventType = "CHANNEL_PINS_UPDATE"
	GatewayEventChannelUpdate            GatewayEventType = "CHANNEL_UPDATE"
	GatewayEventGuildBanAdd              GatewayEventType = "GUILD_BAN_ADD"
	GatewayEventGuildBanRemove           GatewayEventType = "GUILD_BAN_REMOVE"
	GatewayEventGuildCreate              GatewayEventType = "GUILD_CREATE"
	GatewayEventGuildDelete              GatewayEventType = "GUILD_DELETE"
	GatewayEventGuildEmojisUpdate        GatewayEventType = "GUILD_EMOJIS_UPDATE"
	GatewayEventGuildIntegrationsUpdate  GatewayEventType = "GUILD_INTEGRATIONS_UPDATE"
	GatewayEventGuildMemberAdd           GatewayEventType = "GUILD_MEMBER_ADD"
	GatewayEventGuildMemberRemove        GatewayEventType = "GUILD_MEMBER_REMOVE"
	GatewayEventGuildMemberUpdate        GatewayEventType = "GUILD_MEMBER_UPDATE"
	GatewayEventGuildMembersChunk        GatewayEventType = "GUILD_MEMBERS_CHUNK"
	GatewayEventGuildRoleCreate          GatewayEventType = "GUILD_ROLE_CREATE"
	GatewayEventGuildRoleDelete          GatewayEventType = "GUILD_ROLE_DELETE"
	GatewayEventGuildRoleUpdate          GatewayEventType = "GUILD_ROLE_UPDATE"
	GatewayEventGuildUpdate              GatewayEventType = "GUILD_UPDATE"
	GatewayEventInteractionCreate        GatewayEventType = "INTERACTION_CREATE"
	WebhookEventInteractionCreate        GatewayEventType = "INTERACTION_WEBHOOK_CREATE"
	GatewayEventMessageAck               GatewayEventType = "MESSAGE_ACK"
	GatewayEventMessageCreate            GatewayEventType = "MESSAGE_CREATE"
	GatewayEventMessageDelete            GatewayEventType = "MESSAGE_DELETE"
	GatewayEventMessageDeleteBulk        GatewayEventType = "MESSAGE_DELETE_BULK"
	GatewayEventMessageReactionAdd       GatewayEventType = "MESSAGE_REACTION_ADD"
	GatewayEventMessageReactionRemove    GatewayEventType = "MESSAGE_REACTION_REMOVE"
	GatewayEventMessageReactionRemoveAll GatewayEventType = "MESSAGE_REACTION_REMOVE_ALL"
	GatewayEventMessageUpdate            GatewayEventType = "MESSAGE_UPDATE"
	GatewayEventPresenceUpdate           GatewayEventType = "PRESENCE_UPDATE"
	GatewayEventPresencesReplace         GatewayEventType = "PRESENCES_REPLACE"
	GatewayEventReady                    GatewayEventType = "READY"
	GatewayEventResumed                  GatewayEventType = "RESUMED"
	GatewayEventTypingStart              GatewayEventType = "TYPING_START"
	GatewayEventUserGuildSettingsUpdate  GatewayEventType = "USER_GUILD_SETTINGS_UPDATE"
	GatewayEventUserNoteUpdate           GatewayEventType = "USER_NOTE_UPDATE"
	GatewayEventUserSettingsUpdate       GatewayEventType = "USER_SETTINGS_UPDATE"
	GatewayEventUserUpdate               GatewayEventType = "USER_UPDATE"
	GatewayEventVoiceServerUpdate        GatewayEventType = "VOICE_SERVER_UPDATE"
	GatewayEventVoiceStateUpdate         GatewayEventType = "VOICE_STATE_UPDATE"
	GatewayEventWebhooksUpdate           GatewayEventType = "WEBHOOKS_UPDATE"
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
