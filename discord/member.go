package discord

import (
	"time"

	"github.com/disgoorg/disgo/json"
	"github.com/disgoorg/disgo/rest/route"
	"github.com/disgoorg/snowflake"
)

var _ Mentionable = (*Member)(nil)

// Member is a discord GuildMember
type Member struct {
	User                       User                  `json:"user"`
	Nick                       *string               `json:"nick"`
	Avatar                     *string               `json:"avatar"`
	RoleIDs                    []snowflake.Snowflake `json:"roles,omitempty"`
	JoinedAt                   time.Time             `json:"joined_at"`
	PremiumSince               *time.Time            `json:"premium_since,omitempty"`
	Deaf                       bool                  `json:"deaf,omitempty"`
	Mute                       bool                  `json:"mute,omitempty"`
	Pending                    bool                  `json:"pending"`
	CommunicationDisabledUntil *time.Time            `json:"communication_disabled_until"`
}

func (m Member) String() string {
	return m.User.String()
}

func (m Member) Mention() string {
	return m.String()
}

// EffectiveName returns either the nickname or username depending on if the user has a nickname
func (m Member) EffectiveName() string {
	if m.Nick != nil {
		return *m.Nick
	}
	return m.User.Username
}

func (m Member) EffectiveAvatarURL(opts ...CDNOpt) string {
	if m.Avatar == nil {
		return m.User.EffectiveAvatarURL(opts...)
	}
	if avatar := m.AvatarURL(opts...); avatar != nil {
		return *avatar
	}
	return ""
}

func (m Member) AvatarURL(opts ...CDNOpt) *string {
	if m.Avatar == nil {
		return nil
	}
	return formatAssetURL(route.MemberAvatar, opts, m.User.ID, *m.Avatar)
}

// MemberAdd is used to add a member via the oauth2 access token to a guild
type MemberAdd struct {
	AccessToken string                `json:"access_token"`
	Nick        string                `json:"nick,omitempty"`
	Roles       []snowflake.Snowflake `json:"roles,omitempty"`
	Mute        bool                  `json:"mute,omitempty"`
	Deaf        bool                  `json:"deaf,omitempty"`
}

// MemberUpdate is used to modify
type MemberUpdate struct {
	ChannelID                  *snowflake.Snowflake      `json:"channel_id,omitempty"`
	Nick                       *string                   `json:"nick,omitempty"`
	Roles                      []snowflake.Snowflake     `json:"roles,omitempty"`
	Mute                       *bool                     `json:"mute,omitempty"`
	Deaf                       *bool                     `json:"deaf,omitempty"`
	CommunicationDisabledUntil *json.Nullable[time.Time] `json:"communication_disabled_until,omitempty"`
}

// SelfNickUpdate is used to update your own nick
type SelfNickUpdate struct {
	Nick string `json:"nick"`
}
