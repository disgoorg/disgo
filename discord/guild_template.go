package discord

// GuildTemplate is a template used for copying guilds https://discord.com/developers/docs/resources/guild-template
type GuildTemplate struct {
	Code         string       `json:"code"`
	Name         string       `json:"name"`
	Description  *string      `json:"description,omitempty"`
	UsageCount   int          `json:"usage_count"`
	CreatorID    Snowflake    `json:"creator_id"`
	Creator      User         `json:"creator"`
	CreatedAt    Time         `json:"created_at"`
	UpdatedAt    Time         `json:"updated_at"`
	GuildID      Snowflake    `json:"source_guild_id"`
	PartialGuild PartialGuild `json:"serialized_source_guild"`
	IsDirty      bool         `json:"is_dirty,omitempty"`
}

// GuildTemplateCreate is the data used to create a GuildTemplate
type GuildTemplateCreate struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// GuildTemplateUpdate is the data used to update a GuildTemplate
type GuildTemplateUpdate struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// GuildFromTemplateCreate is the data used to create a Guild from a GuildTemplate
type GuildFromTemplateCreate struct {
	Name      string `json:"name"`
	ImageData []byte `json:"icon,omitempty"`
}
