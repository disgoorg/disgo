package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
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

		genericDMChannelEvent := events.GenericDMChannelEvent{
			GenericChannelEvent: genericChannelEvent,
			DMChannel:           disgo.EntityBuilder().CreateDMChannel(channel, api.CacheStrategyYes),
		}
		eventManager.Dispatch(genericDMChannelEvent)

		eventManager.Dispatch(events.DMChannelUpdateEvent{
			GenericDMChannelEvent: genericDMChannelEvent,
			OldDMChannel:          oldDMChannel,
		})

	case api.ChannelTypeGroupDM:
		disgo.Logger().Warnf("ChannelTypeGroupDM received what the hell discord")

	case api.ChannelTypeText, api.ChannelTypeNews:
		oldTextChannel := disgo.Cache().TextChannel(channel.ID)
		if oldTextChannel != nil {
			oldTextChannel = &*oldTextChannel
		}

		genericTextChannelEvent := events.GenericTextChannelEvent{
			GenericChannelEvent: genericChannelEvent,
			TextChannel:         disgo.EntityBuilder().CreateTextChannel(channel, api.CacheStrategyYes),
		}
		eventManager.Dispatch(genericTextChannelEvent)

		eventManager.Dispatch(events.TextChannelUpdateEvent{
			GenericTextChannelEvent: genericTextChannelEvent,
			OldTextChannel:          oldTextChannel,
		})

	case api.ChannelTypeStore:
		oldStoreChannel := disgo.Cache().StoreChannel(channel.ID)
		if oldStoreChannel != nil {
			oldStoreChannel = &*oldStoreChannel
		}

		genericStoreChannelEvent := events.GenericStoreChannelEvent{
			GenericChannelEvent: genericChannelEvent,
			StoreChannel:        disgo.EntityBuilder().CreateStoreChannel(channel, api.CacheStrategyYes),
		}
		eventManager.Dispatch(genericStoreChannelEvent)

		eventManager.Dispatch(events.StoreChannelUpdateEvent{
			GenericStoreChannelEvent: genericStoreChannelEvent,
			OldStoreChannel:          oldStoreChannel,
		})

	case api.ChannelTypeCategory:
		oldCategory := disgo.Cache().Category(channel.ID)
		if oldCategory != nil {
			oldCategory = &*oldCategory
		}

		genericCategoryEvent := events.GenericCategoryEvent{
			GenericChannelEvent: genericChannelEvent,
			Category:            disgo.EntityBuilder().CreateCategory(channel, api.CacheStrategyYes),
		}
		eventManager.Dispatch(genericCategoryEvent)

		eventManager.Dispatch(events.CategoryUpdateEvent{
			GenericCategoryEvent: genericCategoryEvent,
			OldCategory:          oldCategory,
		})

	case api.ChannelTypeVoice:
		oldVoiceChannel := disgo.Cache().VoiceChannel(channel.ID)
		if oldVoiceChannel != nil {
			oldVoiceChannel = &*oldVoiceChannel
		}

		genericVoiceChannelEvent := events.GenericVoiceChannelEvent{
			GenericChannelEvent: genericChannelEvent,
			VoiceChannel:        disgo.EntityBuilder().CreateVoiceChannel(channel, api.CacheStrategyYes),
		}
		eventManager.Dispatch(genericVoiceChannelEvent)

		eventManager.Dispatch(events.VoiceChannelUpdateEvent{
			GenericVoiceChannelEvent: genericVoiceChannelEvent,
			OldVoiceChannel:          oldVoiceChannel,
		})

	default:
		disgo.Logger().Warnf("unknown channel type received: %d", channel.Type)
	}
}
