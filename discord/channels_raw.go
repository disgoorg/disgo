package discord

import (
	"github.com/DisgoOrg/snowflake"
	"github.com/disgoorg/disgo/json"
)

type dmChannel struct {
	ID               snowflake.Snowflake  `json:"id"`
	Type             ChannelType          `json:"type"`
	LastMessageID    *snowflake.Snowflake `json:"last_message_id"`
	Recipients       []User               `json:"recipients"`
	LastPinTimestamp *Time                `json:"last_pin_timestamp"`
}

type guildTextChannel struct {
	ID                         snowflake.Snowflake   `json:"id"`
	Type                       ChannelType           `json:"type"`
	GuildID                    snowflake.Snowflake   `json:"guild_id"`
	Position                   int                   `json:"position"`
	PermissionOverwrites       []PermissionOverwrite `json:"permission_overwrites"`
	Name                       string                `json:"name"`
	Topic                      *string               `json:"topic"`
	NSFW                       bool                  `json:"nsfw"`
	LastMessageID              *snowflake.Snowflake  `json:"last_message_id"`
	RateLimitPerUser           int                   `json:"rate_limit_per_user"`
	ParentID                   *snowflake.Snowflake  `json:"parent_id"`
	LastPinTimestamp           *Time                 `json:"last_pin_timestamp"`
	DefaultAutoArchiveDuration AutoArchiveDuration   `json:"default_auto_archive_duration"`
}

