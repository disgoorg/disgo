package gateway

import (
	"io"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/json"
	"github.com/disgoorg/snowflake/v2"
)

type EventData interface {
	MessageData
	eventData()
}

// EventReady is the event sent by discord when you successfully Identify
type EventReady struct {
	Version          int                        `json:"v"`
	User             discord.OAuth2User         `json:"user"`
	Guilds           []discord.UnavailableGuild `json:"guilds"`
	SessionID        string                     `json:"session_id"`
	ResumeGatewayURL string                     `json:"resume_gateway_url"`
	Shard            []int                      `json:"shard,omitempty"`
	Application      discord.PartialApplication `json:"application"`
}

func (EventReady) messageData() {}
func (EventReady) eventData()   {}

type EventApplicationCommandPermissionsUpdate struct {
	discord.ApplicationCommandPermissions
}

func (EventApplicationCommandPermissionsUpdate) messageData() {}
func (EventApplicationCommandPermissionsUpdate) eventData()   {}

type EventChannelCreate struct {
	discord.Channel
}

func (e *EventChannelCreate) UnmarshalJSON(data []byte) error {
	var v discord.UnmarshalChannel
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	e.Channel = v.Channel
	return nil
}

func (EventChannelCreate) messageData() {}
func (EventChannelCreate) eventData()   {}

type EventChannelUpdate struct {
	discord.Channel
}

func (e *EventChannelUpdate) UnmarshalJSON(data []byte) error {
	var v discord.UnmarshalChannel
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	e.Channel = v.Channel
	return nil
}

func (EventChannelUpdate) messageData() {}
func (EventChannelUpdate) eventData()   {}

type EventChannelDelete struct {
	discord.Channel
}

func (e *EventChannelDelete) UnmarshalJSON(data []byte) error {
	var v discord.UnmarshalChannel
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	e.Channel = v.Channel
	return nil
}

func (EventChannelDelete) messageData() {}
func (EventChannelDelete) eventData()   {}

type EventThreadCreate struct {
	discord.GuildThread
	ThreadMember discord.ThreadMember `json:"thread_member"`
}

func (EventThreadCreate) messageData() {}
func (EventThreadCreate) eventData()   {}

type EventThreadUpdate struct {
	discord.GuildThread
}

func (EventThreadUpdate) messageData() {}
func (EventThreadUpdate) eventData()   {}

type EventThreadDelete struct {
	ID       snowflake.ID        `json:"id"`
	GuildID  snowflake.ID        `json:"guild_id"`
	ParentID snowflake.ID        `json:"parent_id"`
	Type     discord.ChannelType `json:"type"`
}

func (EventThreadDelete) messageData() {}
func (EventThreadDelete) eventData()   {}

type EventThreadListSync struct {
	GuildID    snowflake.ID           `json:"guild_id"`
	ChannelIDs []snowflake.ID         `json:"channel_ids"`
	Threads    []discord.GuildThread  `json:"threads"`
	Members    []discord.ThreadMember `json:"members"`
}

func (EventThreadListSync) messageData() {}
func (EventThreadListSync) eventData()   {}

type EventThreadMemberUpdate struct {
	discord.ThreadMember
}

func (EventThreadMemberUpdate) messageData() {}
func (EventThreadMemberUpdate) eventData()   {}

type AddedThreadMember struct {
	discord.ThreadMember
	Member   discord.Member    `json:"member"`
	Presence *discord.Presence `json:"presence"`
}

type EventThreadMembersUpdate struct {
	ID               snowflake.ID        `json:"id"`
	GuildID          snowflake.ID        `json:"guild_id"`
	MemberCount      int                 `json:"member_count"`
	AddedMembers     []AddedThreadMember `json:"added_members"`
	RemovedMemberIDs []snowflake.ID      `json:"removed_member_ids"`
}

func (EventThreadMembersUpdate) messageData() {}
func (EventThreadMembersUpdate) eventData()   {}

type EventGuildCreate struct {
	discord.GatewayGuild
}

func (EventGuildCreate) messageData() {}
func (EventGuildCreate) eventData()   {}

type EventGuildUpdate struct {
	discord.GatewayGuild
}

func (EventGuildUpdate) messageData() {}
func (EventGuildUpdate) eventData()   {}

type EventGuildDelete struct {
	discord.GatewayGuild
}

func (EventGuildDelete) messageData() {}
func (EventGuildDelete) eventData()   {}

type EventMessageReactionAdd struct {
	UserID    snowflake.ID          `json:"user_id"`
	ChannelID snowflake.ID          `json:"channel_id"`
	MessageID snowflake.ID          `json:"message_id"`
	GuildID   *snowflake.ID         `json:"guild_id"`
	Member    *discord.Member       `json:"member"`
	Emoji     discord.ReactionEmoji `json:"emoji"`
}

