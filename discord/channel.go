package discord

// ChannelType for interacting with discord's channels
type ChannelType int

// Channel constants
//goland:noinspection GoUnusedConst
const (
	ChannelTypeText ChannelType = iota
	ChannelTypeDM
	ChannelTypeVoice
	ChannelTypeGroupDM
	ChannelTypeCategory
	ChannelTypeNews
	ChannelTypeStore
	ChannelTypeStage
)

// Channel is a generic discord channel object
type Channel struct {
	ID                     Snowflake             `json:"id"`
	Name                   *string               `json:"name,omitempty"`
	Type                   ChannelType           `json:"type"`
	LastMessageID          *Snowflake            `json:"last_message_id,omitempty"`
	GuildID                *Snowflake            `json:"guild_id,omitempty"`
	Position               *int                  `json:"position,omitempty"`
	PermissionOverwrites   []PermissionOverwrite `json:"permission_overwrites"`
	Topic                  *string               `json:"topic,omitempty"`
	NSFW                   *bool                 `json:"nsfw,omitempty"`
	Bitrate                *int                  `json:"bitrate,omitempty"`
	UserLimit              *int                  `json:"user_limit,omitempty"`
	RateLimitPerUser       *int                  `json:"rate_limit_per_user,omitempty"`
	Recipients             []*User               `json:"recipients,omitempty"`
	Icon                   *string               `json:"icon,omitempty"`
	OwnerID                *Snowflake            `json:"owner_id,omitempty"`
	ApplicationID          *Snowflake            `json:"application_id,omitempty"`
	ParentID               *Snowflake            `json:"parent_id,omitempty"`
	InteractionPermissions *Permissions          `json:"permissions,omitempty"`
	LastPinTimestamp       *Time                 `json:"last_pin_timestamp,omitempty"`
}

type ChannelCreate struct {
	Name                 string                `json:"name"`
	Type                 ChannelType           `json:"type,omitempty"`
	Topic                string                `json:"topic,omitempty"`
	Bitrate              int                   `json:"bitrate,omitempty"`
	UserLimit            int                   `json:"user_limit,omitempty"`
	RateLimitPerUser     int                   `json:"rate_limit_per_user,omitempty"`
	Position             int                   `json:"position,omitempty"`
	PermissionOverwrites []PermissionOverwrite `json:"permission_overwrites,omitempty"`
	ParentID             Snowflake             `json:"parent_id,omitempty"`
	NSFW                 bool                  `json:"nsfw,omitempty"`
}

// PartialChannel contains basic info about a Channel
type PartialChannel struct {
	ID   Snowflake   `json:"id"`
	Type ChannelType `json:"type"`
	Name string      `json:"name"`
}
