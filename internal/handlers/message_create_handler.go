package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

// MessageCreateHandler handles api.GatewayEventMessageCreate
type MessageCreateHandler struct{}

// Event returns the raw gateway event Event
func (h MessageCreateHandler) Event() api.GatewayEventType {
	return api.GatewayEventMessageCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h MessageCreateHandler) New() interface{} {
	return &api.Message{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h MessageCreateHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
	message, ok := i.(*api.Message)
	if !ok {
		return
	}

	message = disgo.EntityBuilder().CreateMessage(message, api.CacheStrategyYes)

	genericMessageEvent := events.GenericMessageEvent{
		GenericEvent: events.NewEvent(disgo, sequenceNumber),
		Message:      message,
	}
	eventManager.Dispatch(genericMessageEvent)

	eventManager.Dispatch(events.MessageCreateEvent{
		GenericMessageEvent: genericMessageEvent,
	})

	if message.GuildID == nil {
		genericDMMessageEvent := events.GenericDMMessageEvent{
			GenericMessageEvent: genericMessageEvent,
		}
		eventManager.Dispatch(genericDMMessageEvent)

		eventManager.Dispatch(events.DMMessageCreateEvent{
			GenericDMMessageEvent: genericDMMessageEvent,
		})
	} else {
		genericGuildMessageEvent := events.GenericGuildMessageEvent{
			GenericMessageEvent: genericMessageEvent,
			GuildID:             *message.GuildID,
		}
		eventManager.Dispatch(genericGuildMessageEvent)

		eventManager.Dispatch(events.GuildMessageCreateEvent{
			GenericGuildMessageEvent: genericGuildMessageEvent,
		})
	}

}
