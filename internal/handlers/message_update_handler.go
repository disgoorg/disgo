package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

// MessageUpdateHandler handles api.GatewayEventMessageUpdate
type MessageUpdateHandler struct{}

// Event returns the raw gateway event Event
func (h *MessageUpdateHandler) Event() api.GatewayEventType {
	return api.GatewayEventMessageUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *MessageUpdateHandler) New() interface{} {
	return &api.Message{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *MessageUpdateHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
	message, ok := i.(*api.Message)
	if !ok {
		return
	}

	oldMessage := disgo.Cache().Message(message.ChannelID, message.ID)
	if oldMessage != nil {
		oldMessage = &*oldMessage
	}

	message = disgo.EntityBuilder().CreateMessage(message, api.CacheStrategyYes)

	genericMessageEvent := &events.GenericMessageEvent{
		GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
		Message:      message,
	}

	eventManager.Dispatch(&events.MessageUpdateEvent{
		GenericMessageEvent: genericMessageEvent,
		OldMessage:          oldMessage,
	})

	if message.GuildID == nil {
		eventManager.Dispatch(&events.DMMessageUpdateEvent{
			GenericDMMessageEvent: &events.GenericDMMessageEvent{
				GenericMessageEvent: genericMessageEvent,
			},
			OldMessage: oldMessage,
		})
	} else {

		eventManager.Dispatch(&events.GuildMessageUpdateEvent{
			GenericGuildMessageEvent: &events.GenericGuildMessageEvent{
				GenericMessageEvent: genericMessageEvent,
				GuildID:             *message.GuildID,
			},
			OldMessage: oldMessage,
		})
	}
}
