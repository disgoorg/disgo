package api

import (
	"strings"

	"github.com/DiscoOrg/disgo/api/endpoints"
)

// Guild represents a discord guild_events
type Guild struct {
	Disgo Disgo
	ID      Snowflake
	Name    string
	Icon    *string
	OwnerID Snowflake
}

// IconURL returns the Icon of a guild_events
func (g Guild) IconURL() *string {
	if g.Icon == nil {
		return nil
	}
	animated := strings.HasPrefix(*g.Icon, "a_")
	format := endpoints.PNG
	if animated {
		format = endpoints.GIF
	}
	u := endpoints.GuildIcon.Compile(format, g.ID.String(), *g.Icon)
	return &u
}

func (g Guild) CreateCommand(name string, description string) GuildCommandBuilder {
	return NewGuildCommandBuilder(g.Disgo, g.ID, name, description)
}

// UnavailableGuild represents a unavailable discord guild_events
type UnavailableGuild struct {
	ID          Snowflake
	Unavailable bool
}
