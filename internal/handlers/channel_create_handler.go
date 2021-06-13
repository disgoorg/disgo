package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
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

	channel.Disgo = disgo

	genericChannelEvent := events.GenericChannelEvent{
		GenericEvent: events.NewEvent(disgo, sequenceNumber),
		ChannelID:    channel.ID,
		Channel:      channel,
	}
	eventManager.Dispatch(genericChannelEvent)

	var genericGuildChannelEvent events.GenericGuildChannelEvent
	if channel.GuildID != nil {
		genericGuildChannelEvent = events.GenericGuildChannelEvent{
			GuildID:             *channel.GuildID,
			GenericChannelEvent: genericChannelEvent,
			GuildChannel: &api.GuildChannel{
				Channel: *channel,
			},
		}
		eventManager.Dispatch(genericGuildChannelEvent)

		eventManager.Dispatch(events.GuildChannelCreateEvent{
			GenericGuildChannelEvent: genericGuildChannelEvent,
		})
	}

	switch channel.Type {
	case api.ChannelTypeDM:
		dmChannel := disgo.EntityBuilder().CreateDMChannel(channel, api.CacheStrategyYes)

		genericDMChannelEvent := events.GenericDMChannelEvent{
			GenericChannelEvent: genericChannelEvent,
			DMChannel:           dmChannel,
		}
		eventManager.Dispatch(genericDMChannelEvent)

		eventManager.Dispatch(events.DMChannelCreateEvent{
			GenericDMChannelEvent: genericDMChannelEvent,
		})

	case api.ChannelTypeGroupDM:
		disgo.Logger().Warnf("ChannelTypeGroupDM received what the hell discord")

	case api.ChannelTypeText, api.ChannelTypeNews:
		textChannel := disgo.EntityBuilder().CreateTextChannel(channel, api.CacheStrategyYes)

		genericTextChannelEvent := events.GenericTextChannelEvent{
			GenericGuildChannelEvent: genericGuildChannelEvent,
			TextChannel:              textChannel,
		}
		eventManager.Dispatch(genericTextChannelEvent)

		eventManager.Dispatch(events.TextChannelCreateEvent{
			GenericTextChannelEvent: genericTextChannelEvent,
		})

	case api.ChannelTypeStore:
		storeChannel := disgo.EntityBuilder().CreateStoreChannel(channel, api.CacheStrategyYes)

		genericStoreChannelEvent := events.GenericStoreChannelEvent{
			GenericGuildChannelEvent: genericGuildChannelEvent,
			StoreChannel:             storeChannel,
		}
		eventManager.Dispatch(genericStoreChannelEvent)

		eventManager.Dispatch(events.StoreChannelCreateEvent{
			GenericStoreChannelEvent: genericStoreChannelEvent,
		})

	case api.ChannelTypeCategory:
		category := disgo.EntityBuilder().CreateCategory(channel, api.CacheStrategyYes)

		genericCategoryEvent := events.GenericCategoryEvent{
			GenericGuildChannelEvent: genericGuildChannelEvent,
			Category:                 category,
		}
		eventManager.Dispatch(genericCategoryEvent)

		eventManager.Dispatch(events.CategoryCreateEvent{
			GenericCategoryEvent: genericCategoryEvent,
		})

	case api.ChannelTypeVoice:
		voiceChannel := disgo.EntityBuilder().CreateVoiceChannel(channel, api.CacheStrategyYes)

		genericVoiceChannelEvent := events.GenericVoiceChannelEvent{
			GenericGuildChannelEvent: genericGuildChannelEvent,
			VoiceChannel:             voiceChannel,
		}
		eventManager.Dispatch(genericVoiceChannelEvent)

		eventManager.Dispatch(events.VoiceChannelCreateEvent{
			GenericVoiceChannelEvent: genericVoiceChannelEvent,
		})

	default:
		disgo.Logger().Warnf("unknown channel type received: %d", channel.Type)
	}

}
