package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
	log "github.com/sirupsen/logrus"
)

// ChannelUpdateHandler handles api.GatewayEventChannelUpdate
type ChannelUpdateHandler struct{}

// Event returns the raw gateway event Event
func (h ChannelUpdateHandler) Event() api.GatewayEventType {
	return api.GatewayEventChannelUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h ChannelUpdateHandler) New() interface{} {
	return &api.Channel{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h ChannelUpdateHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
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
		oldDMChannel := disgo.Cache().DMChannel(channel.ID)
		if oldDMChannel != nil {
			oldDMChannel = &*oldDMChannel
		}
		newDMChannel := disgo.EntityBuilder().CreateDMChannel(channel, api.CacheStrategyYes)

		genericDMChannelEvent := events.GenericDMChannelEvent{
			GenericChannelEvent: genericChannelEvent,
		}
		eventManager.Dispatch(genericDMChannelEvent)

		eventManager.Dispatch(events.DMChannelUpdateEvent{
			GenericDMChannelEvent: genericDMChannelEvent,
			NewDMChannel:             newDMChannel,
			OldDMChannel:             oldDMChannel,
		})

	case api.ChannelTypeGroupDM:
		log.Warnf("ChannelTypeGroupDM received what the hell discord")

	case api.ChannelTypeText, api.ChannelTypeNews:
		oldTextChannel := disgo.Cache().TextChannel(channel.ID)
		if oldTextChannel != nil {
			oldTextChannel = &*oldTextChannel
		}
		newTextChannel := disgo.EntityBuilder().CreateTextChannel(channel, api.CacheStrategyYes)

		genericTextChannelEvent := events.GenericTextChannelEvent{
			GenericChannelEvent: genericChannelEvent,
		}
		eventManager.Dispatch(genericTextChannelEvent)

		eventManager.Dispatch(events.TextChannelUpdateEvent{
			GenericTextChannelEvent: genericTextChannelEvent,
			NewTextChannel:             newTextChannel,
			OldTextChannel:             oldTextChannel,
		})

	case api.ChannelTypeStore:
		oldStoreChannel := disgo.Cache().StoreChannel(channel.ID)
		if oldStoreChannel != nil {
			oldStoreChannel = &*oldStoreChannel
		}
		newStoreChannel := disgo.EntityBuilder().CreateStoreChannel(channel, api.CacheStrategyYes)

		genericStoreChannelEvent := events.GenericStoreChannelEvent{
			GenericChannelEvent: genericChannelEvent,
		}
		eventManager.Dispatch(genericStoreChannelEvent)

		eventManager.Dispatch(events.StoreChannelUpdateEvent{
			GenericStoreChannelEvent: genericStoreChannelEvent,
			NewStoreChannel:             newStoreChannel,
			OldStoreChannel:             oldStoreChannel,
		})

	case api.ChannelTypeCategory:
		oldCategory := disgo.Cache().Category(channel.ID)
		if oldCategory != nil {
			oldCategory = &*oldCategory
		}
		newCategory := disgo.EntityBuilder().CreateCategory(channel, api.CacheStrategyYes)

		genericCategoryEvent := events.GenericCategoryEvent{
			GenericChannelEvent: genericChannelEvent,
		}
		eventManager.Dispatch(genericCategoryEvent)

		eventManager.Dispatch(events.CategoryUpdateEvent{
			GenericCategoryEvent: genericCategoryEvent,
			NewCategory:             newCategory,
			OldCategory:             oldCategory,
		})

	case api.ChannelTypeVoice:
		oldVoiceChannel := disgo.Cache().VoiceChannel(channel.ID)
		if oldVoiceChannel != nil {
			oldVoiceChannel = &*oldVoiceChannel
		}
		newVoiceChannel := disgo.EntityBuilder().CreateVoiceChannel(channel, api.CacheStrategyYes)

		genericVoiceChannelEvent := events.GenericVoiceChannelEvent{
			GenericChannelEvent: genericChannelEvent,
		}
		eventManager.Dispatch(genericVoiceChannelEvent)

		eventManager.Dispatch(events.VoiceChannelUpdateEvent{
			GenericVoiceChannelEvent: genericVoiceChannelEvent,
			NewVoiceChannel:             newVoiceChannel,
			OldVoiceChannel:             oldVoiceChannel,
		})

	default:
		log.Warnf("unknown channel type received: %d", channel.Type)
	}
}
