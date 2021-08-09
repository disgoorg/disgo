package entities

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/restclient"
)



// Mention parses the Role as a Mention
func (r *Role) Mention() string {
	return "<@&" + r.ID.String() + ">"
}

// String parses the Role to a String representation
func (r *Role) String() string {
	return r.Mention()
}

// Guild returns the Guild of this role from the Cache
func (r *Role) Guild() *Guild {
	return r.Disgo.Cache().Guild(r.GuildID)
}

// Update updates the Role with specific values
func (r *Role) Update(roleUpdate UpdateRole) (*Role, rest.Error) {
	return r.Disgo.RestServices().UpdateRole(r.GuildID, r.ID, roleUpdate)
}

// SetPosition sets the position of the Role
func (r *Role) SetPosition(rolePositionUpdate UpdateRolePosition) ([]*Role, rest.Error) {
	return r.Disgo.RestServices().UpdateRolePositions(r.GuildID, rolePositionUpdate)
}

// Delete deletes the Role
func (r *Role) Delete() rest.Error {
	return r.Disgo.RestServices().DeleteRole(r.GuildID, r.ID)
}

// RoleTag are tags a Role has
type RoleTag struct {
	BotID             *discord.Snowflake `json:"bot_id,omitempty"`
	IntegrationID     *discord.Snowflake `json:"integration_id,omitempty"`
	PremiumSubscriber bool               `json:"premium_subscriber"`
}


