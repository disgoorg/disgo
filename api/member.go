package api

import (
	"time"

	"github.com/DisgoOrg/restclient"
)

// Member is a discord GuildMember
type Member struct {
	Disgo              Disgo
	GuildID            Snowflake    `json:"guild_id"`
	User               *User        `json:"user"`
	Nick               *string      `json:"nick"`
	RoleIDs            []Snowflake  `json:"roles,omitempty"`
	JoinedAt           time.Time    `json:"joined_at"`
	PremiumSince       *time.Time   `json:"premium_since,omitempty"`
	Deaf               *bool        `json:"deaf,omitempty"`
	Mute               *bool        `json:"mute,omitempty"`
	Pending            bool         `json:"pending"`
	ChannelPermissions *Permissions `json:"permissions,omitempty"`
}

// Permissions returns the Permissions the Member has in the Guild
func (m *Member) Permissions() Permissions {
	return GetMemberPermissions(m)
}

// Roles return all Role(s)the Member has
func (m *Member) Roles() []*Role {
	var roles []*Role
	allRoles := m.Disgo.Cache().RoleCache(m.GuildID)
	for _, roleID := range m.RoleIDs {
		roles = append(roles, allRoles[roleID])
	}
	return roles
}

// VoiceState returns the VoiceState for this Member from the Cache(requires CacheFlagVoiceState and GatewayIntentsGuildVoiceStates)
func (m *Member) VoiceState() *VoiceState {
	return m.Disgo.Cache().VoiceState(m.GuildID, m.User.ID)
}

// EffectiveName returns either the nickname or username depending on if the user has a nickname
func (m *Member) EffectiveName() string {
	if m.Nick != nil {
		return *m.Nick
	}
	return m.User.Username
}

// Guild returns the members guild from the cache
func (m *Member) Guild() *Guild {
	return m.Disgo.Cache().Guild(m.GuildID)
}

// IsOwner returns whether the member is the owner of the guild_events that it belongs to
func (m *Member) IsOwner() bool {
	if guild := m.Guild(); guild != nil {
		return guild.OwnerID == m.User.ID
	}
	return false
}

// Update updates the Member
func (m *Member) Update(memberUpdate MemberUpdate) (*Member, restclient.RestError) {
	return m.Disgo.RestClient().UpdateMember(m.GuildID, m.User.ID, memberUpdate)
}

// Kick kicks the Member from the Guild
func (m *Member) Kick(reason string) restclient.RestError {
	return m.Disgo.RestClient().RemoveMember(m.GuildID, m.User.ID, reason)
}

// Move moves/kicks the member to/from a voice channel
func (m *Member) Move(channelID *Snowflake) (*Member, restclient.RestError) {
	return m.Disgo.RestClient().MoveMember(m.GuildID, m.User.ID, channelID)
}

// AddRole adds a specific role the member
func (m *Member) AddRole(roleID Snowflake) restclient.RestError {
	return m.Disgo.RestClient().AddMemberRole(m.GuildID, m.User.ID, roleID)
}

// RemoveRole removes a specific role the member
func (m *Member) RemoveRole(roleID Snowflake) restclient.RestError {
	return m.Disgo.RestClient().AddMemberRole(m.GuildID, m.User.ID, roleID)
}

// MemberAdd is used to add a Member via the oauth2 access token to a Guild
type MemberAdd struct {
	AccessToken string      `json:"access_token"`
	Nick        string      `json:"nick,omitempty"`
	Roles       []Snowflake `json:"roles,omitempty"`
	Mute        bool        `json:"mute,omitempty"`
	Deaf        bool        `json:"deaf,omitempty"`
}

// MemberUpdate is used to modify a Member
type MemberUpdate struct {
	*MemberMove
	Nick  *string     `json:"nick,omitempty"`
	Roles []Snowflake `json:"roles,omitempty"`
	Mute  *bool       `json:"mute,omitempty"`
	Deaf  *bool       `json:"deaf,omitempty"`
}

// MemberMove is used to move a member
type MemberMove struct {
	ChannelID *Snowflake `json:"channel_id,omitempty"`
}

// UpdateSelfNick is used to update your own nick
type UpdateSelfNick struct {
	Nick string `json:"nick"`
}
