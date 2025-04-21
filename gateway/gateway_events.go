package gateway

import (
	"io"
	"time"

	"github.com/disgoorg/json/v2"
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
)

type EventData interface {
	MessageData
	eventData()
}

// EventUnknown is an event that is not known to disgo
type EventUnknown json.RawMessage

func (e EventUnknown) MarshalJSON() ([]byte, error) {
	return json.RawMessage(e).MarshalJSON()
}

func (e *EventUnknown) UnmarshalJSON(data []byte) error {
	return (*json.RawMessage)(e).UnmarshalJSON(data)
}

func (EventUnknown) messageData() {}
func (EventUnknown) eventData()   {}

// EventReady is the event sent by discord when you successfully Identify
type EventReady struct {
	Version          int                        `json:"v"`
	User             discord.OAuth2User         `json:"user"`
	Guilds           []discord.UnavailableGuild `json:"guilds"`
	SessionID        string                     `json:"session_id"`
	ResumeGatewayURL string                     `json:"resume_gateway_url"`
	Shard            [2]int                     `json:"shard,omitempty"`
	Application      discord.PartialApplication `json:"application"`
}

func (EventReady) messageData() {}
func (EventReady) eventData()   {}

// EventResumed is the event sent by discord when you successfully resume
type EventResumed struct{}

func (EventResumed) messageData() {}
func (EventResumed) eventData()   {}

type EventApplicationCommandPermissionsUpdate struct {
	discord.ApplicationCommandPermissions
}

func (EventApplicationCommandPermissionsUpdate) messageData() {}
func (EventApplicationCommandPermissionsUpdate) eventData()   {}

type EventChannelCreate struct {
	discord.GuildChannel
}

func (e *EventChannelCreate) UnmarshalJSON(data []byte) error {
	var v discord.UnmarshalChannel
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	e.GuildChannel = v.Channel.(discord.GuildChannel)
	return nil
}

func (EventChannelCreate) messageData() {}
func (EventChannelCreate) eventData()   {}

type EventChannelUpdate struct {
	discord.GuildChannel
}

func (e *EventChannelUpdate) UnmarshalJSON(data []byte) error {
	var v discord.UnmarshalChannel
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	e.GuildChannel = v.Channel.(discord.GuildChannel)
	return nil
}

func (EventChannelUpdate) messageData() {}
func (EventChannelUpdate) eventData()   {}

type EventChannelDelete struct {
	discord.GuildChannel
}

func (e *EventChannelDelete) UnmarshalJSON(data []byte) error {
	var v discord.UnmarshalChannel
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	e.GuildChannel = v.Channel.(discord.GuildChannel)
	return nil
}

func (EventChannelDelete) messageData() {}
func (EventChannelDelete) eventData()   {}

type EventThreadCreate struct {
	discord.GuildThread
	ThreadMember discord.ThreadMember `json:"thread_member"`
	NewlyCreated bool                 `json:"newly_created"`
}

