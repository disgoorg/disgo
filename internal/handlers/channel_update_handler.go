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

	channel.Disgo = disgo

	genericChannelEvent := &events.GenericChannelEvent{
		GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
		ChannelID:    channel.ID,
		Channel:      channel,
	}

	var genericGuildChannelEvent *events.GenericGuildChannelEvent
	if channel.GuildID != nil {
		genericGuildChannelEvent = &events.GenericGuildChannelEvent{
			GenericChannelEvent: genericChannelEvent,
			GuildID:             *channel.GuildID,
			GuildChannel: &api.GuildChannel{
				Channel: *channel,
			},
		}

		eventManager.Dispatch(&events.GuildChannelUpdateEvent{
			GenericGuildChannelEvent: genericGuildChannelEvent,
			OldGuildChannel:          disgo.Cache().GuildChannel(channel.ID),
		})
	}

	switch channel.Type {
	case api.ChannelTypeDM:
		oldDMChannel := disgo.Cache().DMChannel(channel.ID)
		if oldDMChannel != nil {
			oldDMChannel = &*oldDMChannel
		}

		eventManager.Dispatch(&events.DMChannelUpdateEvent{
			GenericDMChannelEvent: &events.GenericDMChannelEvent{
				GenericChannelEvent: genericChannelEvent,
				DMChannel:           disgo.EntityBuilder().CreateDMChannel(channel, api.CacheStrategyYes),
			},
			OldDMChannel: oldDMChannel,
		})

	case api.ChannelTypeGroupDM:
		disgo.Logger().Warnf("ChannelTypeGroupDM received what the hell discord")

	case api.ChannelTypeText, api.ChannelTypeNews:
		oldTextChannel := disgo.Cache().TextChannel(channel.ID)
		if oldTextChannel != nil {
			oldTextChannel = &*oldTextChannel
		}

		eventManager.Dispatch(&events.TextChannelUpdateEvent{
			GenericTextChannelEvent: &events.GenericTextChannelEvent{
				GenericGuildChannelEvent: genericGuildChannelEvent,
				TextChannel:              disgo.EntityBuilder().CreateTextChannel(channel, api.CacheStrategyYes),
			},
			OldTextChannel: oldTextChannel,
		})

	case api.ChannelTypeStore:
		oldStoreChannel := disgo.Cache().StoreChannel(channel.ID)
		if oldStoreChannel != nil {
			oldStoreChannel = &*oldStoreChannel
		}

		eventManager.Dispatch(&events.StoreChannelUpdateEvent{
			GenericStoreChannelEvent: &events.GenericStoreChannelEvent{
				GenericGuildChannelEvent: genericGuildChannelEvent,
				StoreChannel:             disgo.EntityBuilder().CreateStoreChannel(channel, api.CacheStrategyYes),
			},
			OldStoreChannel: oldStoreChannel,
		})

	case api.ChannelTypeCategory:
		oldCategory := disgo.Cache().Category(channel.ID)
		if oldCategory != nil {
			oldCategory = &*oldCategory
		}

		eventManager.Dispatch(&events.CategoryUpdateEvent{
			GenericCategoryEvent: &events.GenericCategoryEvent{
				GenericGuildChannelEvent: genericGuildChannelEvent,
				Category:                 disgo.EntityBuilder().CreateCategory(channel, api.CacheStrategyYes),
			},
			OldCategory: oldCategory,
		})

	case api.ChannelTypeVoice:
		oldVoiceChannel := disgo.Cache().VoiceChannel(channel.ID)
		if oldVoiceChannel != nil {
			oldVoiceChannel = &*oldVoiceChannel
		}

		eventManager.Dispatch(&events.VoiceChannelUpdateEvent{
			GenericVoiceChannelEvent: &events.GenericVoiceChannelEvent{
				GenericGuildChannelEvent: genericGuildChannelEvent,
				VoiceChannel:             disgo.EntityBuilder().CreateVoiceChannel(channel, api.CacheStrategyYes),
			},
			OldVoiceChannel: oldVoiceChannel,
		})

	default:
		disgo.Logger().Warnf("unknown channel type received: %d", channel.Type)
	}
}
