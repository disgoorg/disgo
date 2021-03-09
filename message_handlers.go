package disgo

import (
	log "github.com/sirupsen/logrus"

	"github.com/DiscoOrg/disgo/constants"
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

func (h MessageCreateHandler) Handle(eventManager EventManager, i interface{}) {
	message, ok := i.(*MessageCreateEvent)
	if !ok {
		return
	}
	switch message.ChannelType {
	case constants.GuildTextChannel:
		eventManager.Dispatch(GuildMessageReceivedEvent{
			Message: message.Message,
			GenericGuildMessageEvent: GenericGuildMessageEvent{
				TextChannel: TextChannel{
					GuildChannel:   GuildChannel{
						Channel{
							Disgo: eventManager.Disgo(),
							ID:    message.ChannelID,
						},
					},
					MessageChannel: MessageChannel{
						Channel{
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
	eventManager.Dispatch(MessageReceivedEvent{
		Message: message.Message,
	})
}
