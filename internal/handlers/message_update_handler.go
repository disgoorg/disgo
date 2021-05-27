package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

// MessageUpdateHandler handles api.GatewayEventMessageUpdate
type MessageUpdateHandler struct{}

// Event returns the raw gateway event Event
func (h MessageUpdateHandler) Event() api.GatewayEventType {
	return api.GatewayEventMessageUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h MessageUpdateHandler) New() interface{} {
	return &api.FullMessage{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h MessageUpdateHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
	fullMessage, ok := i.(*api.FullMessage)
	if !ok {
		return
	}

	oldMessage := disgo.Cache().Message(fullMessage.ChannelID, fullMessage.ID)
	if oldMessage != nil {
		oldMessage = &*oldMessage
	}

	message := disgo.EntityBuilder().CreateMessage(fullMessage, api.CacheStrategyYes)

	genericMessageEvent := events.GenericMessageEvent{
		GenericEvent: events.NewEvent(disgo, sequenceNumber),
		Message:      message,
	}
	eventManager.Dispatch(genericMessageEvent)

	eventManager.Dispatch(events.MessageUpdateEvent{
		GenericMessageEvent: genericMessageEvent,
		OldMessage:          oldMessage,
	})

	if message.GuildID == nil {
		genericDMMessageEvent := events.GenericDMMessageEvent{
			GenericMessageEvent: genericMessageEvent,
		}
		eventManager.Dispatch(genericDMMessageEvent)

		eventManager.Dispatch(events.DMMessageUpdateEvent{
			GenericDMMessageEvent: genericDMMessageEvent,
			OldMessage:            oldMessage,
		})
	} else {
		genericGuildMessageEvent := events.GenericGuildMessageEvent{
			GenericMessageEvent: genericMessageEvent,
			GuildID:             *message.GuildID,
		}
		eventManager.Dispatch(genericGuildMessageEvent)

		eventManager.Dispatch(events.GuildMessageUpdateEvent{
			GenericGuildMessageEvent: genericGuildMessageEvent,
			OldMessage:               oldMessage,
		})
	}
}
