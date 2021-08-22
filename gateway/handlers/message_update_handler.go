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
	return discord.Message{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *MessageUpdateHandler) HandleGatewayEvent(disgo core.Disgo, eventManager core.EventManager, sequenceNumber int, i interface{}) {
	message, ok := i.(discord.Message)
	if !ok {
		return
	}

	oldCoreMessage := disgo.Cache().MessageCache().GetCopy(message.ChannelID, message.ID)

	genericMessageEvent := &events.GenericMessageEvent{
		GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
		Message:      disgo.EntityBuilder().CreateMessage(message, core.CacheStrategyYes),
	}

	eventManager.Dispatch(&events.MessageUpdateEvent{
		GenericMessageEvent: genericMessageEvent,
		OldMessage:          oldCoreMessage,
	})

	if message.GuildID == nil {
		eventManager.Dispatch(&events.DMMessageUpdateEvent{
			GenericDMMessageEvent: &events.GenericDMMessageEvent{
				GenericMessageEvent: genericMessageEvent,
			},
			OldMessage: oldCoreMessage,
		})
	} else {
		eventManager.Dispatch(&events.GuildMessageUpdateEvent{
			GenericGuildMessageEvent: &events.GenericGuildMessageEvent{
				GenericMessageEvent: genericMessageEvent,
				GuildID:             *message.GuildID,
			},
			OldMessage: oldCoreMessage,
		})
	}
}
