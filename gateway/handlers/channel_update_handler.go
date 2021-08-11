package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
)

// ChannelUpdateHandler handles api.GatewayEventChannelUpdate
type ChannelUpdateHandler struct{}

// EventType returns the api.GatewayEventType
func (h *ChannelUpdateHandler) EventType() gateway.EventType {
	return gateway.EventTypeChannelUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *ChannelUpdateHandler) New() interface{} {
	return &discord.Channel{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *ChannelUpdateHandler) HandleGatewayEvent(disgo core.Disgo, eventManager core.EventManager, sequenceNumber int, i interface{}) {
	channel, ok := i.(*discord.Channel)
	if !ok {
		return
	}

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
			GuildChannel: &core.GuildChannel{
				Channel: *channel,
			},
		}

		eventManager.Dispatch(&events.GuildChannelUpdateEvent{
			GenericGuildChannelEvent: genericGuildChannelEvent,
			OldGuildChannel:          disgo.Cache().GuildChannel(channel.ID),
		})
	}

	switch channel.Type {
	case discord.ChannelTypeDM:
		oldDMChannel := disgo.Cache().DMChannel(channel.ID)
		if oldDMChannel != nil {
			oldDMChannel = &*oldDMChannel
		}

		eventManager.Dispatch(&events.DMChannelUpdateEvent{
			GenericDMChannelEvent: &events.GenericDMChannelEvent{
				GenericChannelEvent: genericChannelEvent,
				DMChannel:           disgo.EntityBuilder().CreateDMChannel(channel, core.CacheStrategyYes),
			},
			OldDMChannel: oldDMChannel,
		})

	case discord.ChannelTypeGroupDM:
		disgo.Logger().Warnf("ChannelTypeGroupDM received what the hell discord")

	case discord.ChannelTypeText, discord.ChannelTypeNews:
		oldTextChannel := disgo.Cache().TextChannel(channel.ID)
		if oldTextChannel != nil {
			oldTextChannel = &*oldTextChannel
		}

		eventManager.Dispatch(&events.TextChannelUpdateEvent{
			GenericTextChannelEvent: &events.GenericTextChannelEvent{
				GenericGuildChannelEvent: genericGuildChannelEvent,
				TextChannel:              disgo.EntityBuilder().CreateTextChannel(channel, core.CacheStrategyYes),
			},
			OldTextChannel: oldTextChannel,
		})

	case discord.ChannelTypeStore:
		oldStoreChannel := disgo.Cache().StoreChannel(channel.ID)
		if oldStoreChannel != nil {
			oldStoreChannel = &*oldStoreChannel
		}

		eventManager.Dispatch(&events.StoreChannelUpdateEvent{
			GenericStoreChannelEvent: &events.GenericStoreChannelEvent{
				GenericGuildChannelEvent: genericGuildChannelEvent,
				StoreChannel:             disgo.EntityBuilder().CreateStoreChannel(channel, core.CacheStrategyYes),
			},
			OldStoreChannel: oldStoreChannel,
		})

	case discord.ChannelTypeCategory:
		oldCategory := disgo.Cache().Category(channel.ID)
		if oldCategory != nil {
			oldCategory = &*oldCategory
		}

		eventManager.Dispatch(&events.CategoryUpdateEvent{
			GenericCategoryEvent: &events.GenericCategoryEvent{
				GenericGuildChannelEvent: genericGuildChannelEvent,
				Category:                 disgo.EntityBuilder().CreateCategory(channel, core.CacheStrategyYes),
			},
			OldCategory: oldCategory,
		})

	case discord.ChannelTypeVoice:
		oldVoiceChannel := disgo.Cache().VoiceChannel(channel.ID)
		if oldVoiceChannel != nil {
			oldVoiceChannel = &*oldVoiceChannel
		}

		eventManager.Dispatch(&events.VoiceChannelUpdateEvent{
			GenericVoiceChannelEvent: &events.GenericVoiceChannelEvent{
				GenericGuildChannelEvent: genericGuildChannelEvent,
				VoiceChannel:             disgo.EntityBuilder().CreateVoiceChannel(channel, core.CacheStrategyYes),
			},
			OldVoiceChannel: oldVoiceChannel,
		})

	default:
		disgo.Logger().Warnf("unknown channel type received: %d", channel.Type)
	}
}
