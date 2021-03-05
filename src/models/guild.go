package models

import (
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
func (g Guild) IconURL(size int) *string {
	if g.Icon == nil {
		return nil
	}
	u :=endpoints.CDNGuildIcon(g.ID.String(), *g.Icon, size)
	return &u
}