package gateway

import (
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/json"
	"github.com/disgoorg/snowflake/v2"
)

// EventReady is the event sent by discord when you successfully Identify
type EventReady struct {
	Version     int                        `json:"v"`
	User        discord.OAuth2User         `json:"user"`
	Guilds      []discord.UnavailableGuild `json:"guilds"`
	SessionID   string                     `json:"session_id"`
	Shard       []int                      `json:"shard,omitempty"`
	Application discord.PartialApplication `json:"application"`
}

type EventThreadCreate struct {
	discord.GuildThread
	ThreadMember discord.ThreadMember `json:"thread_member"`
}

func (e *EventThreadCreate) UnmarshalJSON(data []byte) error {
	var v struct {
		discord.UnmarshalChannel
		ThreadMember discord.ThreadMember `json:"thread_member"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	e.GuildThread = v.UnmarshalChannel.Channel.(discord.GuildThread)
	e.ThreadMember = v.ThreadMember
	return nil
}

type EventThreadUpdate struct {
	discord.GuildThread
}

func (e *EventThreadUpdate) UnmarshalJSON(data []byte) error {
	var v discord.UnmarshalChannel
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	e.GuildThread = v.Channel.(discord.GuildThread)
	return nil
}

type EventThreadDelete struct {
	ID       snowflake.ID        `json:"id"`
	GuildID  snowflake.ID        `json:"guild_id"`
	ParentID snowflake.ID        `json:"parent_id"`
	Type     discord.ChannelType `json:"type"`
}

type EventThreadListSync struct {
	GuildID    snowflake.ID           `json:"guild_id"`
	ChannelIDs []snowflake.ID         `json:"channel_ids"`
	Threads    []discord.GuildThread  `json:"threads"`
	Members    []discord.ThreadMember `json:"members"`
}

func (e *EventThreadListSync) UnmarshalJSON(data []byte) error {
	type gatewayEventThreadListSync EventThreadListSync
	var v struct {
		Threads []discord.UnmarshalChannel `json:"threads"`
		gatewayEventThreadListSync
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*e = EventThreadListSync(v.gatewayEventThreadListSync)
	if len(v.Threads) > 0 {
		e.Threads = make([]discord.GuildThread, len(v.Threads))
		for i := range v.Threads {
			e.Threads[i] = v.Threads[i].Channel.(discord.GuildThread)
		}
	}
	return nil
}

type EventThreadMembersUpdate struct {
	ID               snowflake.ID        `json:"id"`
	GuildID          snowflake.ID        `json:"guild_id"`
	MemberCount      int                 `json:"member_count"`
	AddedMembers     []AddedThreadMember `json:"added_members"`
	RemovedMemberIDs []snowflake.ID      `json:"removed_member_ids"`
}

type AddedThreadMember struct {
	discord.ThreadMember
	Member   discord.Member    `json:"member"`
	Presence *discord.Presence `json:"presence"`
}

type EventMessageReactionAdd struct {
	UserID    snowflake.ID          `json:"user_id"`
	ChannelID snowflake.ID          `json:"channel_id"`
	MessageID snowflake.ID          `json:"message_id"`
	GuildID   *snowflake.ID         `json:"guild_id"`
	Member    *discord.Member       `json:"member"`
	Emoji     discord.ReactionEmoji `json:"emoji"`
}

type EventMessageReactionRemove struct {
	UserID    snowflake.ID          `json:"user_id"`
	ChannelID snowflake.ID          `json:"channel_id"`
	MessageID snowflake.ID          `json:"message_id"`
	GuildID   *snowflake.ID         `json:"guild_id"`
	Emoji     discord.ReactionEmoji `json:"emoji"`
}

type EventMessageReactionRemoveEmoji struct {
	ChannelID snowflake.ID          `json:"channel_id"`
	MessageID snowflake.ID          `json:"message_id"`
	GuildID   *snowflake.ID         `json:"guild_id"`
	Emoji     discord.ReactionEmoji `json:"emoji"`
}

type EventMessageReactionRemoveAll struct {
	ChannelID snowflake.ID  `json:"channel_id"`
	MessageID snowflake.ID  `json:"message_id"`
	GuildID   *snowflake.ID `json:"guild_id"`
}

type EventChannelPinsUpdate struct {
	GuildID          *snowflake.ID `json:"guild_id"`
	ChannelID        snowflake.ID  `json:"channel_id"`
	LastPinTimestamp *time.Time    `json:"last_pin_timestamp"`
}

type EventGuildMembersChunk struct {
	GuildID    snowflake.ID       `json:"guild_id"`
	Members    []discord.Member   `json:"members"`
	ChunkIndex int                `json:"chunk_index"`
	ChunkCount int                `json:"chunk_count"`
	NotFound   []snowflake.ID     `json:"not_found"`
	Presences  []discord.Presence `json:"presences"`
	Nonce      string             `json:"nonce"`
}

type EventGuildBanAdd struct {
	GuildID snowflake.ID `json:"guild_id"`
	User    discord.User `json:"user"`
}

type EventGuildBanRemove struct {
	GuildID snowflake.ID `json:"guild_id"`
	User    discord.User `json:"user"`
}

type EventGuildEmojisUpdate struct {
	GuildID snowflake.ID    `json:"guild_id"`
	Emojis  []discord.Emoji `json:"emojis"`
}

type EventGuildStickersUpdate struct {
	GuildID  snowflake.ID      `json:"guild_id"`
	Stickers []discord.Sticker `json:"stickers"`
}

type EventGuildIntegrationsUpdate struct {
	GuildID snowflake.ID `json:"guild_id"`
}

type EventGuildMemberRemove struct {
	GuildID snowflake.ID `json:"guild_id"`
	User    discord.User `json:"user"`
}

type EventGuildRoleCreate struct {
	GuildID snowflake.ID `json:"guild_id"`
	Role    discord.Role `json:"role"`
}

type EventGuildRoleDelete struct {
	GuildID snowflake.ID `json:"guild_id"`
	RoleID  snowflake.ID `json:"role_id"`
}

type EventGuildRoleUpdate struct {
	GuildID snowflake.ID `json:"guild_id"`
	Role    discord.Role `json:"role"`
}

type EventGuildScheduledEventUser struct {
	GuildScheduledEventID snowflake.ID `json:"guild_scheduled_event_id"`
	UserID                snowflake.ID `json:"user_id"`
	GuildID               snowflake.ID `json:"guild_id"`
}

type EventInviteDelete struct {
	ChannelID snowflake.ID  `json:"channel_id"`
	GuildID   *snowflake.ID `json:"guild_id"`
	Code      string        `json:"code"`
}

type EventMessageDelete struct {
	ID        snowflake.ID  `json:"id"`
	ChannelID snowflake.ID  `json:"channel_id"`
	GuildID   *snowflake.ID `json:"guild_id,omitempty"`
}

type EventMessageDeleteBulk struct {
	IDs       []snowflake.ID `json:"id"`
	ChannelID snowflake.ID   `json:"channel_id"`
	GuildID   *snowflake.ID  `json:"guild_id,omitempty"`
}

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

type EventWebhooksUpdate struct {
	GuildID   snowflake.ID `json:"guild_id"`
	ChannelID snowflake.ID `json:"channel_id"`
}

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

type EventIntegrationDelete struct {
	ID            snowflake.ID  `json:"id"`
	GuildID       snowflake.ID  `json:"guild_id"`
	ApplicationID *snowflake.ID `json:"application_id"`
}

type EventAutoModerationRuleCreate struct {
	discord.AutoModerationRule
}

type EventAutoModerationRuleUpdate struct {
	discord.AutoModerationRule
}

type EventAutoModerationRuleDelete struct {
	discord.AutoModerationRule
}

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
