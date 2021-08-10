package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
)

// MessageCreateHandler handles api.GatewayEventMessageCreate
type MessageCreateHandler struct{}

// Event returns the api.GatewayEventType
func (h *MessageCreateHandler) EventType() gateway.EventType {
	return gateway.EventTypeMessageCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h *MessageCreateHandler) New() interface{} {
	return &discord.Message{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *MessageCreateHandler) HandleGatewayEvent(disgo core.Disgo, eventManager core.EventManager, sequenceNumber int, i interface{}) {
	message, ok := i.(*discord.Message)
	if !ok {
		return
	}

	message = disgo.EntityBuilder().CreateMessage(message, core.CacheStrategyYes)

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
