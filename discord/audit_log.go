package discord

import (
	"github.com/disgoorg/json"
	"github.com/disgoorg/snowflake/v2"
)

// AuditLogEvent is an 8-bit unsigned integer representing an audit log event.
type AuditLogEvent int

const (
	AuditLogEventGuildUpdate AuditLogEvent = 1
)

const (
	AuditLogEventChannelCreate AuditLogEvent = iota + 10
	AuditLogEventChannelUpdate
	AuditLogEventChannelDelete
	AuditLogEventChannelOverwriteCreate
	AuditLogEventChannelOverwriteUpdate
	AuditLogEventChannelOverwriteDelete
)

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

const (
	AuditLogEventRoleCreate AuditLogEvent = iota + 30
	AuditLogEventRoleUpdate
	AuditLogEventRoleDelete
)

const (
	AuditLogEventInviteCreate AuditLogEvent = iota + 40
	AuditLogEventInviteUpdate
	AuditLogEventInviteDelete
)

const (
	AuditLogEventWebhookCreate AuditLogEvent = iota + 50
	AuditLogEventWebhookUpdate
	AuditLogEventWebhookDelete
)

const (
	AuditLogEventEmojiCreate AuditLogEvent = iota + 60
	AuditLogEventEmojiUpdate
	AuditLogEventEmojiDelete
)

const (
	AuditLogEventMessageDelete AuditLogEvent = iota + 72
	AuditLogEventMessageBulkDelete
	AuditLogEventMessagePin
	AuditLogEventMessageUnpin
)

const (
	AuditLogEventIntegrationCreate AuditLogEvent = iota + 80
	AuditLogEventIntegrationUpdate
	AuditLogEventIntegrationDelete
	AuditLogEventStageInstanceCreate
	AuditLogEventStageInstanceUpdate
	AuditLogEventStageInstanceDelete
)

const (
	AuditLogEventStickerCreate AuditLogEvent = iota + 90
	AuditLogEventStickerUpdate
	AuditLogEventStickerDelete
)

const (
	AuditLogGuildScheduledEventCreate AuditLogEvent = iota + 100
	AuditLogGuildScheduledEventUpdate
	AuditLogGuildScheduledEventDelete
)

const (
	AuditLogThreadCreate AuditLogEvent = iota + 110
	AuditLogThreadUpdate
	AuditLogThreadDelete
)

const (
	AuditLogApplicationCommandPermissionUpdate AuditLogEvent = 121
)

const (
	AuditLogSoundboardSoundCreate AuditLogEvent = iota + 130
	AuditLogSoundboardSoundUpdate
	AuditLogSoundboardSoundDelete
)

const (
	AuditLogAutoModerationRuleCreate AuditLogEvent = iota + 140
	AuditLogAutoModerationRuleUpdate
	AuditLogAutoModerationRuleDelete
	AuditLogAutoModerationBlockMessage
	AuditLogAutoModerationFlagToChannel
	AuditLogAutoModerationUserCommunicationDisabled
)

const (
	AuditLogCreatorMonetizationRequestCreated AuditLogEvent = iota + 150
	AuditLogCreatorMonetizationTermsAccepted
)

const (
	AuditLogOnboardingPromptCreate AuditLogEvent = iota + 163
	AuditLogOnboardingPromptUpdate
	AuditLogOnboardingPromptDelete
	AuditLogOnboardingCreate
	AuditLogOnboardingUpdate
)

const (
	AuditLogHomeSettingsCreate AuditLogEvent = iota + 190
	AuditLogHomeSettingsUpdate
)

// AuditLogChangeKey is a string representing a key in the audit log change object.
type AuditLogChangeKey string

