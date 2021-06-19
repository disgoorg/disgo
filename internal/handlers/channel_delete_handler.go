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
	return &api.ChannelImpl{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h ChannelDeleteHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
	channel, ok := i.(*api.ChannelImpl)
	if !ok {
		return
	}

	channel.Disgo = disgo

	genericChannelEvent := events.GenericChannelEvent{
		GenericEvent: events.NewEvent(disgo, sequenceNumber),
		ChannelID:    channel.ID(),
		Channel:      channel,
	}
	eventManager.Dispatch(genericChannelEvent)

	var genericGuildChannelEvent events.GenericGuildChannelEvent
	if channel.GuildID != nil {
		genericGuildChannelEvent = events.GenericGuildChannelEvent{
			GenericChannelEvent: genericChannelEvent,
			GuildID:             *channel.GuildID(),
			GuildChannel: &api.GuildChannel{
				Channel: *channel,
			},
		}
		eventManager.Dispatch(genericGuildChannelEvent)

		eventManager.Dispatch(events.GuildChannelDeleteEvent{
			GenericGuildChannelEvent: genericGuildChannelEvent,
		})
	}

	switch channel.Type() {
	case api.ChannelTypeDM:
		disgo.Cache().UncacheDMChannel(channel.ID())

		genericDMChannelEvent := events.GenericDMChannelEvent{
			GenericChannelEvent: genericChannelEvent,
			DMChannel:           disgo.EntityBuilder().CreateDMChannel(channel, api.CacheStrategyNo),
		}
		eventManager.Dispatch(genericDMChannelEvent)

		eventManager.Dispatch(events.DMChannelDeleteEvent{
			GenericDMChannelEvent: genericDMChannelEvent,
		})

	case api.ChannelTypeGroupDM:
		disgo.Logger().Warnf("ChannelTypeGroupDM received what the hell discord")

	case api.ChannelTypeText, api.ChannelTypeNews:
		disgo.Cache().UncacheTextChannel(channel.GuildID(), channel.ID())

		genericTextChannelEvent := events.GenericTextChannelEvent{
			GenericGuildChannelEvent: genericGuildChannelEvent,
			TextChannel:              disgo.EntityBuilder().CreateTextChannel(channel, api.CacheStrategyNo),
		}
		eventManager.Dispatch(genericTextChannelEvent)

		eventManager.Dispatch(events.TextChannelDeleteEvent{
			GenericTextChannelEvent: genericTextChannelEvent,
		})

	case api.ChannelTypeStore:
		disgo.Cache().UncacheStoreChannel(channel.GuildID(), channel.ID())

		genericStoreChannelEvent := events.GenericStoreChannelEvent{
			GenericGuildChannelEvent: genericGuildChannelEvent,
			StoreChannel:             disgo.EntityBuilder().CreateStoreChannel(channel, api.CacheStrategyNo),
		}
		eventManager.Dispatch(genericStoreChannelEvent)

		eventManager.Dispatch(events.StoreChannelDeleteEvent{
			GenericStoreChannelEvent: genericStoreChannelEvent,
		})

	case api.ChannelTypeCategory:
		disgo.Cache().UncacheCategory(channel.GuildID(), channel.ID())

		genericCategoryEvent := events.GenericCategoryEvent{
			GenericGuildChannelEvent: genericGuildChannelEvent,
			Category:                 disgo.EntityBuilder().CreateCategory(channel, api.CacheStrategyNo),
		}
		eventManager.Dispatch(genericCategoryEvent)

		eventManager.Dispatch(events.CategoryDeleteEvent{
			GenericCategoryEvent: genericCategoryEvent,
		})

	case api.ChannelTypeVoice:
		disgo.Cache().UncacheVoiceChannel(channel.GuildID(), channel.ID())

		genericVoiceChannelEvent := events.GenericVoiceChannelEvent{
			GenericGuildChannelEvent: genericGuildChannelEvent,
			VoiceChannel:             disgo.EntityBuilder().CreateVoiceChannel(channel, api.CacheStrategyNo),
		}
		eventManager.Dispatch(genericVoiceChannelEvent)

		eventManager.Dispatch(events.VoiceChannelDeleteEvent{
			GenericVoiceChannelEvent: genericVoiceChannelEvent,
		})

	default:
		disgo.Logger().Warnf("unknown channel type received: %d", channel.Type)
	}
}
