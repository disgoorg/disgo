package discord

import "github.com/disgoorg/snowflake/v2"

type ThreadMember struct {
	ThreadID      snowflake.ID      `json:"id"`
	UserID        snowflake.ID      `json:"user_id"`
	JoinTimestamp Time              `json:"join_timestamp"`
	Flags         ThreadMemberFlags `json:"flags"`
}

type ThreadMemberFlags int
