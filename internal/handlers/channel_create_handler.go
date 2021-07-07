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

	genericChannelEvent := &events.GenericChannelEvent{
		GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
		ChannelID:    channel.ID,
		Channel:      channel,
	}

	var genericGuildChannelEvent *events.GenericGuildChannelEvent
	if channel.GuildID != nil {
		genericGuildChannelEvent = &events.GenericGuildChannelEvent{
			GuildID:             *channel.GuildID,
			GenericChannelEvent: genericChannelEvent,
			GuildChannel: &api.GuildChannel{
				Channel: *channel,
			},
		}

		eventManager.Dispatch(&events.GuildChannelCreateEvent{
			GenericGuildChannelEvent: genericGuildChannelEvent,
		})
	}

	switch channel.Type {
	case api.ChannelTypeDM:
		eventManager.Dispatch(&events.DMChannelCreateEvent{
			GenericDMChannelEvent: &events.GenericDMChannelEvent{
				GenericChannelEvent: genericChannelEvent,
				DMChannel:           disgo.EntityBuilder().CreateDMChannel(channel, api.CacheStrategyYes),
			},
		})

	case api.ChannelTypeGroupDM:
		disgo.Logger().Warnf("ChannelTypeGroupDM received what the hell discord")

	case api.ChannelTypeText, api.ChannelTypeNews:
		eventManager.Dispatch(&events.TextChannelCreateEvent{
			GenericTextChannelEvent: &events.GenericTextChannelEvent{
				GenericGuildChannelEvent: genericGuildChannelEvent,
				TextChannel:              disgo.EntityBuilder().CreateTextChannel(channel, api.CacheStrategyYes),
			},
		})

	case api.ChannelTypeStore:
		eventManager.Dispatch(&events.StoreChannelCreateEvent{
			GenericStoreChannelEvent: &events.GenericStoreChannelEvent{
				GenericGuildChannelEvent: genericGuildChannelEvent,
				StoreChannel:             disgo.EntityBuilder().CreateStoreChannel(channel, api.CacheStrategyYes),
			},
		})

	case api.ChannelTypeCategory:
		eventManager.Dispatch(&events.CategoryCreateEvent{
			GenericCategoryEvent: &events.GenericCategoryEvent{
				GenericGuildChannelEvent: genericGuildChannelEvent,
				Category:                 disgo.EntityBuilder().CreateCategory(channel, api.CacheStrategyYes),
			},
		})

	case api.ChannelTypeVoice:
		eventManager.Dispatch(&events.VoiceChannelCreateEvent{
			GenericVoiceChannelEvent: &events.GenericVoiceChannelEvent{
				GenericGuildChannelEvent: genericGuildChannelEvent,
				VoiceChannel:             disgo.EntityBuilder().CreateVoiceChannel(channel, api.CacheStrategyYes),
			},
		})

	default:
		disgo.Logger().Warnf("unknown channel type received: %d", channel.Type)
	}

}
