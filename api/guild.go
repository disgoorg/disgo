package api

import (
	"strings"

	"github.com/DiscoOrg/disgo/api/endpoints"
)

// Guild represents a discord guild
type Guild struct {
	ID      Snowflake
	Name    string
	Icon    *string
	OwnerID Snowflake
}

// IconURL returns the Icon of a guild
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

// UnavailableGuild represents a unavailable discord guild
type UnavailableGuild struct {
	ID          Snowflake
	Unavailable bool
}
