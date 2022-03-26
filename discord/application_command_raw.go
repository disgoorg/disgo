package discord

import (
	"github.com/disgoorg/disgo/json"
	"github.com/disgoorg/snowflake"
)

type rawSlashCommand struct {
	ID                       snowflake.Snowflake        `json:"id"`
	Type                     ApplicationCommandType     `json:"type"`
	ApplicationID            snowflake.Snowflake        `json:"application_id"`
	GuildID                  *snowflake.Snowflake       `json:"guild_id,omitempty"`
	Name                     string                     `json:"name"`
	NameLocalizations        map[Locale]string          `json:"name_localizations,omitempty"`
	NameLocalized            string                     `json:"name_localized,omitempty"`
	Description              string                     `json:"description,omitempty"`
	DescriptionLocalizations map[Locale]string          `json:"description_localizations,omitempty"`
	DescriptionLocalized     string                     `json:"description_localized,omitempty"`
	Options                  []ApplicationCommandOption `json:"options,omitempty"`
	DefaultPermission        bool                       `json:"default_permission,omitempty"`
	Version                  snowflake.Snowflake        `json:"version"`
}

func (c *rawSlashCommand) UnmarshalJSON(data []byte) error {
	type alias rawSlashCommand
	var sc struct {
		Options []UnmarshalApplicationCommandOption `json:"options,omitempty"`
		alias
	}

	if err := json.Unmarshal(data, &sc); err != nil {
		return err
	}

	*c = rawSlashCommand(sc.alias)

	if len(sc.Options) > 0 {
		c.Options = make([]ApplicationCommandOption, len(sc.Options))
		for i := range sc.Options {
			c.Options[i] = sc.Options[i].ApplicationCommandOption
		}
	}
	return nil
}

type rawContextCommand struct {
	ID                snowflake.Snowflake    `json:"id"`
	Type              ApplicationCommandType `json:"type"`
	ApplicationID     snowflake.Snowflake    `json:"application_id"`
	GuildID           *snowflake.Snowflake   `json:"guild_id,omitempty"`
	Name              string                 `json:"name"`
	NameLocalizations map[Locale]string      `json:"name_localizations,omitempty"`
	NameLocalized     string                 `json:"name_localized,omitempty"`
	DefaultPermission bool                   `json:"default_permission,omitempty"`
	Version           snowflake.Snowflake    `json:"version"`
}
