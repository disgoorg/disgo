package discord

import (
	"github.com/DisgoOrg/disgo/json"
)

var _ Mentionable = (*Role)(nil)

// Role is a Guild Role object
type Role struct {
	ID          Snowflake   `json:"id"`
	GuildID     Snowflake   `json:"guild_id"`
	Name        string      `json:"name"`
	Color       int         `json:"color"`
	Hoist       bool        `json:"hoist"`
	Position    int         `json:"position"`
	Permissions Permissions `json:"permissions"`
	Managed     bool        `json:"managed"`
	Icon        *string     `json:"icon"`
	Emoji       *string     `json:"unicode_emoji"`
	Mentionable bool        `json:"mentionable"`
	Tags        *RoleTag    `json:"tags,omitempty"`
}

func (r Role) String() string {
	return roleMention(r.ID)
}

func (r Role) Mention() string {
	return r.String()
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
	Icon        *Icon       `json:"icon,omitempty"`
	Emoji       *string     `json:"unicode_emoji,omitempty"`
	Mentionable bool        `json:"mentionable,omitempty"`
}

// RoleUpdate is the payload to update a Role
type RoleUpdate struct {
	Name        *string          `json:"name"`
	Permissions *Permissions     `json:"permissions"`
	Color       *int             `json:"color"`
	Hoist       *bool            `json:"hoist"`
	Icon        *NullIcon        `json:"icon,omitempty"`
	Emoji       *json.NullString `json:"unicode_emoji,omitempty"`
	Mentionable *bool            `json:"mentionable"`
}

// RolePositionUpdate is the payload to update a Role(s) position
type RolePositionUpdate struct {
	ID       Snowflake `json:"id"`
	Position *int      `json:"position"`
}

// PartialRole holds basic info about a Role
type PartialRole struct {
	ID   Snowflake `json:"id"`
	Name string    `json:"name"`
}
