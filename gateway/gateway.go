package gateway

import (
	"io"
	"time"

	"github.com/DisgoOrg/log"
)

// Status is the state that the client is currently in
type Status int

// IsConnected returns whether you can send payloads to the Gateway
func (s Status) IsConnected() bool {
	switch s {
	case StatusWaitingForGuilds, StatusReady:
		return true
	default:
		return false
	}
}

// Indicates how far along the client is to connecting
const (
	StatusUnconnected Status = iota
	StatusConnecting
	StatusReconnecting
	StatusIdentifying
	StatusWaitingForHello
	StatusWaitingForReady
	StatusWaitingForGuilds
	StatusReady
	StatusDisconnected
	StatusResuming
)

type Config struct {
	LargeThreshold int
	GatewayIntents Intents
	OS             string
	Browser        string
	Device         string
}

type EventHandlerFunc func(eventType EventType, sequenceNumber int, payload io.Reader)

// Gateway is what is used to connect to discord
type Gateway interface {
	Logger() log.Logger
	Config() Config
	Open() error
	Close()
	Status() Status
	Send(command GatewayCommand) error
	Latency() time.Duration
}

// Op are opcodes used by discord
type Op int

// Constants for the gateway opcodes
//goland:noinspection GoUnusedConst
const (
	OpDispatch Op = iota
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

// EventType wraps all EventType types
type EventType string

// Constants for the gateway events
//goland:noinspection GoUnusedConst
const (
	EventTypeHello                         EventType = "HELLO"
	EventTypeReady                         EventType = "READY"
	EventTypeResumed                       EventType = "RESUMED"
	EventTypeReconnect                     EventType = "RECONNECT"
	EventTypeInvalidSession                EventType = "INVALID_SESSION"
	EventTypeCommandCreate                 EventType = "APPLICATION_COMMAND_CREATE"
	EventTypeCommandUpdate                 EventType = "APPLICATION_COMMAND_UPDATE"
	EventTypeCommandDelete                 EventType = "APPLICATION_COMMAND_DELETE"
	EventTypeChannelCreate                 EventType = "CHANNEL_CREATE"
	EventTypeChannelUpdate                 EventType = "CHANNEL_UPDATE"
	EventTypeChannelDelete                 EventType = "CHANNEL_DELETE"
	EventTypeChannelPinsUpdate             EventType = "CHANNEL_PINS_UPDATE"
	EventTypeGuildCreate                   EventType = "GUILD_CREATE"
	EventTypeGuildUpdate                   EventType = "GUILD_UPDATE"
	EventTypeGuildDelete                   EventType = "GUILD_DELETE"
	EventTypeGuildBanAdd                   EventType = "GUILD_BAN_ADD"
	EventTypeGuildBanRemove                EventType = "GUILD_BAN_REMOVE"
	EventTypeGuildEmojisUpdate             EventType = "GUILD_EMOJIS_UPDATE"
	EventTypeGuildIntegrationsUpdate       EventType = "GUILD_INTEGRATIONS_UPDATE"
	EventTypeGuildMemberAdd                EventType = "GUILD_MEMBER_ADD"
	EventTypeGuildMemberRemove             EventType = "GUILD_MEMBER_REMOVE"
	EventTypeGuildMemberUpdate             EventType = "GUILD_MEMBER_UPDATE"
	EventTypeGuildMembersChunk             EventType = "GUILD_MEMBERS_CHUNK"
	EventTypeGuildRoleCreate               EventType = "GUILD_ROLE_CREATE"
	EventTypeGuildRoleUpdate               EventType = "GUILD_ROLE_UPDATE"
	EventTypeGuildRoleDelete               EventType = "GUILD_ROLE_DELETE"
	EventTypeIntegrationCreate             EventType = "INTEGRATION_CREATE"
	EventTypeIntegrationUpdate             EventType = "INTEGRATION_UPDATE"
	EventTypeIntegrationDelete             EventType = "INTEGRATION_DELETE"
	EventTypeInteractionCreate             EventType = "INTERACTION_CREATE"
	EventTypeInviteCreate                  EventType = "INVITE_CREATE"
	EventTypeInviteDelete                  EventType = "INVITE_DELETE"
	EventTypeMessageCreate                 EventType = "MESSAGE_CREATE"
	EventTypeMessageUpdate                 EventType = "MESSAGE_UPDATE"
	EventTypeMessageDelete                 EventType = "MESSAGE_DELETE"
	EventTypeMessageDeleteBulk             EventType = "MESSAGE_DELETE_BULK"
	EventTypeMessageReactionAdd            EventType = "MESSAGE_REACTION_ADD"
	EventTypeMessageReactionRemove         EventType = "MESSAGE_REACTION_REMOVE"
	EventTypeMessageReactionRemoveAll      EventType = "MESSAGE_REACTION_REMOVE_ALL"
	EventTypeMessageReactionRemoveAllEmoji EventType = "MESSAGE_REACTION_REMOVE_ALL_EMOJI"
	EventTypePresenceUpdate                EventType = "PRESENCE_UPDATE"
	EventTypeTypingStart                   EventType = "TYPING_START"
	EventTypeUserUpdate                    EventType = "USER_UPDATE"
	EventTypeVoiceStateUpdate              EventType = "VOICE_STATE_UPDATE"
	EventTypeVoiceServerUpdate             EventType = "VOICE_SERVER_UPDATE"
	EventTypeWebhooksUpdate                EventType = "WEBHOOKS_UPDATE"
)
