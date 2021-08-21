package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
)

// MessageUpdateHandler handles api.GatewayEventMessageUpdate
type MessageUpdateHandler struct{}

// EventType returns the api.GatewayEventType
func (h *MessageUpdateHandler) EventType() gateway.EventType {
	return gateway.EventTypeMessageUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *MessageUpdateHandler) New() interface{} {
	return &discord.Message{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *MessageUpdateHandler) HandleGatewayEvent(disgo core.Disgo, eventManager core.EventManager, sequenceNumber int, i interface{}) {
	message, ok := i.(*discord.Message)
	if !ok {
		return
	}

	oldMessage := disgo.Cache().Message(message.ChannelID, message.ID)
	if oldMessage != nil {
		oldMessage = &*oldMessage
	}

	message = disgo.EntityBuilder().CreateMessage(message, core.CacheStrategyYes)

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
