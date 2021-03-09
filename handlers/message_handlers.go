package handlers

import (
	log "github.com/sirupsen/logrus"

	"github.com/DiscoOrg/disgo"
	"github.com/DiscoOrg/disgo/constants"
	"github.com/DiscoOrg/disgo/events"
	"github.com/DiscoOrg/disgo/models"
)

// GuildCreatePayload payload from GUILD_CREATE gateways event sent by discord
type MessageCreateEvent struct {
	models.Message
}

type MessageCreateHandler struct{}

func (h MessageCreateHandler) New() interface{} {
	return &MessageCreateEvent{}
}

func (h MessageCreateHandler) Handle(eventManager disgo.EventManager, i interface{}) {
	message, ok := i.(*MessageCreateEvent)
	if !ok {
		return
	}
	switch message.ChannelType {
	case constants.GuildTextChannel:
		message := message.Message
		message.Disgo = eventManager.Disgo()
		message.User.Disgo = eventManager.Disgo()
		eventManager.Dispatch(events.GuildMessageReceivedEvent{
			Message: message,
			GenericGuildMessageEvent: events.GenericGuildMessageEvent{
				TextChannel: models.TextChannel{
					GuildChannel:   models.GuildChannel{
						Channel: models.Channel{
							Disgo: eventManager.Disgo(),
							ID:    message.ChannelID,
						},
						GuildID: message.GuildId,
						Guild:   models.Guild{
							ID: message.GuildId,
						},
					},
					MessageChannel: models.MessageChannel{
						Channel: models.Channel{
							Disgo: eventManager.Disgo(),
							ID:    message.ChannelID,
						},
					},
				},
			},
		})
	case constants.DMChannel:
	case constants.GroupDMChannel:
	default:
		log.Errorf("unknown channel type received: %d", message.ChannelType)
	}
	eventManager.Dispatch(events.MessageReceivedEvent{
		Message: message.Message,
	})
}
