package handlers

import (
	log "github.com/sirupsen/logrus"

	"github.com/DiscoOrg/disgo/api"
	"github.com/DiscoOrg/disgo/internal/events"
)

// GuildCreatePayload payload from GUILD_CREATE gateways event sent by discord
type MessageCreateEvent struct {
	api.Message
}

type MessageCreateHandler struct{}

func (h MessageCreateHandler) New() interface{} {
	return &MessageCreateEvent{}
}

func (h MessageCreateHandler) Handle(disgo api.Disgo, eventManager api.EventManager, i interface{}) {
	message, ok := i.(*MessageCreateEvent)
	if !ok {
		return
	}
	switch message.ChannelType {
	case api.GuildTextChannel:
		message := message.Message
		message.Disgo = disgo
		message.Author.Disgo = disgo
		eventManager.Dispatch(events.GuildMessageReceivedEvent{
			Message: message,
			GenericGuildMessageEvent: events.GenericGuildMessageEvent{
				TextChannel: api.TextChannel{
					GuildChannel:   api.GuildChannel{
						Channel: api.Channel{
							Disgo: disgo,
							ID:    message.ChannelID,
						},
						GuildID: message.GuildId,
						Guild:   api.Guild{
							ID: message.GuildId,
						},
					},
					MessageChannel: api.MessageChannel{
						Channel: api.Channel{
							Disgo: disgo,
							ID:    message.ChannelID,
						},
					},
				},
			},
		})
	case api.DMTextChannel:
	case api.GroupDMChannel:
	default:
		log.Errorf("unknown channel type received: %d", message.ChannelType)
	}
	eventManager.Dispatch(events.MessageReceivedEvent{
		Message: message.Message,
	})
}
