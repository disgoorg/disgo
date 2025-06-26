package discord

import (
	"time"

	"github.com/disgoorg/omit"
	"github.com/disgoorg/snowflake/v2"
)

var _ Mentionable = (*Role)(nil)

// Role is a Guild Role object
type Role struct {
	ID          snowflake.ID `json:"id"`
	GuildID     snowflake.ID `json:"guild_id,omitempty"` // not present in the API but we need it
	Name        string       `json:"name"`
	Description *string      `json:"description,omitempty"`
	Color       int          `json:"color"`
	RoleColors  RoleColors   `json:"colors"`
	Hoist       bool         `json:"hoist"`
	Position    int          `json:"position"`
	Permissions Permissions  `json:"permissions"`
	Managed     bool         `json:"managed"`
	Icon        *string      `json:"icon"`
	Emoji       *string      `json:"unicode_emoji"`
	Mentionable bool         `json:"mentionable"`
	Tags        *RoleTag     `json:"tags,omitempty"`
	Flags       RoleFlags    `json:"flags"`
}

func (r Role) String() string {
	return RoleMention(r.ID)
}

func (r Role) Mention() string {
	return r.String()
}

func (r Role) IconURL(opts ...CDNOpt) *string {
	if r.Icon == nil {
		return nil
	}
	url := formatAssetURL(RoleIcon, opts, r.ID, *r.Icon)
	return &url
}

func (r Role) CreatedAt() time.Time {
	return r.ID.Time()
}

type RoleColors struct {
	PrimaryColor   int  `json:"primary_color"`
	SecondaryColor *int `json:"secondary_color"`
	TertiaryColor  *int `json:"tertiary_color"`
}

// RoleTag are tags a Role has
type RoleTag struct {
	BotID                 *snowflake.ID `json:"bot_id,omitempty"`
	IntegrationID         *snowflake.ID `json:"integration_id,omitempty"`
	PremiumSubscriber     bool          `json:"premium_subscriber"`
	SubscriptionListingID *snowflake.ID `json:"subscription_listing_id,omitempty"`
	AvailableForPurchase  bool          `json:"available_for_purchase"`
	GuildConnections      bool          `json:"guild_connections"`
}

type RoleFlags int

const (
	RoleFlagInPrompt RoleFlags = 1 << iota
	RoleFlagsNone    RoleFlags = 0
)

// RoleCreate is the payload to create a Role
type RoleCreate struct {
	Name        string       `json:"name,omitempty"`
	Permissions *Permissions `json:"permissions,omitempty"`
	Color       int          `json:"color,omitempty"`
	Colors      RoleColors   `json:"colors,omitempty"`
	Hoist       bool         `json:"hoist,omitempty"`
	Icon        *Icon        `json:"icon,omitempty"`
	Emoji       string       `json:"unicode_emoji,omitempty"`
	Mentionable bool         `json:"mentionable,omitempty"`
}

// RoleUpdate is the payload to update a Role
type RoleUpdate struct {
	Name        *string                `json:"name,omitempty"`
	Permissions *Permissions           `json:"permissions,omitempty"`
	Color       *int                   `json:"color,omitempty"`
	Colors      omit.Omit[*RoleColors] `json:"colors,omitzero"`
	Hoist       *bool                  `json:"hoist,omitempty"`
	Icon        omit.Omit[*Icon]       `json:"icon,omitzero"`
	Emoji       *string                `json:"unicode_emoji,omitempty"`
	Mentionable *bool                  `json:"mentionable,omitempty"`
}

// RolePositionUpdate is the payload to update a Role(s) position
type RolePositionUpdate struct {
	ID       snowflake.ID `json:"id"`
	Position *int         `json:"position,omitempty"`
}

// PartialRole holds basic info about a Role
type PartialRole struct {
	ID   snowflake.ID `json:"id"`
	Name string       `json:"name"`
}
