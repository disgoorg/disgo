package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

// MessageCreateHandler handles api.GatewayEventMessageCreate
type MessageCreateHandler struct{}

// Event returns the raw gateway event Event
func (h *MessageCreateHandler) Event() api.GatewayEventType {
	return api.GatewayEventMessageCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h *MessageCreateHandler) New() interface{} {
	return &api.Message{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *MessageCreateHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
	message, ok := i.(*api.Message)
	if !ok {
		return
	}

	message = disgo.EntityBuilder().CreateMessage(message, api.CacheStrategyYes)

	genericMessageEvent := &events.GenericMessageEvent{
		GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
		MessageID:    message.ID,
		Message:      message,
		ChannelID:    message.ChannelID,
	}

	eventManager.Dispatch(&events.MessageCreateEvent{
		GenericMessageEvent: genericMessageEvent,
	})

	if message.GuildID == nil {
		eventManager.Dispatch(&events.DMMessageCreateEvent{
			GenericDMMessageEvent: &events.GenericDMMessageEvent{
				GenericMessageEvent: genericMessageEvent,
			},
		})
	} else {
		eventManager.Dispatch(&events.GuildMessageCreateEvent{
			GenericGuildMessageEvent: &events.GenericGuildMessageEvent{
				GenericMessageEvent: genericMessageEvent,
				GuildID:             *message.GuildID,
			},
		})
	}

}
