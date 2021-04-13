package api

import (
	"time"

	"github.com/DisgoOrg/disgo/api/endpoints"
)

// ExpandedInvite is a full Invite struct
type ExpandedInvite struct {
	Invite
	Uses      int       `json:"uses"`
	MaxUses   int       `json:"max_uses"`
	MaxAge    int       `json:"max_age"`
	Temporary bool      `json:"temporary"`
	CreatedAt time.Time `json:"created_at"`
}

// Invite is a partial invite struct
type Invite struct {
	Disgo                    Disgo
	Code                     string        `json:"code"`
	Guild                    *InviteGuild  `json:"guild"`
	Channel                  InviteChannel `json:"channel"`
	Inviter                  *User         `json:"inviter"`
	TargetUser               *InviteUser   `json:"target_user"`
	TargetType               *TargetType   `json:"target_user_type"`
	ApproximatePresenceCount *int          `json:"approximate_presence_count"`
	ApproximateMemberCount   *int          `json:"approximate_member_count"`
}

func (i Invite) URL() string {
	url, err := endpoints.InviteURL.Compile(i.Code)
	if err != nil {
		return ""
	}
	return url.Route()
}

// TargetType is type of target an Invite uses
type TargetType int

// Constants for TargetType
const (
	TargetTypeStream TargetType = iota + 1
	TargetTypeEmbeddedApplication
)

// An InviteGuild is the Guild of an Invite
type InviteGuild struct {
	ID                Snowflake         `json:"id"`
	Name              string            `json:"name"`
	Splash            *string           `json:"splash"`
	Banner            *string           `json:"banner"`
	Description       *string           `json:"description"`
	Icon              *string           `json:"icon"`
	Features          []GuildFeature    `json:"features"`
	VerificationLevel VerificationLevel `json:"verification_level"`
	VanityURLCode     *string           `json:"vanity_url_code"`
}

// InviteChannel is the Channel of an invite
type InviteChannel struct {
	ID   string      `json:"id"`
	Name string      `json:"name"`
	Type ChannelType `json:"type"`
}

// InviteUser is the user who created an invite
type InviteUser struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Avatar        string `json:"avatar"`
	Discriminator string `json:"discriminator"`
}