const (
	AuditLogChangeKeyAFKChannelID AuditLogChangeKey = "afk_channel_id"
	AuditLogChangeKeyAFKTimeout   AuditLogChangeKey = "afk_timeout"
	// AuditLogChangeKeyAllow is sent when a role's permission overwrites changed (stringy int)
	AuditLogChangeKeyAllow         AuditLogChangeKey = "allow"
	AuditLogChangeKeyApplicationID AuditLogChangeKey = "application_id"
	// AuditLogChangeKeyArchived is sent when a channel thread is archived/unarchived (bool)
	AuditLogChangeKeyArchived AuditLogChangeKey = "archived"
	AuditLogChangeKeyAsset    AuditLogChangeKey = "asset"
	// AuditLogChangeKeyAutoArchiveDuration is sent when a thread's auto archive duration is changed (int)
	AuditLogChangeKeyAutoArchiveDuration AuditLogChangeKey = "auto_archive_duration"
	AuditLogChangeKeyAvailable           AuditLogChangeKey = "available"
	AuditLogChangeKeyAvatarHash          AuditLogChangeKey = "avatar_hash"
	AuditLogChangeKeyBannerHash          AuditLogChangeKey = "banner_hash"
	AuditLogChangeKeyBitrate             AuditLogChangeKey = "bitrate"
	AuditLogChangeKeyChannelID           AuditLogChangeKey = "channel_id"
	AuditLogChangeKeyCode                AuditLogChangeKey = "code"
	// AuditLogChangeKeyColor is sent when a role's color is changed (int)
	AuditLogChangeKeyColor AuditLogChangeKey = "color"
	// AuditLogChangeKeyCommunicationDisabledUntil is sent when a user's communication disabled until datetime is changed (stringy ISO8601 datetime)
	AuditLogChangeKeyCommunicationDisabledUntil AuditLogChangeKey = "communication_disabled_until"
	// AuditLogChangeKeyDeaf is sent when a user is set to be server deafened/undeafened (bool)
	AuditLogChangeKeyDeaf                        AuditLogChangeKey = "deaf"
	AuditLogChangeKeyDefaultAutoArchiveDuration  AuditLogChangeKey = "default_auto_archive_duration"
	AuditLogChangeKeyDefaultMessageNotifications AuditLogChangeKey = "default_message_notifications"
	// AuditLogChangeKeyDeny is sent when a role's permission overwrites changed (stringed int)
	AuditLogChangeKeyDeny                  AuditLogChangeKey = "deny"
	AuditLogChangeKeyDescription           AuditLogChangeKey = "description"
	AuditLogChangeKeyDiscoverySplashHash   AuditLogChangeKey = "discovery_splash_hash"
	AuditLogChangeKeyEnableEmoticons       AuditLogChangeKey = "enable_emoticons"
	AuditLogChangeKeyEntityType            AuditLogChangeKey = "entity_type"
	AuditLogChangeKeyExpireBehavior        AuditLogChangeKey = "expire_behavior"
	AuditLogChangeKeyExpireGracePeriod     AuditLogChangeKey = "expire_grace_period"
	AuditLogChangeKeyExplicitContentFilter AuditLogChangeKey = "explicit_content_filter"
	AuditLogChangeKeyFormatType            AuditLogChangeKey = "format_type"
	AuditLogChangeKeyGuildID               AuditLogChangeKey = "guild_id"
	// AuditLogChangeKeyHoist is sent when a role is set to be displayed separately from online members (bool)
	AuditLogChangeKeyHoist     AuditLogChangeKey = "hoist"
	AuditLogChangeKeyIconHash  AuditLogChangeKey = "icon_hash"
	AuditLogChangeKeyID        AuditLogChangeKey = "id"
	AuditLogChangeKeyInvitable AuditLogChangeKey = "invitable"
	AuditLogChangeKeyInviterID AuditLogChangeKey = "inviter_id"
	AuditLogChangeKeyLocation  AuditLogChangeKey = "location"
	// AuditLogChangeKeyLocked is sent when a channel thread is locked/unlocked (bool)
	AuditLogChangeKeyLocked  AuditLogChangeKey = "locked"
	AuditLogChangeKeyMaxAge  AuditLogChangeKey = "max_age"
	AuditLogChangeKeyMaxUses AuditLogChangeKey = "max_uses"
	// AuditLogChangeKeyMentionable is sent when a role changes its mentionable state (bool)
	AuditLogChangeKeyMentionable AuditLogChangeKey = "mentionable"
	AuditLogChangeKeyMFALevel    AuditLogChangeKey = "mfa_level"
	// AuditLogChangeKeyMute is sent when a user is server muted/unmuted (bool)
	AuditLogChangeKeyMute AuditLogChangeKey = "mute"
	AuditLogChangeKeyName AuditLogChangeKey = "name"
	// AuditLogChangeKeyNick is sent when a user's nickname is changed (string)
	AuditLogChangeKeyNick AuditLogChangeKey = "nick"
	AuditLogChangeKeyNSFW AuditLogChangeKey = "nsfw"
	// AuditLogChangeKeyOwnerID is sent when owner id of a guild changed (snowflake.ID)
	AuditLogChangeKeyOwnerID AuditLogChangeKey = "owner_id"
	// AuditLogChangeKeyPermissionOverwrites is sent when a role's permission overwrites changed (string)
	AuditLogChangeKeyPermissionOverwrites AuditLogChangeKey = "permission_overwrites"
	// AuditLogChangeKeyPermissions is sent when a role's permissions changed (string)
	AuditLogChangeKeyPermissions AuditLogChangeKey = "permissions"
	// AuditLogChangeKeyPosition is sent when channel position changed (int)
	AuditLogChangeKeyPosition               AuditLogChangeKey = "position"
	AuditLogChangeKeyPreferredLocale        AuditLogChangeKey = "preferred_locale"
	AuditLogChangeKeyPrivacyLevel           AuditLogChangeKey = "privacy_level"
	AuditLogChangeKeyPruneDeleteDays        AuditLogChangeKey = "prune_delete_days"
	AuditLogChangeKeyPublicUpdatesChannelID AuditLogChangeKey = "public_updates_channel_id"
	AuditLogChangeKeyRateLimitPerUser       AuditLogChangeKey = "rate_limit_per_user"
	AuditLogChangeKeyRegion                 AuditLogChangeKey = "region"
	AuditLogChangeKeyRulesChannelID         AuditLogChangeKey = "rules_channel_id"
	AuditLogChangeKeySplashHash             AuditLogChangeKey = "splash_hash"
	AuditLogChangeKeyStatus                 AuditLogChangeKey = "status"
	// AuditLogChangeKeySystemChannelID is sent when system channel id of a guild changed (snowflake.ID)
	AuditLogChangeKeySystemChannelID AuditLogChangeKey = "system_channel_id"
	AuditLogChangeKeyTags            AuditLogChangeKey = "tags"
	AuditLogChangeKeyTemporary       AuditLogChangeKey = "temporary"
	// AuditLogChangeKeyTopic is sent when channel topic changed (string)
	AuditLogChangeKeyTopic        AuditLogChangeKey = "topic"
	AuditLogChangeKeyType         AuditLogChangeKey = "type"
	AuditLogChangeKeyUnicodeEmoji AuditLogChangeKey = "unicode_emoji"
	// AuditLogChangeKeyUserLimit is sent when user limit of a voice channel changed (int)
	AuditLogChangeKeyUserLimit     AuditLogChangeKey = "user_limit"
	AuditLogChangeKeyUses          AuditLogChangeKey = "uses"
	AuditLogChangeKeyVanityURLCode AuditLogChangeKey = "vanity_url_code"
	// AuditLogChangeKeyVerificationLevel is sent when verification level of the server changed (int)
	AuditLogChangeKeyVerificationLevel AuditLogChangeKey = "verification_level"
	AuditLogChangeKeyWidgetChannelID   AuditLogChangeKey = "widget_channel_id"
	// AuditLogChangeKeyWidgetEnabled is sent when a server widget is enabled/disabled (bool)
	AuditLogChangeKeyWidgetEnabled AuditLogChangeKey = "widget_enabled"
	// AuditLogChangeKeyRoleAdd is sent when roles are added to a user (array of discord.PartialRole JSON)
	AuditLogChangeKeyRoleAdd AuditLogChangeKey = "$add"
	// AuditLogChangeKeyRoleRemove is sent when roles are removed from a user (array of discord.PartialRole JSON)
	AuditLogChangeKeyRoleRemove AuditLogChangeKey = "$remove"
)

