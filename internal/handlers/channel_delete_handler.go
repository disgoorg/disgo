package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

// ChannelDeleteHandler handles api.GatewayEventChannelDelete
type ChannelDeleteHandler struct{}

// Event returns the raw gateway event Event
func (h *ChannelDeleteHandler) Event() api.GatewayEventType {
	return api.GatewayEventChannelDelete
}

// New constructs a new payload receiver for the raw gateway event
func (h *ChannelDeleteHandler) New() interface{} {
	return &api.Channel{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *ChannelDeleteHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
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

		eventManager.Dispatch(&events.GuildChannelDeleteEvent{
			GenericGuildChannelEvent: genericGuildChannelEvent,
		})
	}

	switch channel.Type {
	case api.ChannelTypeDM:
		disgo.Cache().UncacheDMChannel(channel.ID)

		eventManager.Dispatch(&events.DMChannelCreateEvent{
			GenericDMChannelEvent: &events.GenericDMChannelEvent{
				GenericChannelEvent: genericChannelEvent,
				DMChannel:           disgo.EntityBuilder().CreateDMChannel(channel, api.CacheStrategyNo),
			},
		})

	case api.ChannelTypeGroupDM:
		disgo.Logger().Warnf("ChannelTypeGroupDM received what the hell discord")

	case api.ChannelTypeText, api.ChannelTypeNews:
		disgo.Cache().UncacheTextChannel(*channel.GuildID, channel.ID)

		eventManager.Dispatch(&events.TextChannelCreateEvent{
			GenericTextChannelEvent: &events.GenericTextChannelEvent{
				GenericGuildChannelEvent: genericGuildChannelEvent,
				TextChannel:              disgo.EntityBuilder().CreateTextChannel(channel, api.CacheStrategyNo),
			},
		})

	case api.ChannelTypeStore:
		disgo.Cache().UncacheStoreChannel(*channel.GuildID, channel.ID)

		eventManager.Dispatch(&events.StoreChannelCreateEvent{
			GenericStoreChannelEvent: &events.GenericStoreChannelEvent{
				GenericGuildChannelEvent: genericGuildChannelEvent,
				StoreChannel:             disgo.EntityBuilder().CreateStoreChannel(channel, api.CacheStrategyNo),
			},
		})

	case api.ChannelTypeCategory:
		disgo.Cache().UncacheCategory(*channel.GuildID, channel.ID)

		eventManager.Dispatch(&events.CategoryCreateEvent{
			GenericCategoryEvent: &events.GenericCategoryEvent{
				GenericGuildChannelEvent: genericGuildChannelEvent,
				Category:                 disgo.EntityBuilder().CreateCategory(channel, api.CacheStrategyNo),
			},
		})

	case api.ChannelTypeVoice:
		disgo.Cache().UncacheVoiceChannel(*channel.GuildID, channel.ID)

		eventManager.Dispatch(&events.VoiceChannelCreateEvent{
			GenericVoiceChannelEvent: &events.GenericVoiceChannelEvent{
				GenericGuildChannelEvent: genericGuildChannelEvent,
				VoiceChannel:             disgo.EntityBuilder().CreateVoiceChannel(channel, api.CacheStrategyNo),
			},
		})

	default:
		disgo.Logger().Warnf("unknown channel type received: %d", channel.Type)
	}
}
