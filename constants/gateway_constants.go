package constants

// Gateway opcodes used by discord
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

const (
	ChannelCreateEvent            = "CHANNEL_CREATE"
	ChannelDeleteEvent            = "CHANNEL_DELETE"
	ChannelPinsUpdateEvent        = "CHANNEL_PINS_UPDATE"
	ChannelUpdateEvent            = "CHANNEL_UPDATE"
	GuildBanAddEvent              = "GUILD_BAN_ADD"
	GuildBanRemoveEvent           = "GUILD_BAN_REMOVE"
	GuildCreateEvent              = "GUILD_CREATE"
	GuildDeleteEvent              = "GUILD_DELETE"
	GuildEmojisUpdateEvent        = "GUILD_EMOJIS_UPDATE"
	GuildIntegrationsUpdateEvent  = "GUILD_INTEGRATIONS_UPDATE"
	GuildMemberAddEvent           = "GUILD_MEMBER_ADD"
	GuildMemberRemoveEvent        = "GUILD_MEMBER_REMOVE"
	GuildMemberUpdateEvent        = "GUILD_MEMBER_UPDATE"
	GuildMembersChunkEvent        = "GUILD_MEMBERS_CHUNK"
	GuildRoleCreateEvent          = "GUILD_ROLE_CREATE"
	GuildRoleDeleteEvent          = "GUILD_ROLE_DELETE"
	GuildRoleUpdateEvent          = "GUILD_ROLE_UPDATE"
	GuildUpdateEvent              = "GUILD_UPDATE"
	InteractionCreateEvent        = "INTERACTION_CREATE"
	MessageAckEvent               = "MESSAGE_ACK"
	MessageCreateEvent            = "MESSAGE_CREATE"
	MessageDeleteEvent            = "MESSAGE_DELETE"
	MessageDeleteBulkEvent        = "MESSAGE_DELETE_BULK"
	MessageReactionAddEvent       = "MESSAGE_REACTION_ADD"
	MessageReactionRemoveEvent    = "MESSAGE_REACTION_REMOVE"
	MessageReactionRemoveAllEvent = "MESSAGE_REACTION_REMOVE_ALL"
	MessageUpdateEvent            = "MESSAGE_UPDATE"
	PresenceUpdateEvent           = "PRESENCE_UPDATE"
	PresencesReplaceEvent         = "PRESENCES_REPLACE"
	ReadyEvent                    = "READY"
	ResumedEvent                  = "RESUMED"
	TypingStartEvent              = "TYPING_START"
	UserGuildSettingsUpdateEvent  = "USER_GUILD_SETTINGS_UPDATE"
	UserNoteUpdateEvent           = "USER_NOTE_UPDATE"
	UserSettingsUpdateEvent       = "USER_SETTINGS_UPDATE"
	UserUpdateEvent               = "USER_UPDATE"
	VoiceServerUpdateEvent        = "VOICE_SERVER_UPDATE"
	VoiceStateUpdateEvent         = "VOICE_STATE_UPDATE"
	WebhooksUpdateEvent           = "WEBHOOKS_UPDATE"
)
