package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
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

		genericDMChannelEvent := events.GenericDMChannelEvent{
			GenericChannelEvent: genericChannelEvent,
			DMChannel:           disgo.EntityBuilder().CreateDMChannel(channel, api.CacheStrategyNo),
		}
		eventManager.Dispatch(genericDMChannelEvent)

		eventManager.Dispatch(events.DMChannelCreateEvent{
			GenericDMChannelEvent: genericDMChannelEvent,
		})

	case api.ChannelTypeGroupDM:
		disgo.Logger().Warnf("ChannelTypeGroupDM received what the hell discord")

	case api.ChannelTypeText, api.ChannelTypeNews:
		disgo.Cache().UncacheTextChannel(*channel.GuildID, channel.ID)

		genericTextChannelEvent := events.GenericTextChannelEvent{
			GenericChannelEvent: genericChannelEvent,
			TextChannel:         disgo.EntityBuilder().CreateTextChannel(channel, api.CacheStrategyNo),
		}
		eventManager.Dispatch(genericTextChannelEvent)

		eventManager.Dispatch(events.TextChannelCreateEvent{
			GenericTextChannelEvent: genericTextChannelEvent,
		})

	case api.ChannelTypeStore:
		disgo.Cache().UncacheStoreChannel(*channel.GuildID, channel.ID)

		genericStoreChannelEvent := events.GenericStoreChannelEvent{
			GenericChannelEvent: genericChannelEvent,
			StoreChannel:        disgo.EntityBuilder().CreateStoreChannel(channel, api.CacheStrategyNo),
		}
		eventManager.Dispatch(genericStoreChannelEvent)

		eventManager.Dispatch(events.StoreChannelCreateEvent{
			GenericStoreChannelEvent: genericStoreChannelEvent,
		})

	case api.ChannelTypeCategory:
		disgo.Cache().UncacheCategory(*channel.GuildID, channel.ID)

		genericCategoryEvent := events.GenericCategoryEvent{
			GenericChannelEvent: genericChannelEvent,
			Category:            disgo.EntityBuilder().CreateCategory(channel, api.CacheStrategyNo),
		}
		eventManager.Dispatch(genericCategoryEvent)

		eventManager.Dispatch(events.CategoryCreateEvent{
			GenericCategoryEvent: genericCategoryEvent,
		})

	case api.ChannelTypeVoice:
		disgo.Cache().UncacheVoiceChannel(*channel.GuildID, channel.ID)

		genericVoiceChannelEvent := events.GenericVoiceChannelEvent{
			GenericChannelEvent: genericChannelEvent,
			VoiceChannel:        disgo.EntityBuilder().CreateVoiceChannel(channel, api.CacheStrategyNo),
		}
		eventManager.Dispatch(genericVoiceChannelEvent)

		eventManager.Dispatch(events.VoiceChannelCreateEvent{
			GenericVoiceChannelEvent: genericVoiceChannelEvent,
		})

	default:
		disgo.Logger().Warnf("unknown channel type received: %d", channel.Type)
	}
}
