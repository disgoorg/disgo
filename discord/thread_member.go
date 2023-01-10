package discord

import (
	"time"

	"github.com/disgoorg/snowflake/v2"
)

type ThreadMember struct {
	ThreadID      snowflake.ID      `json:"id"`
	UserID        snowflake.ID      `json:"user_id"`
	JoinTimestamp time.Time         `json:"join_timestamp"`
	Flags         ThreadMemberFlags `json:"flags"`
	Member        *Member           `json:"member,omitempty"`
}

type ThreadMemberFlags int
