package events

import (
	"github.com/snekROmonoro/snowflake"

	"github.com/snekROmonoro/disgo/discord"
)

// UserStatusUpdate generic Status event
type UserStatusUpdate struct {
	*GenericEvent
	UserID    snowflake.ID
	OldStatus discord.OnlineStatus
	Status    discord.OnlineStatus
}

// UserClientStatusUpdate generic client-specific Status event
type UserClientStatusUpdate struct {
	*GenericEvent
	UserID          snowflake.ID
	OldClientStatus discord.ClientStatus
	ClientStatus    discord.ClientStatus
}
