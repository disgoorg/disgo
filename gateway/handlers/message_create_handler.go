package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

// MessageCreateHandler handles api.GatewayEventMessageCreate
type MessageCreateHandler struct{}

// EventType returns the api.GatewayGatewayEventType
func (h *MessageCreateHandler) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeMessageCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h *MessageCreateHandler) New() interface{} {
	return discord.Message{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *MessageCreateHandler) HandleGatewayEvent(disgo core.Disgo, eventManager core.EventManager, sequenceNumber int, i interface{}) {
	message, ok := i.(discord.Message)
	if !ok {
		return
	}

	genericMessageEvent := &events.GenericMessageEvent{
		GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
		MessageID:    message.ID,
		Message:      disgo.EntityBuilder().CreateMessage(message, core.CacheStrategyYes),
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
