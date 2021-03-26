package api

// Role is a Guild Role object
type Role struct {
	Disgo       Disgo
	GuildID     Snowflake
	ID          Snowflake   `json:"id"`
	Name        string      `json:"name"`
	Color       Color       `json:"color"`
	Hoist       bool        `json:"hoist"`
	Position    int         `json:"position"`
	Permissions Permissions `json:"permissions"`
	Managed     bool        `json:"managed"`
	Mentionable bool        `json:"mentionable"`
	Tags        []*RoleTag  `json:"tags,omitempty"`
}

func (r Role) Guild() *Guild {
	return r.Disgo.Cache().Guild(r.GuildID)
}

func (r Role) Update(roleUpdate UpdateRole) (*Role, error) {
	return r.Disgo.RestClient().UpdateRole(r.GuildID, r.ID, roleUpdate)
}

func (r Role) SetPosition(rolePositionUpdate UpdateRolePosition) ([]*Role, error) {
	return r.Disgo.RestClient().UpdateRolePositions(r.GuildID, rolePositionUpdate)
}

func (r Role) Delete() error {
	return r.Disgo.RestClient().DeleteRole(r.GuildID, r.ID)
}

type RoleTag struct {
	BotID             *Snowflake `json:"bot_id,omitempty"`
	IntegrationID     *Snowflake `json:"integration_id,omitempty"`
	PremiumSubscriber bool       `json:"premium_subscriber"`
}

type UpdateRole struct {
	Name        *string      `json:"name,omitempty"`
	Permissions *Permissions `json:"permissions,omitempty"`
	Color       *Color       `json:"color,omitempty"`
	Hoist       *bool        `json:"hoist,omitempty"`
	Mentionable *bool        `json:"mentionable,omitempty"`
}

type UpdateRolePosition struct {
	ID       Snowflake `json:"id"`
	Position *int      `json:"position"`
}
