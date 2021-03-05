package models

import (
	"strings"

	"github.com/DiscoOrg/disgo/src/endpoints"
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
	format := "png"
	if animated {
		format = "gif"
	}
	u := endpoints.GuildIcon.Compile(g.ID.String(), *g.Icon, format)
	return &u
}
