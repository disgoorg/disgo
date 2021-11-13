package discord

type ThreadMember struct {
	ThreadID      Snowflake         `json:"id"`
	UserID        Snowflake         `json:"user_id"`
	JoinTimestamp Time              `json:"join_timestamp"`
	Flags         ThreadMemberFlags `json:"flags"`
}

type ThreadMemberFlags int
