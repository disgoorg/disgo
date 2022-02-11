package discord

import (
	"github.com/DisgoOrg/disgo/json"
	"github.com/DisgoOrg/snowflake"
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

type CChannelCreate struct {
	Name                 string                `json:"name"`
	Type                 ChannelType           `json:"type,omitempty"`
	Topic                string                `json:"topic,omitempty"`
	Bitrate              int                   `json:"bitrate,omitempty"`
	UserLimit            int                   `json:"user_limit,omitempty"`
	RateLimitPerUser     int                   `json:"rate_limit_per_user,omitempty"`
	Position             int                   `json:"position,omitempty"`
	PermissionOverwrites []PermissionOverwrite `json:"permission_overwrites,omitempty"`
	ParentID             snowflake.Snowflake   `json:"parent_id,omitempty"`
	NSFW                 bool                  `json:"nsfw,omitempty"`
}

type DMChannelCreate struct {
	RecipientID snowflake.Snowflake `json:"recipient_id"`
}
