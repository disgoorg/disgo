package api

import "time"

type ThreadMember struct {
	Disgo         Disgo
	GuildID       Snowflake
	ThreadID      Snowflake `json:"thread_id"`
	UserID        Snowflake `json:"user_id"`
	JoinTimestamp time.Time `json:"join_timestamp"`
	Flags         int       `json:"flags"`
}

func (m *ThreadMember) Thread() Thread {
	return m.Disgo.Cache().Thread(m.ThreadID)
}

func (m *ThreadMember) User() *User {
	return m.Disgo.Cache().User(m.UserID)
}

func (m *ThreadMember) Guild() *Guild {
	return m.Disgo.Cache().Guild(m.UserID)
}
