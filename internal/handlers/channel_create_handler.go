package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
	log "github.com/sirupsen/logrus"
)

// ChannelCreateHandler handles api.GatewayEventChannelCreate
type ChannelCreateHandler struct{}

// Event returns the raw gateway event Event
func (h ChannelCreateHandler) Event() api.GatewayEventType {
	return api.GatewayEventChannelCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h ChannelCreateHandler) New() interface{} {
	return &api.Channel{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h ChannelCreateHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
	channel, ok := i.(*api.Channel)
	if !ok {
		return
	}

	genericChannelEvent := events.GenericChannelEvent{
		GenericEvent: events.NewEvent(disgo, sequenceNumber),
		ChannelID:    channel.ID,
	}
	eventManager.Dispatch(genericChannelEvent)

	switch channel.Type {
	case api.ChannelTypeDM:
		dmChannel := disgo.EntityBuilder().CreateDMChannel(channel, api.CacheStrategyYes)

		genericDMChannelEvent := events.GenericDMChannelEvent{
			GenericChannelEvent: genericChannelEvent,
		}
		eventManager.Dispatch(genericDMChannelEvent)

		eventManager.Dispatch(events.DMChannelCreateEvent{
			GenericDMChannelEvent: genericDMChannelEvent,
			DMChannel:             dmChannel,
		})

	case api.ChannelTypeGroupDM:
		log.Warnf("ChannelTypeGroupDM received what the hell discord")

	case api.ChannelTypeText, api.ChannelTypeNews:
		textChannel := disgo.EntityBuilder().CreateTextChannel(channel, api.CacheStrategyYes)

		genericTextChannelEvent := events.GenericTextChannelEvent{
			GenericChannelEvent: genericChannelEvent,
		}
		eventManager.Dispatch(genericTextChannelEvent)

		eventManager.Dispatch(events.TextChannelCreateEvent{
			GenericTextChannelEvent: genericTextChannelEvent,
			TextChannel:             textChannel,
		})

	case api.ChannelTypeStore:
		storeChannel := disgo.EntityBuilder().CreateStoreChannel(channel, api.CacheStrategyYes)

		genericStoreChannelEvent := events.GenericStoreChannelEvent{
			GenericChannelEvent: genericChannelEvent,
		}
		eventManager.Dispatch(genericStoreChannelEvent)

		eventManager.Dispatch(events.StoreChannelCreateEvent{
			GenericStoreChannelEvent: genericStoreChannelEvent,
			StoreChannel:             storeChannel,
		})

	case api.ChannelTypeCategory:
		category := disgo.EntityBuilder().CreateCategory(channel, api.CacheStrategyYes)

		genericCategoryEvent := events.GenericCategoryEvent{
			GenericChannelEvent: genericChannelEvent,
		}
		eventManager.Dispatch(genericCategoryEvent)

		eventManager.Dispatch(events.CategoryCreateEvent{
			GenericCategoryEvent: genericCategoryEvent,
			Category:             category,
		})

	case api.ChannelTypeVoice:
		voiceChannel := disgo.EntityBuilder().CreateVoiceChannel(channel, api.CacheStrategyYes)

		genericVoiceChannelEvent := events.GenericVoiceChannelEvent{
			GenericChannelEvent: genericChannelEvent,
		}
		eventManager.Dispatch(genericVoiceChannelEvent)

		eventManager.Dispatch(events.VoiceChannelCreateEvent{
			GenericVoiceChannelEvent: genericVoiceChannelEvent,
			VoiceChannel:             voiceChannel,
		})

	default:
		log.Warnf("unknown channel type received: %d", channel.Type)
	}

}