// AuditLog (https://discord.com/developers/docs/resources/audit-log) These are logs of events that occurred, accessible via the Discord
type AuditLog struct {
	ApplicationCommands  []ApplicationCommand  `json:"application_commands"`
	AuditLogEntries      []AuditLogEntry       `json:"audit_log_entries"`
	AutoModerationRules  []AutoModerationRule  `json:"auto_moderation_rules"`
	GuildScheduledEvents []GuildScheduledEvent `json:"guild_scheduled_events"`
	Integrations         []Integration         `json:"integrations"`
	Threads              []GuildThread         `json:"threads"`
	Users                []User                `json:"users"`
	Webhooks             []Webhook             `json:"webhooks"`
}

func (l *AuditLog) UnmarshalJSON(data []byte) error {
	type auditLog AuditLog
	var v struct {
		ApplicationCommands []UnmarshalApplicationCommand `json:"application_commands"`
		Integrations        []UnmarshalIntegration        `json:"integrations"`
		Webhooks            []UnmarshalWebhook            `json:"webhooks"`
		auditLog
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*l = AuditLog(v.auditLog)

	if v.ApplicationCommands != nil {
		l.ApplicationCommands = make([]ApplicationCommand, len(v.ApplicationCommands))
		for i := range v.ApplicationCommands {
			l.ApplicationCommands[i] = v.ApplicationCommands[i].ApplicationCommand
		}
	}

	if v.Integrations != nil {
		l.Integrations = make([]Integration, len(v.Integrations))
		for i := range v.Integrations {
			l.Integrations[i] = v.Integrations[i].Integration
		}
	}

	if v.Webhooks != nil {
		l.Webhooks = make([]Webhook, len(v.Webhooks))
		for i := range v.Webhooks {
			l.Webhooks[i] = v.Webhooks[i].Webhook
		}
	}

	return nil
}

// AuditLogEntry (https://discord.com/developers/docs/resources/audit-log#audit-log-entry-object)
type AuditLogEntry struct {
	TargetID   *snowflake.ID              `json:"target_id"`
	Changes    []AuditLogChange           `json:"changes"`
	UserID     snowflake.ID               `json:"user_id"`
	ID         snowflake.ID               `json:"id"`
	ActionType AuditLogEvent              `json:"action_type"`
	Options    *OptionalAuditLogEntryInfo `json:"options"`
	Reason     *string                    `json:"reason"`
}

// AuditLogChange (https://discord.com/developers/docs/resources/audit-log#audit-log-change-object) contains what was changed.
// For a list of possible keys & values see the discord documentation.
type AuditLogChange struct {
	// NewValue is the new value of the key after the change as a json.RawMessage.
	NewValue json.RawMessage `json:"new_value"`
	// OldValue is the old value of the key before the change as a json.RawMessage.
	OldValue json.RawMessage `json:"old_value"`
	// Key is the key of the change.
	Key AuditLogChangeKey `json:"key"`
}

// UnmarshalNewValue unmarshals the NewValue field into the provided type.
func (c *AuditLogChange) UnmarshalNewValue(v any) error {
	return json.Unmarshal(c.NewValue, v)
}

// UnmarshalOldValue unmarshals the OldValue field into the provided type.
func (c *AuditLogChange) UnmarshalOldValue(v any) error {
	return json.Unmarshal(c.OldValue, v)
}

// OptionalAuditLogEntryInfo (https://discord.com/developers/docs/resources/audit-log#audit-log-entry-object-optional-audit-entry-info)
type OptionalAuditLogEntryInfo struct {
	DeleteMemberDays              *string                    `json:"delete_member_days"`
	MembersRemoved                *string                    `json:"members_removed"`
	ChannelID                     *snowflake.ID              `json:"channel_id"`
	MessageID                     *snowflake.ID              `json:"message_id"`
	Count                         *string                    `json:"count"`
	ID                            *string                    `json:"id"`
	Type                          *string                    `json:"type"`
	RoleName                      *string                    `json:"role_name"`
	ApplicationID                 *snowflake.ID              `json:"application_id"`
	AutoModerationRuleName        *string                    `json:"auto_moderation_rule_name"`
	AutoModerationRuleTriggerType *AutoModerationTriggerType `json:"auto_moderation_rule_trigger_type,string"`
	IntegrationType               *IntegrationType           `json:"integration_type"`
}
