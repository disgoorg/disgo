package discord

import (
	"time"

	"github.com/disgoorg/json"
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/internal/flags"
)

var _ Mentionable = (*Member)(nil)

// Member is a discord GuildMember
type Member struct {
	User                       User                  `json:"user"`
	Nick                       *string               `json:"nick"`
	Avatar                     *string               `json:"avatar"`
	RoleIDs                    []snowflake.ID        `json:"roles,omitempty"`
	JoinedAt                   time.Time             `json:"joined_at"`
	PremiumSince               *time.Time            `json:"premium_since,omitempty"`
	Deaf                       bool                  `json:"deaf,omitempty"`
	Mute                       bool                  `json:"mute,omitempty"`
	Flags                      MemberFlags           `json:"flags"`
	Pending                    bool                  `json:"pending"`
	CommunicationDisabledUntil *time.Time            `json:"communication_disabled_until"`
	AvatarDecorationData       *AvatarDecorationData `json:"avatar_decoration_data"`

	// This field is not present everywhere in the API and often populated by disgo
	GuildID snowflake.ID `json:"guild_id"`
}

// String returns a mention of the user
func (m Member) String() string {
	return m.User.String()
}

// Mention returns a mention of the user
func (m Member) Mention() string {
	return m.String()
}

// EffectiveName returns the nickname of the member if set, falling back to User.EffectiveName()
func (m Member) EffectiveName() string {
	if m.Nick != nil {
		return *m.Nick
	}
	return m.User.EffectiveName()
}

// EffectiveAvatarURL returns the guild-specific avatar URL of the user if set, falling back to the effective avatar URL of the user
func (m Member) EffectiveAvatarURL(opts ...CDNOpt) string {
	if m.Avatar == nil {
		return m.User.EffectiveAvatarURL(opts...)
	}
	if avatar := m.AvatarURL(opts...); avatar != nil {
		return *avatar
	}
	return ""
}

// AvatarURL returns the guild-specific avatar URL of the user if set or nil
func (m Member) AvatarURL(opts ...CDNOpt) *string {
	if m.Avatar == nil {
		return nil
	}
	url := formatAssetURL(MemberAvatar, opts, m.GuildID, m.User.ID, *m.Avatar)
	return &url
}

// AvatarDecorationURL returns the avatar decoration URL if set or nil
func (m Member) AvatarDecorationURL(opts ...CDNOpt) *string {
	if m.AvatarDecorationData == nil {
		return nil
	}
	url := formatAssetURL(AvatarDecoration, opts, m.AvatarDecorationData.Asset)
	return &url
}

func (m Member) CreatedAt() time.Time {
	return m.User.CreatedAt()
}

// MemberAdd is used to add a member via the oauth2 access token to a guild
type MemberAdd struct {
	AccessToken string         `json:"access_token"`
	Nick        string         `json:"nick,omitempty"`
	Roles       []snowflake.ID `json:"roles,omitempty"`
	Mute        bool           `json:"mute,omitempty"`
	Deaf        bool           `json:"deaf,omitempty"`
}

// MemberUpdate is used to modify a member
type MemberUpdate struct {
	ChannelID                  *snowflake.ID             `json:"channel_id,omitempty"`
	Nick                       *string                   `json:"nick,omitempty"`
	Roles                      *[]snowflake.ID           `json:"roles,omitempty"`
	Mute                       *bool                     `json:"mute,omitempty"`
	Deaf                       *bool                     `json:"deaf,omitempty"`
	Flags                      *MemberFlags              `json:"flags,omitempty"`
	CommunicationDisabledUntil *json.Nullable[time.Time] `json:"communication_disabled_until,omitempty"`
}

// CurrentMemberUpdate is used to update the current member
type CurrentMemberUpdate struct {
	Nick string `json:"nick"`
}

type MemberFlags int

const (
	MemberFlagDidRejoin MemberFlags = 1 << iota
	MemberFlagCompletedOnboarding
	MemberFlagBypassesVerification
	MemberFlagStartedOnboarding
	MemberFlagsNone MemberFlags = 0
)

// Add allows you to add multiple bits together, producing a new bit
func (f MemberFlags) Add(bits ...MemberFlags) MemberFlags {
	return flags.Add(f, bits...)
}

// Remove allows you to subtract multiple bits from the first, producing a new bit
func (f MemberFlags) Remove(bits ...MemberFlags) MemberFlags {
	return flags.Remove(f, bits...)
}

// Has will ensure that the bit includes all the bits entered
func (f MemberFlags) Has(bits ...MemberFlags) bool {
	return flags.Has(f, bits...)
}

// Missing will check whether the bit is missing any one of the bits
func (f MemberFlags) Missing(bits ...MemberFlags) bool {
	return flags.Missing(f, bits...)
}
