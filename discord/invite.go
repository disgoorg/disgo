package discord

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
	Code                     string            `json:"code"`
	Guild                    *InviteGuild      `json:"guild"`
	Channel                  InviteChannel     `json:"channel"`
	Inviter                  *User             `json:"inviter"`
	TargetUser               *InviteUser       `json:"target_user"`
	TargetType               *InviteTargetType `json:"target_user_type"`
	ApproximatePresenceCount *int              `json:"approximate_presence_count"`
	ApproximateMemberCount   *int              `json:"approximate_member_count"`
}

// ExpandedInvite is a full Invite struct
type ExpandedInvite struct {
	*Invite
	Uses      int  `json:"uses"`
	MaxUses   int  `json:"max_uses"`
	MaxAge    int  `json:"max_age"`
	Temporary bool `json:"temporary"`
	CreatedAt Time `json:"created_at"`
}

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

// InviteChannel is the Channel of an Invite
type InviteChannel struct {
	ID   string      `json:"id"`
	Name string      `json:"name"`
	Type ChannelType `json:"type"`
}

// InviteUser is the user who created an Invite
type InviteUser struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Avatar        string `json:"avatar"`
	Discriminator string `json:"discriminator"`
}

type CreateChannelInvite struct {
	MaxAgree            int              `json:"max_agree,omitempty"`
	MaxUses             int              `json:"max_uses,omitempty"`
	Temporary           bool             `json:"temporary,omitempty"`
	Unique              bool             `json:"unique,omitempty"`
	TargetType          InviteTargetType `json:"target_type,omitempty"`
	TargetUserID        Snowflake        `json:"target_user_id,omitempty"`
	TargetApplicationID Snowflake        `json:"target_application_id,omitempty"`
}
