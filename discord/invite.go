package discord

import (
	"github.com/DisgoOrg/disgo/rest/route"
	"github.com/DisgoOrg/snowflake"
)

// InviteTargetType is type of target an Invite uses
type InviteTargetType int

// Constants for TargetType
//goland:noinspection GoUnusedConst
const (
	InviteTargetTypeStream InviteTargetType = iota + 1
	InviteTargetTypeEmbeddedApplication
)

// Invite is a partial invite struct
type Invite struct {
	Code                     string               `json:"code"`
	Guild                    *InviteGuild         `json:"guild"`
	Channel                  *InviteChannel       `json:"channel"`
	ChannelID                snowflake.Snowflake  `json:"channel_id"`
	Inviter                  *User                `json:"inviter"`
	TargetUser               *User                `json:"target_user"`
	TargetType               InviteTargetType     `json:"target_user_type"`
	ApproximatePresenceCount int                  `json:"approximate_presence_count"`
	ApproximateMemberCount   int                  `json:"approximate_member_count"`
	ExpiresAt                *Time                `json:"created_at"`
	GuildScheduledEvent      *GuildScheduledEvent `json:"guild_scheduled_event"`
}

type ExtendedInvite struct {
	Invite
	Uses      int  `json:"uses"`
	MaxUses   int  `json:"max_uses"`
	MaxAge    int  `json:"max_age"`
	Temporary bool `json:"temporary"`
	CreatedAt Time `json:"created_at"`
}

type InviteChannel struct {
	ID   snowflake.Snowflake `json:"id"`
	Type ChannelType         `json:"type"`
	Name string              `json:"name"`
	Icon *string             `json:"icon,omitempty"`
}

// GetIconURL returns the Icon URL of this channel.
// This will be nil for every ChannelType except ChannelTypeGroupDM
func (c InviteChannel) GetIconURL(size int) *string {
	return FormatAssetURL(route.ChannelIcon, c.ID, c.Icon, size)
}

// An InviteGuild is the Guild of an Invite
type InviteGuild struct {
	ID                snowflake.Snowflake `json:"id"`
	Name              string              `json:"name"`
	Splash            *string             `json:"splash"`
	Banner            *string             `json:"banner"`
	Description       *string             `json:"description"`
	Icon              *string             `json:"icon"`
	Features          []GuildFeature      `json:"features"`
	VerificationLevel VerificationLevel   `json:"verification_level"`
	VanityURLCode     *string             `json:"vanity_url_code"`
}

type InviteCreate struct {
	MaxAgree            int                 `json:"max_agree,omitempty"`
	MaxUses             int                 `json:"max_uses,omitempty"`
	Temporary           bool                `json:"temporary,omitempty"`
	Unique              bool                `json:"unique,omitempty"`
	TargetType          InviteTargetType    `json:"target_type,omitempty"`
	TargetUserID        snowflake.Snowflake `json:"target_user_id,omitempty"`
	TargetApplicationID snowflake.Snowflake `json:"target_application_id,omitempty"`
}
