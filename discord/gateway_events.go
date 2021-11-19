package discord

import (
	"time"

	"github.com/DisgoOrg/disgo/json"
)

// GatewayPayload raw GatewayEvent type
type GatewayPayload struct {
	Op GatewayOpcode    `json:"op"`
	S  int              `json:"s,omitempty"`
	T  GatewayEventType `json:"t,omitempty"`
	D  json.RawMessage  `json:"d,omitempty"`
}

// GatewayEventReady is the event sent by discord when you successfully Identify
type GatewayEventReady struct {
	Version     int                `json:"v"`
	SelfUser    OAuth2User         `json:"user"`
	Guilds      []UnavailableGuild `json:"guilds"`
	SessionID   string             `json:"session_id"`
	Shard       []int              `json:"shard,omitempty"`
	Application PartialApplication `json:"application"`
}

// GatewayEventHello is sent when we connect to the gateway
type GatewayEventHello struct {
	HeartbeatInterval time.Duration `json:"heartbeat_interval"`
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
	ID       Snowflake   `json:"id"`
	GuildID  Snowflake   `json:"guild_id"`
	ParentID Snowflake   `json:"parent_id"`
	Type     ChannelType `json:"type"`
}

type GatewayEventThreadListSync struct {
	GuildID    Snowflake      `json:"guild_id"`
	ChannelIDs []Snowflake    `json:"channel_ids"`
	Threads    []Channel      `json:"threads"`
	Members    []ThreadMember `json:"members"`
}

type GatewayEventThreadMemberUpdate struct {
	ThreadMember
}

type GatewayEventThreadMembersUpdate struct {
	ID               Snowflake                   `json:"id"`
	GuildID          Snowflake                   `json:"guild_id"`
	MemberCount      int                         `json:"member_count"`
	AddedMembers     []ThreadMembersUpdateMember `json:"added_members"`
	RemovedMemberIDs []Snowflake                 `json:"removed_member_ids"`
}

type ThreadMembersUpdateMember struct {
	ThreadMember
	Member   Member    `json:"member"`
	Presence *Presence `json:"presence"`
}

type GatewayEventMessageReactionAdd struct {
	UserID    Snowflake     `json:"user_id"`
	ChannelID Snowflake     `json:"channel_id"`
	MessageID Snowflake     `json:"message_id"`
	GuildID   *Snowflake    `json:"guild_id"`
	Member    *Member       `json:"member"`
	Emoji     ReactionEmoji `json:"emoji"`
}

type GatewayEventMessageReactionRemove struct {
	UserID    Snowflake     `json:"user_id"`
	ChannelID Snowflake     `json:"channel_id"`
	MessageID Snowflake     `json:"message_id"`
	GuildID   *Snowflake    `json:"guild_id"`
	Emoji     ReactionEmoji `json:"emoji"`
}

type GatewayEventMessageReactionRemoveEmoji struct {
	ChannelID Snowflake     `json:"channel_id"`
	MessageID Snowflake     `json:"message_id"`
	GuildID   *Snowflake    `json:"guild_id"`
	Emoji     ReactionEmoji `json:"emoji"`
}

type GatewayEventMessageReactionRemoveAll struct {
	ChannelID Snowflake  `json:"channel_id"`
	MessageID Snowflake  `json:"message_id"`
	GuildID   *Snowflake `json:"guild_id"`
}
