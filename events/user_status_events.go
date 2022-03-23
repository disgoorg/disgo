package events

import (
	"github.com/DisgoOrg/snowflake"
	"github.com/disgoorg/disgo/discord"
)

// UserStatusUpdateEvent generic Status event
type UserStatusUpdateEvent struct {
	*GenericEvent
	UserID    snowflake.Snowflake
	OldStatus discord.OnlineStatus
	Status    discord.OnlineStatus
}

// UserClientStatusUpdateEvent generic client-specific Status event
type UserClientStatusUpdateEvent struct {
	*GenericEvent
	UserID          snowflake.Snowflake
	OldClientStatus *discord.ClientStatus
	ClientStatus    discord.ClientStatus
}
