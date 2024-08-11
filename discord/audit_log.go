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
	Key string `json:"key"`
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