func (EventMessageReactionAdd) messageData() {}
func (EventMessageReactionAdd) eventData()   {}

type EventMessageReactionRemove struct {
	UserID    snowflake.ID          `json:"user_id"`
	ChannelID snowflake.ID          `json:"channel_id"`
	MessageID snowflake.ID          `json:"message_id"`
	GuildID   *snowflake.ID         `json:"guild_id"`
	Emoji     discord.ReactionEmoji `json:"emoji"`
}

func (EventMessageReactionRemove) messageData() {}
func (EventMessageReactionRemove) eventData()   {}

type EventMessageReactionRemoveEmoji struct {
	ChannelID snowflake.ID          `json:"channel_id"`
	MessageID snowflake.ID          `json:"message_id"`
	GuildID   *snowflake.ID         `json:"guild_id"`
	Emoji     discord.ReactionEmoji `json:"emoji"`
}

func (EventMessageReactionRemoveEmoji) messageData() {}
func (EventMessageReactionRemoveEmoji) eventData()   {}

type EventMessageReactionRemoveAll struct {
	ChannelID snowflake.ID  `json:"channel_id"`
	MessageID snowflake.ID  `json:"message_id"`
	GuildID   *snowflake.ID `json:"guild_id"`
}

func (EventMessageReactionRemoveAll) messageData() {}
func (EventMessageReactionRemoveAll) eventData()   {}

type EventChannelPinsUpdate struct {
	GuildID          *snowflake.ID `json:"guild_id"`
	ChannelID        snowflake.ID  `json:"channel_id"`
	LastPinTimestamp *time.Time    `json:"last_pin_timestamp"`
}

func (EventChannelPinsUpdate) messageData() {}
func (EventChannelPinsUpdate) eventData()   {}

type EventGuildMembersChunk struct {
	GuildID    snowflake.ID       `json:"guild_id"`
	Members    []discord.Member   `json:"members"`
	ChunkIndex int                `json:"chunk_index"`
	ChunkCount int                `json:"chunk_count"`
	NotFound   []snowflake.ID     `json:"not_found"`
	Presences  []discord.Presence `json:"presences"`
	Nonce      string             `json:"nonce"`
}

func (EventGuildMembersChunk) messageData() {}
func (EventGuildMembersChunk) eventData()   {}

type EventGuildBanAdd struct {
	GuildID snowflake.ID `json:"guild_id"`
	User    discord.User `json:"user"`
}

func (EventGuildBanAdd) messageData() {}
func (EventGuildBanAdd) eventData()   {}

type EventGuildBanRemove struct {
	GuildID snowflake.ID `json:"guild_id"`
	User    discord.User `json:"user"`
}

func (EventGuildBanRemove) messageData() {}
func (EventGuildBanRemove) eventData()   {}

type EventGuildEmojisUpdate struct {
	GuildID snowflake.ID    `json:"guild_id"`
	Emojis  []discord.Emoji `json:"emojis"`
}

func (EventGuildEmojisUpdate) messageData() {}
func (EventGuildEmojisUpdate) eventData()   {}

type EventGuildStickersUpdate struct {
	GuildID  snowflake.ID      `json:"guild_id"`
	Stickers []discord.Sticker `json:"stickers"`
}

func (EventGuildStickersUpdate) messageData() {}
func (EventGuildStickersUpdate) eventData()   {}

type EventGuildIntegrationsUpdate struct {
	GuildID snowflake.ID `json:"guild_id"`
}

func (EventGuildIntegrationsUpdate) messageData() {}
func (EventGuildIntegrationsUpdate) eventData()   {}

type EventGuildMemberAdd struct {
	discord.Member
}

func (EventGuildMemberAdd) messageData() {}
func (EventGuildMemberAdd) eventData()   {}

type EventGuildMemberUpdate struct {
	discord.Member
}

func (EventGuildMemberUpdate) messageData() {}
func (EventGuildMemberUpdate) eventData()   {}

type EventGuildMemberRemove struct {
	GuildID snowflake.ID `json:"guild_id"`
	User    discord.User `json:"user"`
}

func (EventGuildMemberRemove) messageData() {}
func (EventGuildMemberRemove) eventData()   {}

type EventGuildRoleCreate struct {
	GuildID snowflake.ID `json:"guild_id"`
	Role    discord.Role `json:"role"`
}

func (EventGuildRoleCreate) messageData() {}
func (EventGuildRoleCreate) eventData()   {}

type EventGuildRoleDelete struct {
	GuildID snowflake.ID `json:"guild_id"`
	RoleID  snowflake.ID `json:"role_id"`
}

