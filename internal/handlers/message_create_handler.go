package handlers

import (
	"github.com/DiscoOrg/disgo/api"
	"github.com/DiscoOrg/disgo/api/events"
)

// MessageCreateHandler handles api.MessageCreateGatewayEvent
type MessageCreateHandler struct{}

// Name returns the raw gateway event name
func (h MessageCreateHandler) Name() string {
	return api.MessageCreateGatewayEvent
}

// New constructs a new payload receiver for the raw gateway event
func (h MessageCreateHandler) New() interface{} {
	return &api.Message{}
}

// Handle handles the specific raw gateway event
func (h MessageCreateHandler) Handle(disgo api.Disgo, eventManager api.EventManager, i interface{}) {
	message, ok := i.(*api.Message)
	if !ok {
		return
	}

	genericMessageEvent := events.GenericMessageEvent{
		Event: api.Event{
			Disgo: disgo,
		},
		MessageChannelID: message.ChannelID,
		MessageID:        message.ID,
	}
	eventManager.Dispatch(genericMessageEvent)

	genericGuildEvent := events.GenericGuildEvent{
		Event: api.Event{
			Disgo: disgo,
		},
		GuildID: *message.GuildID,
	}
	eventManager.Dispatch(genericGuildEvent)

	eventManager.Dispatch(events.MessageReceivedEvent{
		GenericMessageEvent: genericMessageEvent,
		Message:             *message,
	})

	if message.GuildID == nil {
		// dm channel
	} else {
		// text channel
		message.Disgo = disgo
		message.Author.Disgo = disgo
		eventManager.Dispatch(events.GuildMessageReceivedEvent{
			Message: *message,
			GenericGuildMessageEvent: events.GenericGuildMessageEvent{
				GenericGuildEvent:   genericGuildEvent,
				GenericMessageEvent: genericMessageEvent,
			},
		})
	}

}
