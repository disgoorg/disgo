package discord

import (
	"github.com/disgoorg/disgo/json"
	"github.com/disgoorg/snowflake"
)

type ChannelCreate interface {
	json.Marshaler
	Type() ChannelType
	channelCreate()
}

type GuildChannelCreate interface {
	ChannelCreate
	guildChannelCreate()
}

var (
	_ ChannelCreate      = (*GuildTextChannelCreate)(nil)
	_ GuildChannelCreate = (*GuildTextChannelCreate)(nil)
)

type GuildTextChannelCreate struct {
	Name                       string                `json:"name"`
	Topic                      string                `json:"topic,omitempty"`
	RateLimitPerUser           int                   `json:"rate_limit_per_user,omitempty"`
	Position                   int                   `json:"position,omitempty"`
	PermissionOverwrites       []PermissionOverwrite `json:"permission_overwrites,omitempty"`
	ParentID                   snowflake.Snowflake   `json:"parent_id,omitempty"`
	NSFW                       bool                  `json:"nsfw,omitempty"`
	DefaultAutoArchiveDuration AutoArchiveDuration   `json:"default_auto_archive_days,omitempty"`
}

func (c GuildTextChannelCreate) Type() ChannelType {
	return ChannelTypeGuildText
}

func (c GuildTextChannelCreate) MarshalJSON() ([]byte, error) {
	type guildTextChannelCreate GuildTextChannelCreate
	return json.Marshal(struct {
		Type ChannelType `json:"type"`
		guildTextChannelCreate
	}{
		Type:                   c.Type(),
		guildTextChannelCreate: guildTextChannelCreate(c),
	})
}

func (GuildTextChannelCreate) channelCreate()      {}
func (GuildTextChannelCreate) guildChannelCreate() {}

var (
	_ ChannelCreate      = (*GuildVoiceChannelCreate)(nil)
	_ GuildChannelCreate = (*GuildVoiceChannelCreate)(nil)
)

type GuildVoiceChannelCreate struct {
	Name                 string                `json:"name"`
	Topic                string                `json:"topic,omitempty"`
	Bitrate              int                   `json:"bitrate,omitempty"`
	UserLimit            int                   `json:"user_limit,omitempty"`
	Position             int                   `json:"position,omitempty"`
	PermissionOverwrites []PermissionOverwrite `json:"permission_overwrites,omitempty"`
	ParentID             snowflake.Snowflake   `json:"parent_id,omitempty"`
}

func (c GuildVoiceChannelCreate) Type() ChannelType {
	return ChannelTypeGuildVoice
}

func (c GuildVoiceChannelCreate) MarshalJSON() ([]byte, error) {
	type guildVoiceChannelCreate GuildVoiceChannelCreate
	return json.Marshal(struct {
		Type ChannelType `json:"type"`
		guildVoiceChannelCreate
	}{
		Type:                    c.Type(),
		guildVoiceChannelCreate: guildVoiceChannelCreate(c),
	})
}

func (GuildVoiceChannelCreate) channelCreate()      {}
func (GuildVoiceChannelCreate) guildChannelCreate() {}

var (
	_ ChannelCreate      = (*GuildCategoryChannelCreate)(nil)
	_ GuildChannelCreate = (*GuildCategoryChannelCreate)(nil)
)

type GuildCategoryChannelCreate struct {
	Name                 string                `json:"name"`
	Topic                string                `json:"topic,omitempty"`
	Position             int                   `json:"position,omitempty"`
	PermissionOverwrites []PermissionOverwrite `json:"permission_overwrites,omitempty"`
}

func (c GuildCategoryChannelCreate) Type() ChannelType {
	return ChannelTypeGuildCategory
}

func (c GuildCategoryChannelCreate) MarshalJSON() ([]byte, error) {
	type guildCategoryChannelCreate GuildCategoryChannelCreate
	return json.Marshal(struct {
		Type ChannelType `json:"type"`
		guildCategoryChannelCreate
	}{
		Type:                       c.Type(),
		guildCategoryChannelCreate: guildCategoryChannelCreate(c),
	})
}

func (GuildCategoryChannelCreate) channelCreate()      {}
func (GuildCategoryChannelCreate) guildChannelCreate() {}

var (
	_ ChannelCreate      = (*GuildNewsChannelCreate)(nil)
	_ GuildChannelCreate = (*GuildNewsChannelCreate)(nil)
)

type GuildNewsChannelCreate struct {
	Name                       string                `json:"name"`
	Topic                      string                `json:"topic,omitempty"`
	RateLimitPerUser           int                   `json:"rate_limit_per_user,omitempty"`
	Position                   int                   `json:"position,omitempty"`
	PermissionOverwrites       []PermissionOverwrite `json:"permission_overwrites,omitempty"`
	ParentID                   snowflake.Snowflake   `json:"parent_id,omitempty"`
	NSFW                       bool                  `json:"nsfw,omitempty"`
	DefaultAutoArchiveDuration AutoArchiveDuration   `json:"default_auto_archive_days,omitempty"`
}

func (c GuildNewsChannelCreate) Type() ChannelType {
	return ChannelTypeGuildNews
}

func (c GuildNewsChannelCreate) MarshalJSON() ([]byte, error) {
	type guildNewsChannelCreate GuildNewsChannelCreate
	return json.Marshal(struct {
		Type ChannelType `json:"type"`
		guildNewsChannelCreate
	}{
		Type:                   c.Type(),
		guildNewsChannelCreate: guildNewsChannelCreate(c),
	})
}

func (GuildNewsChannelCreate) channelCreate()      {}
func (GuildNewsChannelCreate) guildChannelCreate() {}

var (
	_ ChannelCreate      = (*GuildStageVoiceChannelCreate)(nil)
	_ GuildChannelCreate = (*GuildStageVoiceChannelCreate)(nil)
)

type GuildStageVoiceChannelCreate struct {
	Name                 string                `json:"name"`
	Topic                string                `json:"topic,omitempty"`
	Bitrate              int                   `json:"bitrate,omitempty"`
	UserLimit            int                   `json:"user_limit,omitempty"`
	Position             int                   `json:"position,omitempty"`
	PermissionOverwrites []PermissionOverwrite `json:"permission_overwrites,omitempty"`
	ParentID             snowflake.Snowflake   `json:"parent_id,omitempty"`
}

func (c GuildStageVoiceChannelCreate) Type() ChannelType {
	return ChannelTypeGuildStageVoice
}

func (c GuildStageVoiceChannelCreate) MarshalJSON() ([]byte, error) {
	type guildStageChannelCreate GuildStageVoiceChannelCreate
	return json.Marshal(struct {
		Type ChannelType `json:"type"`
		guildStageChannelCreate
	}{
		Type:                    c.Type(),
		guildStageChannelCreate: guildStageChannelCreate(c),
	})
}

func (GuildStageVoiceChannelCreate) channelCreate()      {}
func (GuildStageVoiceChannelCreate) guildChannelCreate() {}

type DMChannelCreate struct {
	RecipientID snowflake.Snowflake `json:"recipient_id"`
}
