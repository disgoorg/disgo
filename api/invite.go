package api

import "time"

type ExpandedInvite struct {
	Invite
	Uses      int       `json:"uses"`
	MaxUses   int       `json:"max_uses"`
	MaxAge    int       `json:"max_age"`
	Temporary bool      `json:"temporary"`
	CreatedAt time.Time `json:"created_at"`
}

type Invite struct {
	Disgo                    Disgo
	Code                     string          `json:"code"`
	Guild                    *InviteGuild    `json:"guild"`
	Channel                  InviteChannel   `json:"channel"`
	Inviter                  *User           `json:"inviter"`
	TargetUser               *InviteUser     `json:"target_user"`
	TargetUserType           *TargetUserType `json:"target_user_type"`
	ApproximatePresenceCount *int            `json:"approximate_presence_count"`
	ApproximateMemberCount   *int            `json:"approximate_member_count"`
}

type TargetUserType int

const (
	TargetUserTypeStream = iota + 1
)

type InviteGuild struct {
	Id                Snowflake         `json:"id"`
	Name              string            `json:"name"`
	Splash            *string           `json:"splash"`
	Banner            *string           `json:"banner"`
	Description       *string           `json:"description"`
	Icon              *string           `json:"icon"`
	Features          []GuildFeature    `json:"features"`
	VerificationLevel VerificationLevel `json:"verification_level"`
	VanityUrlCode     *string           `json:"vanity_url_code"`
}

type InviteChannel struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Type int    `json:"type"`
}

type InviteUser struct {
	Id            string `json:"id"`
	Username      string `json:"username"`
	Avatar        string `json:"avatar"`
	Discriminator string `json:"discriminator"`
}
