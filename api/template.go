package api

import "time"

// GuildTemplate is a template used for copying guilds https://discord.com/developers/docs/resources/guild-template
type GuildTemplate struct {
	Code          string       `json:"code"`
	Name          string       `json:"name"`
	Description   *string      `json:"description,omitempty"`
	UsageCount    int          `json:"usage_count"`
	CreatorID     Snowflake    `json:"creator_id"`
	Creator       *User        `json:"creator"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
	SourceGuildID Snowflake    `json:"source_guild_id"`
	SourceGuild   PartialGuild `json:"serialized_source_guild"`
	IsDirty       *bool        `json:"is_dirty,omitempty"`
}

// CreateGuildTemplate is the data used to create a GuildTemplate
type CreateGuildTemplate struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// UpdateGuildTemplate is the data used to update a GuildTemplate
type UpdateGuildTemplate struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// CreateGuildFromTemplate is the data used to create a Guild from a GuildTemplate
type CreateGuildFromTemplate struct {
	Name      string `json:"name"`
	ImageData []byte `json:"icon,omitempty"`
}
