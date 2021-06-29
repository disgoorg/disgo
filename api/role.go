package api

import "github.com/DisgoOrg/restclient"

// Role is a Guild Role object
type Role struct {
	Disgo       Disgo
	GuildID     Snowflake
	ID          Snowflake   `json:"id"`
	Name        string      `json:"name"`
	Color       int         `json:"color"`
	Hoist       bool        `json:"hoist"`
	Position    int         `json:"position"`
	Permissions Permissions `json:"permissions"`
	Managed     bool        `json:"managed"`
	Mentionable bool        `json:"mentionable"`
	Tags        *RoleTag    `json:"tags,omitempty"`
}

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
func (r *Role) Update(roleUpdate RoleUpdate) (*Role, restclient.RestError) {
	return r.Disgo.RestClient().UpdateRole(r.GuildID, r.ID, roleUpdate)
}

// SetPosition sets the position of the Role
func (r *Role) SetPosition(rolePositionUpdate RolePositionUpdate) ([]*Role, restclient.RestError) {
	return r.Disgo.RestClient().UpdateRolePositions(r.GuildID, rolePositionUpdate)
}

// Delete deletes the Role
func (r *Role) Delete() error {
	return r.Disgo.RestClient().DeleteRole(r.GuildID, r.ID)
}

// RoleTag are tags a Role has
type RoleTag struct {
	BotID             *Snowflake `json:"bot_id,omitempty"`
	IntegrationID     *Snowflake `json:"integration_id,omitempty"`
	PremiumSubscriber bool       `json:"premium_subscriber"`
}

// RoleCreate is the payload to create a Role
type RoleCreate struct {
	Name        string      `json:"name,omitempty"`
	Permissions Permissions `json:"permissions,omitempty"`
	Color       int         `json:"color,omitempty"`
	Hoist       bool        `json:"hoist,omitempty"`
	Mentionable bool        `json:"mentionable,omitempty"`
}

// RoleUpdate is the payload to update a Role
type RoleUpdate struct {
	Name        *string      `json:"name,omitempty"`
	Permissions *Permissions `json:"permissions,omitempty"`
	Color       *int         `json:"color,omitempty"`
	Hoist       *bool        `json:"hoist,omitempty"`
	Mentionable *bool        `json:"mentionable,omitempty"`
}

// RolePositionUpdate is the payload to update a Role(s) position
type RolePositionUpdate struct {
	ID       Snowflake `json:"id"`
	Position *int      `json:"position"`
}
