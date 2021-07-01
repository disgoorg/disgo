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
func (r *Role) Update(roleUpdate UpdateRole) (*Role, restclient.RestError) {
	return r.Disgo.RestClient().UpdateRole(r.GuildID, r.ID, roleUpdate)
}

// SetPosition sets the position of the Role
func (r *Role) SetPosition(rolePositionUpdate UpdateRolePosition) ([]*Role, restclient.RestError) {
	return r.Disgo.RestClient().UpdateRolePositions(r.GuildID, rolePositionUpdate)
}

// Delete deletes the Role
func (r *Role) Delete() restclient.RestError {
	return r.Disgo.RestClient().DeleteRole(r.GuildID, r.ID)
}

// RoleTag are tags a Role has
type RoleTag struct {
	BotID             *Snowflake `json:"bot_id,omitempty"`
	IntegrationID     *Snowflake `json:"integration_id,omitempty"`
	PremiumSubscriber bool       `json:"premium_subscriber"`
}

// CreateRole is the payload to create a Role
type CreateRole struct {
	Name        string      `json:"name,omitempty"`
	Permissions Permissions `json:"permissions,omitempty"`
	Color       int         `json:"color,omitempty"`
	Hoist       bool        `json:"hoist,omitempty"`
	Mentionable bool        `json:"mentionable,omitempty"`
}

// UpdateRole is the payload to update a Role
type UpdateRole struct {
	Name        *string      `json:"name"`
	Permissions *Permissions `json:"permissions"`
	Color       *int         `json:"color"`
	Hoist       *bool        `json:"hoist"`
	Mentionable *bool        `json:"mentionable"`
}

// UpdateRolePosition is the payload to update a Role(s) position
type UpdateRolePosition struct {
	ID       Snowflake `json:"id"`
	Position *int      `json:"position"`
}
