package discord

type DMChannel struct {
	ID               Snowflake  `json:"id"`
	Name             string     `json:"name,omitempty"`
	LastMessageID    *Snowflake `json:"last_message_id,omitempty"`
	Recipients       []User     `json:"recipients,omitempty"`
	Icon             *string    `json:"icon,omitempty"`
	OwnerID          Snowflake  `json:"owner_id,omitempty"`
	ApplicationID    Snowflake  `json:"application_id,omitempty"`
	LastPinTimestamp *Time      `json:"last_pin_timestamp,omitempty"`
}

type GuildTextChannel struct {
	ID                         Snowflake             `json:"id"`
	GuildID                    Snowflake             `json:"guild_id,omitempty"`
	Position                   int                   `json:"position,omitempty"`
	PermissionOverwrites       []PermissionOverwrite `json:"permission_overwrites"`
	Name                       string                `json:"name,omitempty"`
	Topic                      *string               `json:"topic,omitempty"`
	NSFW                       bool                  `json:"nsfw,omitempty"`
	LastMessageID              *Snowflake            `json:"last_message_id,omitempty"`
	RateLimitPerUser           int                   `json:"rate_limit_per_user,omitempty"`
	ParentID                   *Snowflake            `json:"parent_id,omitempty"`
	LastPinTimestamp           *Time                 `json:"last_pin_timestamp,omitempty"`
	DefaultAutoArchiveDuration AutoArchiveDuration   `json:"default_auto_archive_duration"`
	InteractionPermissions     Permissions           `json:"permissions,omitempty"`
}

type GuildNewsChannel struct {
	ID                         Snowflake             `json:"id"`
	GuildID                    Snowflake             `json:"guild_id,omitempty"`
	Position                   int                   `json:"position,omitempty"`
	PermissionOverwrites       []PermissionOverwrite `json:"permission_overwrites"`
	Name                       string                `json:"name,omitempty"`
	Topic                      *string               `json:"topic,omitempty"`
	NSFW                       bool                  `json:"nsfw,omitempty"`
	LastMessageID              *Snowflake            `json:"last_message_id,omitempty"`
	RateLimitPerUser           int                   `json:"rate_limit_per_user,omitempty"`
	ParentID                   *Snowflake            `json:"parent_id,omitempty"`
	LastPinTimestamp           *Time                 `json:"last_pin_timestamp,omitempty"`
	DefaultAutoArchiveDuration AutoArchiveDuration   `json:"default_auto_archive_duration"`
	InteractionPermissions     Permissions           `json:"permissions,omitempty"`
}

type AutoArchiveDuration int

const (
	AutoArchiveDuration1h  AutoArchiveDuration = 60
	AutoArchiveDuration24h AutoArchiveDuration = 1440
	AutoArchiveDuration3d  AutoArchiveDuration = 4320
	AutoArchiveDuration1w  AutoArchiveDuration = 10080
)
