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
	AFKTimeout                  *uint8                      `json:"afk_timeout"`
	RulesChannelID              *Snowflake                  `json:"rules_channel_id"`
	PublicUpdatesChannelID      *Snowflake                  `json:"public_updates_channel_id"`
	MFALevel                    *MFALevel                   `json:"mfa_level"`
	VerificationLevel           *VerificationLevel          `json:"verification_level"`
	ExplicitContentFilterLevel  *ExplicitContentFilterLevel `json:"explicit_content_filter"`
	DefaultMessageNotifications *MessageNotifications       `json:"default_message_notifications"`
	VanityURLCode               *string                     `json:"vanity_url_code"`
	Add                         *[]Role                     `json:"$add"`
	Remove                      *[]Role                     `json:"$remove"`
	PruneDeleteDays             *uint8                      `json:"prune_delete_days"`
	WidgetEnabled               *bool                       `json:"widget_enabled"`
	WidgetChannelID             *string                     `json:"widget_channel_id"`
	SystemChannelID             *string                     `json:"system_channel_id"`
	Position                    *uint8                      `json:"position"`
	Topic                       *string                     `json:"topic"`
	Bitrate                     *uint8                      `json:"bitrate"`
	PermissionOverwrites        *[]PermissionOverwrite      `json:"permission_overwrites"`
	NSFW                        *bool                       `json:"nsfw"`
	ApplicationID               *Snowflake                  `json:"application_id"`
	RateLimitPerUser            *uint                       `json:"ratelimit_per_user"`
	Permissions                 *string                     `json:"permissions"`
	Color                       *uint                       `json:"color"`
	Hoist                       *bool                       `json:"hoist"`
	Mentionable                 *bool                       `json:"mentionable"`
	Allow                       *string                     `json:"allow"`
	Deny                        *string                     `json:"deny"`
	Code                        *string                     `json:"code"`
	ChannelID                   *Snowflake                  `json:"channel_id"`
	InviterID                   *Snowflake                  `json:"inviter_id"`
	MaxUses                     *uint8                      `json:"max_uses"`
	Uses                        *uint                       `json:"uses"`
	MaxAge                      *string                     `json:"max_age"`
	Temporary                   *bool                       `json:"temporary"`
	Deaf                        *bool                       `json:"deaf"`
	Mute                        *bool                       `json:"mute"`
	Nick                        *string                     `json:"nick"`
	AvatarHash                  *string                     `json:"avatar_hash"`
	ID                          *Snowflake                  `json:"id"`
	Type                        *interface{}                `json:"type"`
	EnableEmoticons             *bool                       `json:"enable_emoticons"`
	ExpireBehavior              *uint                       `json:"expire_behavior"`
	ExpireGracePeriod           *uint                       `json:"expire_grace_period"`
	UserLimit                   *uint8                      `json:"user_limit"`
	PrivacyLevel                *uint8                      `json:"privacy_level"`
}

// AuditLogEvent is an 8-bit unsigned integer representing an audit log event.
type AuditLogEvent uint8

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
	Changes    *[]AuditLogChangeKey       `json:"changes"`
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