func (EventGuildRoleDelete) messageData() {}
func (EventGuildRoleDelete) eventData()   {}

type EventGuildRoleUpdate struct {
	GuildID snowflake.ID `json:"guild_id"`
	Role    discord.Role `json:"role"`
}

func (EventGuildRoleUpdate) messageData() {}
func (EventGuildRoleUpdate) eventData()   {}

type EventGuildScheduledEventCreate struct {
	discord.GuildScheduledEvent
}

func (EventGuildScheduledEventCreate) messageData() {}
func (EventGuildScheduledEventCreate) eventData()   {}

type EventGuildScheduledEventUpdate struct {
	discord.GuildScheduledEvent
}

func (EventGuildScheduledEventUpdate) messageData() {}
func (EventGuildScheduledEventUpdate) eventData()   {}

type EventGuildScheduledEventDelete struct {
	discord.GuildScheduledEvent
}

func (EventGuildScheduledEventDelete) messageData() {}
func (EventGuildScheduledEventDelete) eventData()   {}

type EventGuildScheduledEventUserAdd struct {
	GuildScheduledEventID snowflake.ID `json:"guild_scheduled_event_id"`
	UserID                snowflake.ID `json:"user_id"`
	GuildID               snowflake.ID `json:"guild_id"`
}

func (EventGuildScheduledEventUserAdd) messageData() {}
func (EventGuildScheduledEventUserAdd) eventData()   {}

type EventGuildScheduledEventUserRemove struct {
	GuildScheduledEventID snowflake.ID `json:"guild_scheduled_event_id"`
	UserID                snowflake.ID `json:"user_id"`
	GuildID               snowflake.ID `json:"guild_id"`
}

func (EventGuildScheduledEventUserRemove) messageData() {}
func (EventGuildScheduledEventUserRemove) eventData()   {}

type EventInteractionCreate struct {
	discord.Interaction
}

func (e *EventInteractionCreate) UnmarshalJSON(data []byte) error {
	var interaction discord.UnmarshalInteraction
	if err := json.Unmarshal(data, &interaction); err != nil {
		return err
	}
	e.Interaction = interaction.Interaction
	return nil
}

func (EventInteractionCreate) messageData() {}
func (EventInteractionCreate) eventData()   {}

type EventInviteCreate struct {
	discord.Invite
}

func (EventInviteCreate) messageData() {}
func (EventInviteCreate) eventData()   {}

type EventInviteDelete struct {
	ChannelID snowflake.ID  `json:"channel_id"`
	GuildID   *snowflake.ID `json:"guild_id"`
	Code      string        `json:"code"`
}

func (EventInviteDelete) messageData() {}
func (EventInviteDelete) eventData()   {}

type EventMessageCreate struct {
	discord.Message
}

func (EventMessageCreate) messageData() {}
func (EventMessageCreate) eventData()   {}

type EventMessageUpdate struct {
	discord.Message
}

func (EventMessageUpdate) messageData() {}
func (EventMessageUpdate) eventData()   {}

type EventMessageDelete struct {
	ID        snowflake.ID  `json:"id"`
	ChannelID snowflake.ID  `json:"channel_id"`
	GuildID   *snowflake.ID `json:"guild_id,omitempty"`
}

func (EventMessageDelete) messageData() {}
func (EventMessageDelete) eventData()   {}

type EventMessageDeleteBulk struct {
	IDs       []snowflake.ID `json:"id"`
	ChannelID snowflake.ID   `json:"channel_id"`
	GuildID   *snowflake.ID  `json:"guild_id,omitempty"`
}

func (EventMessageDeleteBulk) messageData() {}
func (EventMessageDeleteBulk) eventData()   {}

type EventPresenceUpdate struct {
	discord.Presence
}

func (EventPresenceUpdate) messageData() {}
func (EventPresenceUpdate) eventData()   {}

type EventStageInstanceCreate struct {
	discord.StageInstance
}

func (EventStageInstanceCreate) messageData() {}
func (EventStageInstanceCreate) eventData()   {}

type EventStageInstanceUpdate struct {
	discord.StageInstance
}

func (EventStageInstanceUpdate) messageData() {}
func (EventStageInstanceUpdate) eventData()   {}

type EventStageInstanceDelete struct {
	discord.StageInstance
}

func (EventStageInstanceDelete) messageData() {}
func (EventStageInstanceDelete) eventData()   {}

type EventTypingStart struct {
	ChannelID snowflake.ID    `json:"channel_id"`
	GuildID   *snowflake.ID   `json:"guild_id,omitempty"`
	UserID    snowflake.ID    `json:"user_id"`
	Timestamp time.Time       `json:"timestamp"`
	Member    *discord.Member `json:"member,omitempty"`
	User      discord.User    `json:"user"`
}

