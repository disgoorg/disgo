package discord

import (
	"time"

	"github.com/disgoorg/snowflake"
)

type ThreadMember struct {
	ThreadID      snowflake.Snowflake `json:"id"`
	UserID        snowflake.Snowflake `json:"user_id"`
	JoinTimestamp time.Time           `json:"join_timestamp"`
	Flags         ThreadMemberFlags   `json:"flags"`
}

type ThreadMemberFlags int
