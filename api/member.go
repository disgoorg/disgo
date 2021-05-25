package api

import "time"

// Member is a discord GuildMember
type Member struct {
	Disgo        Disgo
	GuildID      Snowflake    `json:"guild_id"`
	User         *User        `json:"user"`
	Nick         *string      `json:"nick"`
	Roles        []Snowflake  `json:"roles,omitempty"`
	JoinedAt     time.Time    `json:"joined_at"`
	PremiumSince *time.Time   `json:"premium_since,omitempty"`
	Deaf         *bool        `json:"deaf,omitempty"`
	Mute         *bool        `json:"mute,omitempty"`
	Pending      bool         `json:"pending"`
	Permissions  *Permissions `json:"permissions,omitempty"`
}

// VoiceState returns the VoiceState for this Member from the Cache(requires CacheFlagVoiceState and GatewayIntentsGuildVoiceStates)
func (m Member) VoiceState() *VoiceState {
	return m.Disgo.Cache().VoiceState(m.GuildID, m.User.ID)
}

// EffectiveName returns either the nickname or username depending on if the user has a nickname
func (m Member) EffectiveName() string {
	if m.Nick != nil {
		return *m.Nick
	}
	return m.User.Username
}

// Guild returns the members guild from the cache
func (m Member) Guild() *Guild {
	return m.Disgo.Cache().Guild(m.GuildID)
}

// IsOwner returns whether the member is the owner of the guild_events that it belongs to
func (m Member) IsOwner() bool {
	return m.Guild().OwnerID == m.User.ID
}

// Update updates the member
func (m Member) Update(updateGuildMemberData *UpdateGuildMemberData) (*Member, error) {
	return m.Disgo.RestClient().UpdateMember(m.GuildID, m.User.ID, updateGuildMemberData)
}

// Move moves/kicks the member to/from a voice channel
func (m Member) Move(channelID *Snowflake) (*Member, error) {
	return m.Disgo.RestClient().MoveMember(m.GuildID, m.User.ID, channelID)
}

// AddRole adds a specific role the member
func (m Member) AddRole(roleID Snowflake) error {
	return m.Disgo.RestClient().AddMemberRole(m.GuildID, m.User.ID, roleID)
}

// RemoveRole removes a specific role the member
func (m Member) RemoveRole(roleID Snowflake) error {
	return m.Disgo.RestClient().AddMemberRole(m.GuildID, m.User.ID, roleID)
}

// AddGuildMemberData is used to add a member via the oauth2 access token to a guild
type AddGuildMemberData struct {
	AccessToken string      `json:"access_token"`
	Nick        *string     `json:"nick,omitempty"`
	Roles       []Snowflake `json:"roles,omitempty"`
	Mute        *bool       `json:"mute,omitempty"`
	Deaf        *bool       `json:"deaf,omitempty"`
}

// UpdateGuildMemberData is used to modify
type UpdateGuildMemberData struct {
	Nick  *string     `json:"nick,omitempty"`
	Roles []Snowflake `json:"roles,omitempty"`
	Mute  *bool       `json:"mute,omitempty"`
	Deaf  *bool       `json:"deaf,omitempty"`
}

// MoveGuildMemberData is used to move a member
type MoveGuildMemberData struct {
	ChannelID *Snowflake `json:"channel_id"`
}

// UpdateSelfNick is used to update your own nick
type UpdateSelfNick struct {
	Nick *string `json:"nick"`
}