func (t *guildTextChannel) UnmarshalJSON(data []byte) error {
	type guildTextChannelAlias guildTextChannel
	var v struct {
		PermissionOverwrites []UnmarshalPermissionOverwrite `json:"permission_overwrites"`
		guildTextChannelAlias
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*t = guildTextChannel(v.guildTextChannelAlias)
	if len(v.PermissionOverwrites) > 0 {
		t.PermissionOverwrites = make([]PermissionOverwrite, len(v.PermissionOverwrites))
		for i := range v.PermissionOverwrites {
			t.PermissionOverwrites[i] = v.PermissionOverwrites[i].PermissionOverwrite
		}
	}
	return nil
}

type guildNewsChannel struct {
	ID                         snowflake.Snowflake   `json:"id"`
	Type                       ChannelType           `json:"type"`
	GuildID                    snowflake.Snowflake   `json:"guild_id"`
	Position                   int                   `json:"position"`
	PermissionOverwrites       []PermissionOverwrite `json:"permission_overwrites"`
	Name                       string                `json:"name"`
	Topic                      *string               `json:"topic"`
	NSFW                       bool                  `json:"nsfw"`
	RateLimitPerUser           int                   `json:"rate_limit_per_user"`
	ParentID                   *snowflake.Snowflake  `json:"parent_id"`
	LastMessageID              *snowflake.Snowflake  `json:"last_message_id"`
	LastPinTimestamp           *Time                 `json:"last_pin_timestamp"`
	DefaultAutoArchiveDuration AutoArchiveDuration   `json:"default_auto_archive_duration"`
}

func (t *guildNewsChannel) UnmarshalJSON(data []byte) error {
	type guildNewsChannelAlias guildNewsChannel
	var v struct {
		PermissionOverwrites []UnmarshalPermissionOverwrite `json:"permission_overwrites"`
		guildNewsChannelAlias
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*t = guildNewsChannel(v.guildNewsChannelAlias)
	if len(v.PermissionOverwrites) > 0 {
		t.PermissionOverwrites = make([]PermissionOverwrite, len(v.PermissionOverwrites))
		for i := range v.PermissionOverwrites {
			t.PermissionOverwrites[i] = v.PermissionOverwrites[i].PermissionOverwrite
		}
	}
	return nil
}

type guildThread struct {
	ID               snowflake.Snowflake  `json:"id"`
	Type             ChannelType          `json:"type"`
	GuildID          snowflake.Snowflake  `json:"guild_id"`
	Name             string               `json:"name"`
	NSFW             bool                 `json:"nsfw"`
	LastMessageID    *snowflake.Snowflake `json:"last_message_id"`
	RateLimitPerUser int                  `json:"rate_limit_per_user"`
	OwnerID          snowflake.Snowflake  `json:"owner_id"`
	ParentID         snowflake.Snowflake  `json:"parent_id"`
	LastPinTimestamp *Time                `json:"last_pin_timestamp"`
	MessageCount     int                  `json:"message_count"`
	MemberCount      int                  `json:"member_count"`
	ThreadMetadata   ThreadMetadata       `json:"thread_metadata"`
}

type guildCategoryChannel struct {
	ID                   snowflake.Snowflake   `json:"id"`
	Type                 ChannelType           `json:"type"`
	GuildID              snowflake.Snowflake   `json:"guild_id"`
	Position             int                   `json:"position"`
	PermissionOverwrites []PermissionOverwrite `json:"permission_overwrites"`
	Name                 string                `json:"name"`
}

func (t *guildCategoryChannel) UnmarshalJSON(data []byte) error {
	type guildCategoryChannelAlias guildCategoryChannel
	var v struct {
		PermissionOverwrites []UnmarshalPermissionOverwrite `json:"permission_overwrites"`
		guildCategoryChannelAlias
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*t = guildCategoryChannel(v.guildCategoryChannelAlias)
	if len(v.PermissionOverwrites) > 0 {
		t.PermissionOverwrites = make([]PermissionOverwrite, len(v.PermissionOverwrites))
		for i := range v.PermissionOverwrites {
			t.PermissionOverwrites[i] = v.PermissionOverwrites[i].PermissionOverwrite
		}
	}
	return nil
}

type guildVoiceChannel struct {
	ID                   snowflake.Snowflake   `json:"id"`
	Type                 ChannelType           `json:"type"`
	GuildID              snowflake.Snowflake   `json:"guild_id"`
	Position             int                   `json:"position"`
	PermissionOverwrites []PermissionOverwrite `json:"permission_overwrites"`
	Name                 string                `json:"name"`
	Bitrate              int                   `json:"bitrate"`
	UserLimit            int                   `json:"user_limit"`
	ParentID             *snowflake.Snowflake  `json:"parent_id"`
	RTCRegion            string                `json:"rtc_region"`
	VideoQualityMode     VideoQualityMode      `json:"video_quality_mode"`
}

func (t *guildVoiceChannel) UnmarshalJSON(data []byte) error {
	type guildVoiceChannelAlias guildVoiceChannel
	var v struct {
		PermissionOverwrites []UnmarshalPermissionOverwrite `json:"permission_overwrites"`
		guildVoiceChannelAlias
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*t = guildVoiceChannel(v.guildVoiceChannelAlias)
	if len(v.PermissionOverwrites) > 0 {
		t.PermissionOverwrites = make([]PermissionOverwrite, len(v.PermissionOverwrites))
		for i := range v.PermissionOverwrites {
			t.PermissionOverwrites[i] = v.PermissionOverwrites[i].PermissionOverwrite
		}
	}
	return nil
}

type guildStageVoiceChannel struct {
	ID                   snowflake.Snowflake   `json:"id"`
	Type                 ChannelType           `json:"type"`
	GuildID              snowflake.Snowflake   `json:"guild_id"`
	Position             int                   `json:"position"`
	PermissionOverwrites []PermissionOverwrite `json:"permission_overwrites"`
	Name                 string                `json:"name"`
	Bitrate              int                   `json:"bitrate,"`
	ParentID             *snowflake.Snowflake  `json:"parent_id"`
	RTCRegion            string                `json:"rtc_region"`
}

func (t *guildStageVoiceChannel) UnmarshalJSON(data []byte) error {
	type guildStageVoiceChannelAlias guildStageVoiceChannel
	var v struct {
		PermissionOverwrites []UnmarshalPermissionOverwrite `json:"permission_overwrites"`
		guildStageVoiceChannelAlias
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*t = guildStageVoiceChannel(v.guildStageVoiceChannelAlias)
	if len(v.PermissionOverwrites) > 0 {
		t.PermissionOverwrites = make([]PermissionOverwrite, len(v.PermissionOverwrites))
		for i := range v.PermissionOverwrites {
			t.PermissionOverwrites[i] = v.PermissionOverwrites[i].PermissionOverwrite
		}
	}
	return nil
}
