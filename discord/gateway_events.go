package discord

import (
	"time"

	"github.com/disgoorg/disgo/json"
	"github.com/disgoorg/snowflake"
)

// GatewayEventReady is the event sent by discord when you successfully Identify
type GatewayEventReady struct {
	Version     int                `json:"v"`
	User        OAuth2User         `json:"user"`
	Guilds      []UnavailableGuild `json:"guilds"`
	SessionID   string             `json:"session_id"`
	Shard       []int              `json:"shard,omitempty"`
	Application PartialApplication `json:"application"`
}

type GatewayEventThreadCreate struct {
	GuildThread
	ThreadMember ThreadMember `json:"thread_member"`
}

func (e *GatewayEventThreadCreate) UnmarshalJSON(data []byte) error {
	var v struct {
		UnmarshalChannel
		ThreadMember ThreadMember `json:"thread_member"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	e.GuildThread = v.UnmarshalChannel.Channel.(GuildThread)
	e.ThreadMember = v.ThreadMember
	return nil
}

type GatewayEventThreadUpdate struct {
	GuildThread
}

func (e *GatewayEventThreadUpdate) UnmarshalJSON(data []byte) error {
	var v UnmarshalChannel
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	e.GuildThread = v.Channel.(GuildThread)
	return nil
}

type GatewayEventThreadDelete struct {
	ID       snowflake.Snowflake `json:"id"`
	GuildID  snowflake.Snowflake `json:"guild_id"`
	ParentID snowflake.Snowflake `json:"parent_id"`
	Type     ChannelType         `json:"type"`
}

type GatewayEventThreadListSync struct {
	GuildID    snowflake.Snowflake   `json:"guild_id"`
	ChannelIDs []snowflake.Snowflake `json:"channel_ids"`
	Threads    []GuildThread         `json:"threads"`
	Members    []ThreadMember        `json:"members"`
}

func (e *GatewayEventThreadListSync) UnmarshalJSON(data []byte) error {
	type gatewayEventThreadListSync GatewayEventThreadListSync
	var v struct {
		Threads []UnmarshalChannel `json:"threads"`
		gatewayEventThreadListSync
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*e = GatewayEventThreadListSync(v.gatewayEventThreadListSync)
	if len(v.Threads) > 0 {
		e.Threads = make([]GuildThread, len(v.Threads))
		for i := range v.Threads {
			e.Threads[i] = v.Threads[i].Channel.(GuildThread)
		}
	}
	return nil
}

type GatewayEventThreadMembersUpdate struct {
	ID               snowflake.Snowflake   `json:"id"`
	GuildID          snowflake.Snowflake   `json:"guild_id"`
	MemberCount      int                   `json:"member_count"`
	AddedMembers     []AddedThreadMember   `json:"added_members"`
	RemovedMemberIDs []snowflake.Snowflake `json:"removed_member_ids"`
}

type AddedThreadMember struct {
	ThreadMember
	Member   Member    `json:"member"`
	Presence *Presence `json:"presence"`
}

type GatewayEventMessageReactionAdd struct {
	UserID    snowflake.Snowflake  `json:"user_id"`
	ChannelID snowflake.Snowflake  `json:"channel_id"`
	MessageID snowflake.Snowflake  `json:"message_id"`
	GuildID   *snowflake.Snowflake `json:"guild_id"`
	Member    *Member              `json:"member"`
	Emoji     ReactionEmoji        `json:"emoji"`
}

type GatewayEventMessageReactionRemove struct {
	UserID    snowflake.Snowflake  `json:"user_id"`
	ChannelID snowflake.Snowflake  `json:"channel_id"`
	MessageID snowflake.Snowflake  `json:"message_id"`
	GuildID   *snowflake.Snowflake `json:"guild_id"`
	Emoji     ReactionEmoji        `json:"emoji"`
}

type GatewayEventMessageReactionRemoveEmoji struct {
	ChannelID snowflake.Snowflake  `json:"channel_id"`
	MessageID snowflake.Snowflake  `json:"message_id"`
	GuildID   *snowflake.Snowflake `json:"guild_id"`
	Emoji     ReactionEmoji        `json:"emoji"`
}

type GatewayEventMessageReactionRemoveAll struct {
	ChannelID snowflake.Snowflake  `json:"channel_id"`
	MessageID snowflake.Snowflake  `json:"message_id"`
	GuildID   *snowflake.Snowflake `json:"guild_id"`
}

type GatewayEventChannelPinsUpdate struct {
	GuildID          *snowflake.Snowflake `json:"guild_id"`
	ChannelID        snowflake.Snowflake  `json:"channel_id"`
	LastPinTimestamp *Time                `json:"last_pin_timestamp"`
}

type GatewayEventGuildMembersChunk struct {
	GuildID    snowflake.Snowflake   `json:"guild_id"`
	Members    []Member              `json:"members"`
	ChunkIndex int                   `json:"chunk_index"`
	ChunkCount int                   `json:"chunk_count"`
	NotFound   []snowflake.Snowflake `json:"not_found"`
	Presences  []Presence            `json:"presences"`
	Nonce      string                `json:"nonce"`
}

type GatewayEventGuildBanAdd struct {
	GuildID snowflake.Snowflake `json:"guild_id"`
	User    User                `json:"user"`
}

type GatewayEventGuildBanRemove struct {
	GuildID snowflake.Snowflake `json:"guild_id"`
	User    User                `json:"user"`
}

type GatewayEventGuildEmojisUpdate struct {
	GuildID snowflake.Snowflake `json:"guild_id"`
	Emojis  []Emoji             `json:"emojis"`
}

type GatewayEventGuildStickersUpdate struct {
	GuildID  snowflake.Snowflake `json:"guild_id"`
	Stickers []Sticker           `json:"stickers"`
}

type GatewayEventGuildIntegrationsUpdate struct {
	GuildID snowflake.Snowflake `json:"guild_id"`
}

type GatewayEventGuildMemberAdd struct {
	GuildID snowflake.Snowflake `json:"guild_id"`
	Member
}

type GatewayEventGuildMemberRemove struct {
	GuildID snowflake.Snowflake `json:"guild_id"`
	User    User                `json:"user"`
}

type GatewayEventGuildRoleCreate struct {
	GuildID snowflake.Snowflake `json:"guild_id"`
	Role    Role                `json:"role"`
}

type GatewayEventGuildRoleDelete struct {
	GuildID snowflake.Snowflake `json:"guild_id"`
	RoleID  snowflake.Snowflake `json:"role_id"`
}

type GatewayEventGuildRoleUpdate struct {
	GuildID snowflake.Snowflake `json:"guild_id"`
	Role    Role                `json:"role"`
}

type GatewayEventGuildScheduledEventUser struct {
	GuildScheduledEventID snowflake.Snowflake `json:"guild_scheduled_event_id"`
	UserID                snowflake.Snowflake `json:"user_id"`
	GuildID               snowflake.Snowflake `json:"guild_id"`
}

type GatewayEventInviteDelete struct {
	ChannelID snowflake.Snowflake  `json:"channel_id"`
	GuildID   *snowflake.Snowflake `json:"guild_id"`
	Code      string               `json:"code"`
}

type GatewayEventMessageDelete struct {
	ID        snowflake.Snowflake  `json:"id"`
	ChannelID snowflake.Snowflake  `json:"channel_id"`
	GuildID   *snowflake.Snowflake `json:"guild_id,omitempty"`
}

type GatewayEventMessageDeleteBulk struct {
	IDs       []snowflake.Snowflake `json:"id"`
	ChannelID snowflake.Snowflake   `json:"channel_id"`
	GuildID   *snowflake.Snowflake  `json:"guild_id,omitempty"`
}

type GatewayEventTypingStart struct {
	ChannelID snowflake.Snowflake  `json:"channel_id"`
	GuildID   *snowflake.Snowflake `json:"guild_id,omitempty"`
	UserID    snowflake.Snowflake  `json:"user_id"`
	Timestamp time.Time            `json:"timestamp"`
	Member    *Member              `json:"member,omitempty"`
	User      User                 `json:"user"`
}

func (e *GatewayEventTypingStart) UnmarshalJSON(data []byte) error {
	type typingStartGatewayEvent GatewayEventTypingStart
	var v struct {
		Timestamp int64 `json:"timestamp"`
		typingStartGatewayEvent
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*e = GatewayEventTypingStart(v.typingStartGatewayEvent)
	e.Timestamp = time.Unix(v.Timestamp, 0)
	return nil
}

type GatewayEventWebhooksUpdate struct {
	GuildID   snowflake.Snowflake `json:"guild_id"`
	ChannelID snowflake.Snowflake `json:"channel_id"`
}

type GatewayEventIntegrationCreate struct {
	Integration
	GuildID snowflake.Snowflake `json:"guild_id"`
}

func (e *GatewayEventIntegrationCreate) UnmarshalJSON(data []byte) error {
	type integrationCreateGatewayEvent GatewayEventIntegrationCreate
	var v struct {
		UnmarshalIntegration
		integrationCreateGatewayEvent
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*e = GatewayEventIntegrationCreate(v.integrationCreateGatewayEvent)

	e.Integration = v.UnmarshalIntegration.Integration
	return nil
}

type GatewayEventIntegrationUpdate struct {
	Integration
	GuildID snowflake.Snowflake `json:"guild_id"`
}

func (e *GatewayEventIntegrationUpdate) UnmarshalJSON(data []byte) error {
	type integrationUpdateGatewayEvent GatewayEventIntegrationUpdate
	var v struct {
		UnmarshalIntegration
		integrationUpdateGatewayEvent
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*e = GatewayEventIntegrationUpdate(v.integrationUpdateGatewayEvent)

	e.Integration = v.UnmarshalIntegration.Integration
	return nil
}

type GatewayEventIntegrationDelete struct {
	ID            snowflake.Snowflake  `json:"id"`
	GuildID       snowflake.Snowflake  `json:"guild_id"`
	ApplicationID *snowflake.Snowflake `json:"application_id"`
}