func (e *EventThreadCreate) UnmarshalJSON(data []byte) error {
	var guildThread discord.GuildThread
	if err := json.Unmarshal(data, &guildThread); err != nil {
		return err
	}

	e.GuildThread = guildThread

	var v struct {
		ThreadMember discord.ThreadMember `json:"thread_member"`
		NewlyCreated bool                 `json:"newly_created"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	e.ThreadMember = v.ThreadMember
	e.NewlyCreated = v.NewlyCreated
	return nil
}

func (e EventThreadCreate) MarshalJSON() ([]byte, error) {
	data1, err := json.Marshal(e.GuildThread)
	if err != nil {
		return nil, err
	}

	data2, err := json.Marshal(struct {
		ThreadMember discord.ThreadMember `json:"thread_member"`
		NewlyCreated bool                 `json:"newly_created"`
	}{
		ThreadMember: e.ThreadMember,
		NewlyCreated: e.NewlyCreated,
	})
	if err != nil {
		return nil, err
	}

	return json.SimpleMerge(data1, data2)
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
	discord.UnavailableGuild
}

func (EventGuildDelete) messageData() {}
func (EventGuildDelete) eventData()   {}

type EventGuildAuditLogEntryCreate struct {
	discord.AuditLogEntry
	GuildID snowflake.ID `json:"guild_id"`
}

func (EventGuildAuditLogEntryCreate) messageData() {}
func (EventGuildAuditLogEntryCreate) eventData()   {}

type EventMessageReactionAdd struct {
	UserID          snowflake.ID                `json:"user_id"`
	ChannelID       snowflake.ID                `json:"channel_id"`
	MessageID       snowflake.ID                `json:"message_id"`
	GuildID         *snowflake.ID               `json:"guild_id"`
	Member          *discord.Member             `json:"member"`
	Emoji           discord.PartialEmoji        `json:"emoji"`
	MessageAuthorID *snowflake.ID               `json:"message_author_id"`
	BurstColors     []string                    `json:"burst_colors"`
	Burst           bool                        `json:"burst"`
	Type            discord.MessageReactionType `json:"type"`
}

func (e *EventMessageReactionAdd) UnmarshalJSON(data []byte) error {
	type eventMessageReactionAdd EventMessageReactionAdd
	var v eventMessageReactionAdd
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*e = EventMessageReactionAdd(v)
	if e.Member != nil && e.GuildID != nil {
		e.Member.GuildID = *e.GuildID
	}
	return nil
}

func (EventMessageReactionAdd) messageData() {}
func (EventMessageReactionAdd) eventData()   {}

type EventMessageReactionRemove struct {
	UserID      snowflake.ID                `json:"user_id"`
	ChannelID   snowflake.ID                `json:"channel_id"`
	MessageID   snowflake.ID                `json:"message_id"`
	GuildID     *snowflake.ID               `json:"guild_id"`
	Emoji       discord.PartialEmoji        `json:"emoji"`
	BurstColors []string                    `json:"burst_colors"`
	Burst       bool                        `json:"burst"`
	Type        discord.MessageReactionType `json:"type"`
}

func (EventMessageReactionRemove) messageData() {}
func (EventMessageReactionRemove) eventData()   {}

type EventMessageReactionRemoveEmoji struct {
	ChannelID snowflake.ID         `json:"channel_id"`
	MessageID snowflake.ID         `json:"message_id"`
	GuildID   *snowflake.ID        `json:"guild_id"`
	Emoji     discord.PartialEmoji `json:"emoji"`
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

func (e *EventGuildEmojisUpdate) UnmarshalJSON(data []byte) error {
	type eventGuildEmojisUpdate EventGuildEmojisUpdate
	var v eventGuildEmojisUpdate
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*e = EventGuildEmojisUpdate(v)
	for i := range e.Emojis {
		e.Emojis[i].GuildID = e.GuildID
	}
	return nil
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

func (e *EventGuildRoleCreate) UnmarshalJSON(data []byte) error {
	type eventGuildRoleCreate EventGuildRoleCreate
	var v eventGuildRoleCreate
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*e = EventGuildRoleCreate(v)
	e.Role.GuildID = e.GuildID
	return nil
}

func (e *EventGuildRoleCreate) MarshalJSON() ([]byte, error) {
	type eventGuildRoleCreate EventGuildRoleCreate
	e.GuildID = e.Role.GuildID
	return json.Marshal(eventGuildRoleCreate(*e))
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

func (e *EventGuildRoleUpdate) UnmarshalJSON(data []byte) error {
	type eventGuildRoleUpdate EventGuildRoleUpdate
	var v eventGuildRoleUpdate
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*e = EventGuildRoleUpdate(v)
	e.Role.GuildID = e.GuildID
	return nil
}

func (e *EventGuildRoleUpdate) MarshalJSON() ([]byte, error) {
	type eventGuildRoleUpdate EventGuildRoleUpdate
	e.GuildID = e.Role.GuildID
	return json.Marshal(eventGuildRoleUpdate(*e))
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

type EventGuildSoundboardSoundCreate struct {
	discord.SoundboardSound
}

func (EventGuildSoundboardSoundCreate) messageData() {}
func (EventGuildSoundboardSoundCreate) eventData()   {}

type EventGuildSoundboardSoundUpdate struct {
	discord.SoundboardSound
}

func (EventGuildSoundboardSoundUpdate) messageData() {}
func (EventGuildSoundboardSoundUpdate) eventData()   {}

type EventGuildSoundboardSoundDelete struct {
	SoundID snowflake.ID `json:"sound_id"`
	GuildID snowflake.ID `json:"guild_id"`
}

func (EventGuildSoundboardSoundDelete) messageData() {}
func (EventGuildSoundboardSoundDelete) eventData()   {}

type EventGuildSoundboardSoundsUpdate struct {
	SoundboardSounds []discord.SoundboardSound `json:"soundboard_sounds"`
	GuildID          snowflake.ID              `json:"guild_id"`
}

func (EventGuildSoundboardSoundsUpdate) messageData() {}
func (EventGuildSoundboardSoundsUpdate) eventData()   {}

type EventInteractionCreate struct {
	discord.Interaction
}

func (e *EventInteractionCreate) UnmarshalJSON(data []byte) error {
	interaction, err := discord.UnmarshalInteraction(data)
	if err != nil {
		return err
	}
	e.Interaction = interaction
	return nil
}

func (e EventInteractionCreate) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Interaction)
}

func (EventInteractionCreate) messageData() {}
func (EventInteractionCreate) eventData()   {}

type EventInviteCreate struct {
	ChannelID         snowflake.ID                `json:"channel_id"`
	Code              string                      `json:"code"`
	CreatedAt         time.Time                   `json:"created_at"`
	GuildID           *snowflake.ID               `json:"guild_id"`
	Inviter           *discord.User               `json:"inviter"`
	MaxAge            int                         `json:"max_age"`
	MaxUses           int                         `json:"max_uses"`
	TargetType        discord.InviteTargetType    `json:"target_type"`
	TargetUser        *discord.User               `json:"target_user"`
	TargetApplication *discord.PartialApplication `json:"target_application"`
	Temporary         bool                        `json:"temporary"`
	Uses              int                         `json:"uses"`
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
	IDs       []snowflake.ID `json:"ids"`
	ChannelID snowflake.ID   `json:"channel_id"`
	GuildID   *snowflake.ID  `json:"guild_id,omitempty"`
}

func (EventMessageDeleteBulk) messageData() {}
func (EventMessageDeleteBulk) eventData()   {}

type EventMessagePollVoteAdd struct {
	UserID    snowflake.ID  `json:"user_id"`
	ChannelID snowflake.ID  `json:"channel_id"`
	MessageID snowflake.ID  `json:"message_id"`
	GuildID   *snowflake.ID `json:"guild_id"`
	AnswerID  int           `json:"answer_id"`
}

func (EventMessagePollVoteAdd) messageData() {}
func (EventMessagePollVoteAdd) eventData()   {}

type EventMessagePollVoteRemove struct {
	UserID    snowflake.ID  `json:"user_id"`
	ChannelID snowflake.ID  `json:"channel_id"`
	MessageID snowflake.ID  `json:"message_id"`
	GuildID   *snowflake.ID `json:"guild_id"`
	AnswerID  int           `json:"answer_id"`
}

func (EventMessagePollVoteRemove) messageData() {}
func (EventMessagePollVoteRemove) eventData()   {}

type EventPresenceUpdate struct {
	discord.Presence
}

func (EventPresenceUpdate) messageData() {}
func (EventPresenceUpdate) eventData()   {}

type EventSoundboardSounds struct {
	SoundboardSounds []discord.SoundboardSound `json:"soundboard_sounds"`
	GuildID          snowflake.ID              `json:"guild_id"`
}

func (EventSoundboardSounds) messageData() {}
func (EventSoundboardSounds) eventData()   {}

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

type EventVoiceChannelEffectSend struct {
	ChannelID     snowflake.ID                             `json:"channel_id"`
	GuildID       snowflake.ID                             `json:"guild_id"`
	UserID        snowflake.ID                             `json:"user_id"`
	Emoji         *discord.Emoji                           `json:"emoji"`
	AnimationType *discord.VoiceChannelEffectAnimationType `json:"animation_type,omitempty"`
	AnimationID   *int                                     `json:"animation_id,omitempty"`
	SoundID       *int64                                   `json:"-"`
	SoundVolume   *float64                                 `json:"sound_volume,omitempty"`
}

func (e *EventVoiceChannelEffectSend) UnmarshalJSON(data []byte) error {
	type eventVoiceChannelEffectSend EventVoiceChannelEffectSend
	var v struct {
		SoundID *json.Number `json:"sound_id"`
		eventVoiceChannelEffectSend
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*e = EventVoiceChannelEffectSend(v.eventVoiceChannelEffectSend)
	if v.SoundID == nil {
		return nil
	}
	i, err := v.SoundID.Int64()
	if err != nil {
		return err
	}
	e.SoundID = &i
	return nil
}

func (EventVoiceChannelEffectSend) messageData() {}
func (EventVoiceChannelEffectSend) eventData()   {}

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
	var integration discord.UnmarshalIntegration
	if err := json.Unmarshal(data, &integration); err != nil {
		return err
	}

	var v struct {
		GuildID snowflake.ID `json:"guild_id"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	e.Integration = integration.Integration
	e.GuildID = v.GuildID
	return nil
}

func (e EventIntegrationCreate) MarshalJSON() ([]byte, error) {
	data1, err := json.Marshal(e.Integration)
	if err != nil {
		return nil, err
	}

	data2, err := json.Marshal(struct {
		GuildID snowflake.ID `json:"guild_id"`
	}{
		GuildID: e.GuildID,
	})
	if err != nil {
		return nil, err
	}

	return json.SimpleMerge(data1, data2)
}

func (EventIntegrationCreate) messageData() {}
func (EventIntegrationCreate) eventData()   {}

type EventIntegrationUpdate struct {
	discord.Integration
	GuildID snowflake.ID `json:"guild_id"`
}

func (e *EventIntegrationUpdate) UnmarshalJSON(data []byte) error {
	var integration discord.UnmarshalIntegration
	if err := json.Unmarshal(data, &integration); err != nil {
		return err
	}

	var v struct {
		GuildID snowflake.ID `json:"guild_id"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	e.Integration = integration.Integration
	e.GuildID = v.GuildID
	return nil
}

func (e EventIntegrationUpdate) MarshalJSON() ([]byte, error) {
	data1, err := json.Marshal(e.Integration)
	if err != nil {
		return nil, err
	}

	data2, err := json.Marshal(struct {
		GuildID snowflake.ID `json:"guild_id"`
	}{
		GuildID: e.GuildID,
	})
	if err != nil {
		return nil, err
	}

	return json.SimpleMerge(data1, data2)
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

type EventHeartbeatAck struct {
	LastHeartbeat time.Time
	NewHeartbeat  time.Time
}

func (EventHeartbeatAck) messageData() {}
func (EventHeartbeatAck) eventData()   {}

type EventEntitlementCreate struct {
	discord.Entitlement
}

func (EventEntitlementCreate) messageData() {}
func (EventEntitlementCreate) eventData()   {}

type EventEntitlementUpdate struct {
	discord.Entitlement
}

func (EventEntitlementUpdate) messageData() {}
func (EventEntitlementUpdate) eventData()   {}

type EventEntitlementDelete struct {
	discord.Entitlement
}

func (EventEntitlementDelete) messageData() {}
func (EventEntitlementDelete) eventData()   {}

type EventSubscriptionCreate struct {
	discord.Subscription
}

func (EventSubscriptionCreate) messageData() {}
func (EventSubscriptionCreate) eventData()   {}

type EventSubscriptionUpdate struct {
	discord.Subscription
}

func (EventSubscriptionUpdate) messageData() {}
func (EventSubscriptionUpdate) eventData()   {}

type EventSubscriptionDelete struct {
	discord.Subscription
}

func (EventSubscriptionDelete) messageData() {}
func (EventSubscriptionDelete) eventData()   {}
