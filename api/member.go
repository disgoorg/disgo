package api

// Member is a discord GuildMember
type Member struct {
	User
	Guild     Guild
	GuildID   Snowflake
	IsPending bool
}

// isOwner returns whether the member is the owner of the guild_events that it belongs to
func (m Member) isOwner() bool {
	return m.Guild.OwnerID == m.ID
}
