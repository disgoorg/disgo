package discord

import "github.com/DisgoOrg/snowflake"

type ThreadMember struct {
	ThreadID      snowflake.Snowflake `json:"id"`
	UserID        snowflake.Snowflake `json:"user_id"`
	JoinTimestamp Time                `json:"join_timestamp"`
	Flags         ThreadMemberFlags   `json:"flags"`
}

type ThreadMemberFlags int
