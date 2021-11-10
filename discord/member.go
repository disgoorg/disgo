package discord

// Member is a discord GuildMember
type Member struct {
	GuildID       Snowflake    `json:"guild_id"`
	User          User         `json:"user"`
	Nick          *string      `json:"nick"`
	Avatar        *string      `json:"avatar"`
	RoleIDs       []Snowflake  `json:"roles,omitempty"`
	JoinedAt      Time         `json:"joined_at"`
	PremiumSince  *Time        `json:"premium_since,omitempty"`
	Deaf          bool         `json:"deaf,omitempty"`
	Mute          bool         `json:"mute,omitempty"`
	Pending       bool         `json:"pending"`
	TimedOutUntil Time         `json:"communication_disabled_until"`
	Permissions   *Permissions `json:"permissions"` // only sent from slash commands & should not be cached
}

// MemberAdd is used to add a member via the oauth2 access token to a guild
type MemberAdd struct {
	AccessToken string      `json:"access_token"`
	Nick        string      `json:"nick,omitempty"`
	Roles       []Snowflake `json:"roles,omitempty"`
	Mute        bool        `json:"mute,omitempty"`
	Deaf        bool        `json:"deaf,omitempty"`
}

// MemberUpdate is used to modify
type MemberUpdate struct {
	ChannelID *Snowflake  `json:"channel_id,omitempty"`
	Nick      *string     `json:"nick,omitempty"`
	Roles     []Snowflake `json:"roles,omitempty"`
	Mute      *bool       `json:"mute,omitempty"`
	Deaf      *bool       `json:"deaf,omitempty"`
}

// SelfNickUpdate is used to update your own nick
type SelfNickUpdate struct {
	Nick string `json:"nick"`
}
