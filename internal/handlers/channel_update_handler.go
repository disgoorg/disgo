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
		newDMChannel := disgo.EntityBuilder().CreateDMChannel(channel, false)

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
		disgo.Cache().UncacheTextChannel(*channel.GuildID, channel.ID)
		textChannel := disgo.EntityBuilder().CreateTextChannel(channel, false)

		genericTextChannelEvent := events.GenericTextChannelEvent{
			GenericChannelEvent: genericChannelEvent,
		}
		eventManager.Dispatch(genericTextChannelEvent)

		eventManager.Dispatch(events.TextChannelUpdateEvent{
			GenericTextChannelEvent: genericTextChannelEvent,
			TextChannel:             textChannel,
		})

	case api.ChannelTypeStore:
		disgo.Cache().UncacheStoreChannel(*channel.GuildID, channel.ID)
		storeChannel := disgo.EntityBuilder().CreateStoreChannel(channel, false)

		genericStoreChannelEvent := events.GenericStoreChannelEvent{
			GenericChannelEvent: genericChannelEvent,
		}
		eventManager.Dispatch(genericStoreChannelEvent)

		eventManager.Dispatch(events.StoreChannelUpdateEvent{
			GenericStoreChannelEvent: genericStoreChannelEvent,
			StoreChannel:             storeChannel,
		})

	case api.ChannelTypeCategory:
		disgo.Cache().UncacheCategory(*channel.GuildID, channel.ID)
		category := disgo.EntityBuilder().CreateCategory(channel, false)

		genericCategoryEvent := events.GenericCategoryEvent{
			GenericChannelEvent: genericChannelEvent,
		}
		eventManager.Dispatch(genericCategoryEvent)

		eventManager.Dispatch(events.CategoryUpdateEvent{
			GenericCategoryEvent: genericCategoryEvent,
			Category:             category,
		})

	case api.ChannelTypeVoice:
		disgo.Cache().UncacheVoiceChannel(*channel.GuildID, channel.ID)
		voiceChannel := disgo.EntityBuilder().CreateVoiceChannel(channel, false)

		genericVoiceChannelEvent := events.GenericVoiceChannelEvent{
			GenericChannelEvent: genericChannelEvent,
		}
		eventManager.Dispatch(genericVoiceChannelEvent)

		eventManager.Dispatch(events.VoiceChannelUpdateEvent{
			GenericVoiceChannelEvent: genericVoiceChannelEvent,
			VoiceChannel:             voiceChannel,
		})

	default:
		log.Warnf("unknown channel type received: %d", channel.Type)
	}
}
