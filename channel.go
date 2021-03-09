package disgo

import (
	"github.com/chebyrash/promise"

	"github.com/DiscoOrg/disgo/constants"
	"github.com/DiscoOrg/disgo/models"
)

// Channel is a generic discord channel object
type Channel struct {
	Disgo Disgo
	ID    models.Snowflake      `json:"id"`
	Type  constants.ChannelType `json:"type"`
}

//  MessageChannel is used for sending messages to user
type MessageChannel struct {
	Channel
}

func (c MessageChannel) SendMessage(content string) *promise.Promise {
	return c.Disgo.RestClient().SendMessage(c.ID, models.Message{Content: content})
}


// DMChannel is used for interacting in private messages with users
type DMChannel struct {
	MessageChannel
}

// GuildChannel is a generic type for all server channels
type GuildChannel struct {
	Channel
}

// CategoryChannel groups text & voice channels in servers together
type CategoryChannel struct {
	GuildChannel
}

//  VoiceChannel adds methods specifically for interacting with discord's voice
type VoiceChannel struct {
	GuildChannel
}

// TextChannel allows you to interact with discord's text channels
type TextChannel struct {
	GuildChannel
	MessageChannel
}


// StoreChannel allows you to interact with discord's store channels
type StoreChannel struct {
	GuildChannel
}

// NewsChannel allows you to interact with discord's news channels
type NewsChannel struct {
	TextChannel
}
