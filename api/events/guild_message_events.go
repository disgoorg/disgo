package events

import (
	"github.com/DisgoOrg/disgo/api"
)

// GenericDMMessageEvent generic api.DMChannel api.Message api.GenericEvent
type GenericGuildMessageEvent struct {
	GenericMessageEvent
	GuildID api.Snowflake
}

func (e GenericGuildMessageEvent) Guild() *api.Guild {
	return e.Disgo().Cache().Guild(e.GuildID)
}

func (e GenericGuildMessageEvent) TextChannel() *api.TextChannel {
	return e.Disgo().Cache().TextChannel(e.ChannelID)
}

// DMMessageReceivedEvent called upon receiving a api.Message in a api.DMChannel
type GuildMessageReceivedEvent struct {
	GenericGuildMessageEvent
	Message *api.Message
}

// DMMessageUpdateEvent called upon editing a api.Message in a api.DMChannel
type GuildMessageUpdateEvent struct {
	GenericGuildMessageEvent
	NewMessage *api.Message
	OldMessage *api.Message
}

// DMMessageDeleteEvent called upon deleting a api.Message in a api.DMChannel
type GuildMessageDeleteEvent struct {
	GenericGuildMessageEvent
	Message *api.Message
}
