package api

import (
	"time"

	"github.com/DisgoOrg/restclient"
)

// GuildTemplate is a template used for copying guilds https://discord.com/developers/docs/resources/guild-template
type GuildTemplate struct {
	Disgo        Disgo
	Code         string       `json:"code"`
	Name         string       `json:"name"`
	Description  *string      `json:"description,omitempty"`
	UsageCount   int          `json:"usage_count"`
	CreatorID    Snowflake    `json:"creator_id"`
	Creator      *User        `json:"creator"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	GuildID      Snowflake    `json:"source_guild_id"`
	PartialGuild PartialGuild `json:"serialized_source_guild"`
	IsDirty       bool        `json:"is_dirty,omitempty"`
}

// Guild returns the full Guild of the GuildTemplate if in cache
func (t *GuildTemplate) Guild() *Guild {
	return t.Disgo.Cache().Guild(t.GuildID)
}

// Update updates the GuildTemplate with the provided UpdateGuildTemplate
func (t *GuildTemplate) Update(updateGuildTemplate UpdateGuildTemplate) (*GuildTemplate, restclient.RestError) {
	return t.Disgo.RestClient().UpdateGuildTemplate(t.GuildID, t.Code, updateGuildTemplate)
}

// Delete deletes the GuildTemplate
func (t *GuildTemplate) Delete() (*GuildTemplate, restclient.RestError) {
	return t.Disgo.RestClient().DeleteGuildTemplate(t.GuildID, t.Code)
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
