package handlers

import (
	"github.com/DiscoOrg/disgo/api"
	"github.com/DiscoOrg/disgo/api/events"
)

type MessageCreateHandler struct{}

func (h MessageCreateHandler) New() interface{} {
	return &api.Message{}
}

func (h MessageCreateHandler) Handle(disgo api.Disgo, eventManager api.EventManager, i interface{}) {
	message, ok := i.(*api.Message)
	if !ok {
		return
	}
	if message.GuildID == nil {
		// dm channel
	} else{
		// text channel
		message.Disgo = disgo
		message.Author.Disgo = disgo
		eventManager.Dispatch(events.GuildMessageReceivedEvent{
			Message: *message,
			GenericGuildMessageEvent: events.GenericGuildMessageEvent{
				GenericGuildEvent:   events.GenericGuildEvent{
					Event:   api.Event{
						Disgo: disgo,
					},
					GuildID: *message.GuildID,
				},
				GenericMessageEvent: events.GenericMessageEvent{
					Event:   api.Event{
						Disgo: disgo,
					},
					MessageChannelID: message.ChannelID,
					MessageID: message.ID,
				},
			},
		})
	}
	eventManager.Dispatch(events.MessageReceivedEvent{
		GenericMessageEvent: events.GenericMessageEvent{
			Event:   api.Event{
				Disgo: disgo,
			},
			MessageChannelID: message.ChannelID,
			MessageID: message.ID,
		},
		Message: *message,
	})
}
