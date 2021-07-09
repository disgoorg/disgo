package events

import (
	"github.com/DisgoOrg/disgo/api"
)

// GenericTextChannelEvent is called upon receiving TextChannelCreateEvent, TextChannelUpdateEvent or TextChannelDeleteEvent
type GenericTextChannelEvent struct {
	*GenericGuildChannelEvent
	TextChannel api.TextChannel
}

// TextChannelCreateEvent indicates that a new api.TextChannel got created in a api.Guild
type TextChannelCreateEvent struct {
	*GenericTextChannelEvent
}

// TextChannelUpdateEvent indicates that a api.TextChannel got updated in a api.Guild
type TextChannelUpdateEvent struct {
	*GenericTextChannelEvent
	OldTextChannel api.TextChannel
}

// TextChannelDeleteEvent indicates that a api.TextChannel got deleted in a api.Guild
type TextChannelDeleteEvent struct {
	*GenericTextChannelEvent
}

// WebhooksUpdateEvent indicates that a api.Webhook updated in this api.TextChannel
type WebhooksUpdateEvent struct {
	*GenericTextChannelEvent
}
