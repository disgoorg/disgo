package api

type Member struct {
	User
	GuildID   Snowflake
	Guild     Guild
	IsPending bool
}

func (m Member) isOwner() bool {
	return m.Guild.OwnerID == m.ID
}
