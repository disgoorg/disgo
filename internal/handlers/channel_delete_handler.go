package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
	log "github.com/sirupsen/logrus"
)

// ChannelDeleteHandler handles api.GatewayEventChannelDelete
type ChannelDeleteHandler struct{}

// Event returns the raw gateway event Event
func (h ChannelDeleteHandler) Event() api.GatewayEventType {
	return api.GatewayEventChannelDelete
}

// New constructs a new payload receiver for the raw gateway event
func (h ChannelDeleteHandler) New() interface{} {
	return &api.Channel{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h ChannelDeleteHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
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
		disgo.Cache().UncacheDMChannel(channel.ID)
		dmChannel := disgo.EntityBuilder().CreateDMChannel(channel, api.CacheStrategyNo)

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
		disgo.Cache().UncacheTextChannel(*channel.GuildID, channel.ID)
		textChannel := disgo.EntityBuilder().CreateTextChannel(channel, api.CacheStrategyNo)

		genericTextChannelEvent := events.GenericTextChannelEvent{
			GenericChannelEvent: genericChannelEvent,
		}
		eventManager.Dispatch(genericTextChannelEvent)

		eventManager.Dispatch(events.TextChannelCreateEvent{
			GenericTextChannelEvent: genericTextChannelEvent,
			TextChannel:             textChannel,
		})

	case api.ChannelTypeStore:
		disgo.Cache().UncacheStoreChannel(*channel.GuildID, channel.ID)
		storeChannel := disgo.EntityBuilder().CreateStoreChannel(channel, api.CacheStrategyNo)

		genericStoreChannelEvent := events.GenericStoreChannelEvent{
			GenericChannelEvent: genericChannelEvent,
		}
		eventManager.Dispatch(genericStoreChannelEvent)

		eventManager.Dispatch(events.StoreChannelCreateEvent{
			GenericStoreChannelEvent: genericStoreChannelEvent,
			StoreChannel:             storeChannel,
		})

	case api.ChannelTypeCategory:
		disgo.Cache().UncacheCategory(*channel.GuildID, channel.ID)
		category := disgo.EntityBuilder().CreateCategory(channel, api.CacheStrategyNo)

		genericCategoryEvent := events.GenericCategoryEvent{
			GenericChannelEvent: genericChannelEvent,
		}
		eventManager.Dispatch(genericCategoryEvent)

		eventManager.Dispatch(events.CategoryCreateEvent{
			GenericCategoryEvent: genericCategoryEvent,
			Category:             category,
		})

	case api.ChannelTypeVoice:
		disgo.Cache().UncacheVoiceChannel(*channel.GuildID, channel.ID)
		voiceChannel := disgo.EntityBuilder().CreateVoiceChannel(channel, api.CacheStrategyNo)

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
