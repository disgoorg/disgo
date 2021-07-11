package api

// AuditLogChangeKey (https://discord.com/developers/docs/resources/audit-log#audit-log-change-object-audit-log-change-key) is data representing changes values/settings in an audit log.
type AuditLogChangeKey struct {
	Name                        *string                     `json:"name"`
	Description                 *string                     `json:"description"`
	IconHash                    *string                     `json:"icon_hash"`
	SplashHash                  *string                     `json:"splash_hash"`
	DiscoverySplashHash         *string                     `json:"discovery_splash_hash"`
	BannerHash                  *string                     `json:"banner_hash"`
	OwnerID                     *Snowflake                  `json:"owner_id"`
	Region                      *string                     `json:"region"`
	PreferredLocale             *string                     `json:"preferred_locale"`
	AFKChannelID                *Snowflake                  `json:"afk_channel_id"`
	AFKTimeout                  *int                        `json:"afk_timeout"`
	RulesChannelID              *Snowflake                  `json:"rules_channel_id"`
	PublicUpdatesChannelID      *Snowflake                  `json:"public_updates_channel_id"`
	MFALevel                    *MFALevel                   `json:"mfa_level"`
	VerificationLevel           *VerificationLevel          `json:"verification_level"`
	ExplicitContentFilterLevel  *ExplicitContentFilterLevel `json:"explicit_content_filter"`
	DefaultMessageNotifications *MessageNotifications       `json:"default_message_notifications"`
	VanityURLCode               *string                     `json:"vanity_url_code"`
	Add                         []Role                      `json:"$add"`
	Remove                      []Role                      `json:"$remove"`
	PruneDeleteDays             *int                        `json:"prune_delete_days"`
	WidgetEnabled               *bool                       `json:"widget_enabled"`
	WidgetChannelID             *string                     `json:"widget_channel_id"`
	SystemChannelID             *string                     `json:"system_channel_id"`
	Position                    *int                        `json:"position"`
	Topic                       *string                     `json:"topic"`
	Bitrate                     *int                        `json:"bitrate"`
	PermissionOverwrites        []PermissionOverwrite       `json:"permission_overwrites"`
	NSFW                        *bool                       `json:"nsfw"`
	ApplicationID               *Snowflake                  `json:"application_id"`
	RateLimitPerUser            *int                        `json:"ratelimit_per_user"`
	Permissions                 *string                     `json:"permissions"`
	Color                       *int                        `json:"color"`
	Hoist                       *bool                       `json:"hoist"`
	Mentionable                 *bool                       `json:"mentionable"`
	Allow                       *Permissions                `json:"allow"`
	Deny                        *Permissions                `json:"deny"`
	Code                        *string                     `json:"code"`
	ChannelID                   *Snowflake                  `json:"channel_id"`
	InviterID                   *Snowflake                  `json:"inviter_id"`
	MaxUses                     *int                        `json:"max_uses"`
	Uses                        *int                        `json:"uses"`
	MaxAge                      *string                     `json:"max_age"`
	Temporary                   *bool                       `json:"temporary"`
	Deaf                        *bool                       `json:"deaf"`
	Mute                        *bool                       `json:"mute"`
	Nick                        *string                     `json:"nick"`
	AvatarHash                  *string                     `json:"avatar_hash"`
	ID                          *Snowflake                  `json:"id"`
	Type                        interface{}                 `json:"type"`
	EnableEmoticons             *bool                       `json:"enable_emoticons"`
	ExpireBehavior              *int                        `json:"expire_behavior"`
	ExpireGracePeriod           *int                        `json:"expire_grace_period"`
	UserLimit                   *int                        `json:"user_limit"`
	PrivacyLevel                *int                        `json:"privacy_level"`
}

// AuditLogEvent is an 8-bit unsigned integer representing an audit log event.
type AuditLogEvent int

// AuditLogEventGuildUpdate
const (
	AuditLogEventGuildUpdate AuditLogEvent = 1
)

// AuditLogEventChannelCreate
const (
	AuditLogEventChannelCreate AuditLogEvent = iota + 10
	AuditLogEventChannelUpdate
	AuditLogEventChannelDelete
	AuditLogEventChannelOverwriteCreate
	AuditLogEventChannelOverwriteUpdate
	AuditLogEventChannelOverwriteDelete
)

// AuditLogEventMemberKick
const (
	AuditLogEventMemberKick AuditLogEvent = iota + 20
	AuditLogEventMemberPrune
	AuditLogEventMemberBanAdd
	AuditLogEventMemberBanRemove
	AuditLogEventMemberUpdate
	AuditLogEventMemberRoleUpdate
	AuditLogEventMemberMove
	AuditLogEventMemberDisconnect
	AuditLogEventBotAdd
)

// AuditLogEventRoleCreate
const (
	AuditLogEventRoleCreate AuditLogEvent = iota + 30
	AuditLogEventRoleUpdate
	AuditLogEventRoleDelete
)

// AuditLogEventInviteCreate
const (
	AuditLogEventInviteCreate AuditLogEvent = iota + 40
	AuditLogEventInviteUpdate
	AuditLogEventInviteDelete
)

// AuditLogEventWebhookCreate
const (
	AuditLogEventWebhookCreate AuditLogEvent = iota + 50
	AuditLogEventWebhookUpdate
	AuditLogEventWebhookDelete
)

// AuditLogEventEmojiCreate
const (
	AuditLogEventEmojiCreate AuditLogEvent = iota + 60
	AuditLogEventEmojiUpdate
	AuditLogEventEmojiDelete
)

// AuditLogEventMessageDelete
const (
	AuditLogEventMessageDelete AuditLogEvent = iota + 72
	AuditLogEventMessageBulkDelete
	AuditLogEventMessagePin
	AuditLogEventMessageUnpin
)

// AuditLogEventIntegrationCreate
const (
	AuditLogEventIntegrationCreate AuditLogEvent = iota + 80
	AuditLogEventIntegrationUpdate
	AuditLogEventIntegrationDelete
	AuditLogEventStageInstanceCreate
	AuditLogEventStageInstanceUpdate
	AuditLogEventStageInstanceDelete
)

// OptionalAuditLogEntryInfo (https://discord.com/developers/docs/resources/audit-log#audit-log-entry-object-optional-audit-entry-info)
type OptionalAuditLogEntryInfo struct {
	DeleteMemberDays *string    `json:"delete_member_days"`
	MembersRemoved   *string    `json:"members_removed"`
	ChannelID        *Snowflake `json:"channel_id"`
	MessageID        *Snowflake `json:"message_id"`
	Count            *string    `json:"count"`
	ID               *string    `json:"id"`
	Type             *string    `json:"type"`
	RoleName         *string    `json:"role_name"`
}

// AuditLogEntry (https://discord.com/developers/docs/resources/audit-log#audit-log-entry-object)
type AuditLogEntry struct {
	TargetID   *Snowflake                 `json:"target_id"`
	Changes    []AuditLogChangeKey        `json:"changes"`
	UserID     Snowflake                  `json:"user_id"`
	ID         Snowflake                  `json:"id"`
	ActionType AuditLogEvent              `json:"action_type"`
	Options    *OptionalAuditLogEntryInfo `json:"options"`
	Reason     *string                    `json:"reason"`
}

// AuditLog (https://discord.com/developers/docs/resources/audit-log) These are logs of events that occurred, accessible via the Discord API.
type AuditLog struct {
	Webhooks     []Webhook
	Users        []User
	Entries      []AuditLogEntry
	Integrations []Integration
}