func (e *EventTypingStart) UnmarshalJSON(data []byte) error {
	type typingStartEvent EventTypingStart
	var v struct {
		Timestamp int64 `json:"timestamp"`
		typingStartEvent
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*e = EventTypingStart(v.typingStartEvent)
	e.Timestamp = time.Unix(v.Timestamp, 0)
	return nil
}

func (EventTypingStart) messageData() {}
func (EventTypingStart) eventData()   {}

type EventUserUpdate struct {
	discord.OAuth2User
}

func (EventUserUpdate) messageData() {}
func (EventUserUpdate) eventData()   {}

type EventVoiceStateUpdate struct {
	discord.VoiceState
	Member discord.Member `json:"member"`
}

func (EventVoiceStateUpdate) messageData() {}
func (EventVoiceStateUpdate) eventData()   {}

type EventVoiceServerUpdate struct {
	Token    string       `json:"token"`
	GuildID  snowflake.ID `json:"guild_id"`
	Endpoint *string      `json:"endpoint"`
}

func (EventVoiceServerUpdate) messageData() {}
func (EventVoiceServerUpdate) eventData()   {}

type EventWebhooksUpdate struct {
	GuildID   snowflake.ID `json:"guild_id"`
	ChannelID snowflake.ID `json:"channel_id"`
}

func (EventWebhooksUpdate) messageData() {}
func (EventWebhooksUpdate) eventData()   {}

type EventIntegrationCreate struct {
	discord.Integration
	GuildID snowflake.ID `json:"guild_id"`
}

func (e *EventIntegrationCreate) UnmarshalJSON(data []byte) error {
	type integrationCreateEvent EventIntegrationCreate
	var v struct {
		discord.UnmarshalIntegration
		integrationCreateEvent
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*e = EventIntegrationCreate(v.integrationCreateEvent)

	e.Integration = v.UnmarshalIntegration.Integration
	return nil
}

func (EventIntegrationCreate) messageData() {}
func (EventIntegrationCreate) eventData()   {}

type EventIntegrationUpdate struct {
	discord.Integration
	GuildID snowflake.ID `json:"guild_id"`
}

func (e *EventIntegrationUpdate) UnmarshalJSON(data []byte) error {
	type integrationUpdateEvent EventIntegrationUpdate
	var v struct {
		discord.UnmarshalIntegration
		integrationUpdateEvent
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*e = EventIntegrationUpdate(v.integrationUpdateEvent)

	e.Integration = v.UnmarshalIntegration.Integration
	return nil
}

func (EventIntegrationUpdate) messageData() {}
func (EventIntegrationUpdate) eventData()   {}

type EventIntegrationDelete struct {
	ID            snowflake.ID  `json:"id"`
	GuildID       snowflake.ID  `json:"guild_id"`
	ApplicationID *snowflake.ID `json:"application_id"`
}

func (EventIntegrationDelete) messageData() {}
func (EventIntegrationDelete) eventData()   {}

type EventAutoModerationRuleCreate struct {
	discord.AutoModerationRule
}

func (EventAutoModerationRuleCreate) messageData() {}
func (EventAutoModerationRuleCreate) eventData()   {}

type EventAutoModerationRuleUpdate struct {
	discord.AutoModerationRule
}

func (EventAutoModerationRuleUpdate) messageData() {}
func (EventAutoModerationRuleUpdate) eventData()   {}

type EventAutoModerationRuleDelete struct {
	discord.AutoModerationRule
}

func (EventAutoModerationRuleDelete) messageData() {}
func (EventAutoModerationRuleDelete) eventData()   {}

type EventAutoModerationActionExecution struct {
	GuildID              snowflake.ID                      `json:"guild_id"`
	Action               discord.AutoModerationAction      `json:"action"`
	RuleID               snowflake.ID                      `json:"rule_id"`
	RuleTriggerType      discord.AutoModerationTriggerType `json:"rule_trigger_type"`
	UserID               snowflake.ID                      `json:"user_id"`
	ChannelID            *snowflake.ID                     `json:"channel_id,omitempty"`
	MessageID            *snowflake.ID                     `json:"message_id,omitempty"`
	AlertSystemMessageID snowflake.ID                      `json:"alert_system_message_id"`
	Content              string                            `json:"content"`
	MatchedKeywords      *string                           `json:"matched_keywords"`
	MatchedContent       *string                           `json:"matched_content"`
}

func (EventAutoModerationActionExecution) messageData() {}
func (EventAutoModerationActionExecution) eventData()   {}

type EventRaw struct {
	EventType EventType
	Payload   io.Reader
}

func (EventRaw) messageData() {}
func (EventRaw) eventData()   {}
