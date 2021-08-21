package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
)

// ChannelDeleteHandler handles api.GatewayEventChannelDelete
type ChannelDeleteHandler struct{}

// EventType returns the api.GatewayEventType
func (h *ChannelDeleteHandler) EventType() gateway.EventType {
	return gateway.EventTypeChannelDelete
}

// New constructs a new payload receiver for the raw gateway event
func (h *ChannelDeleteHandler) New() interface{} {
	return discord.Channel{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *ChannelDeleteHandler) HandleGatewayEvent(disgo core.Disgo, eventManager core.EventManager, sequenceNumber int, i interface{}) {
	deletedChannel, ok := i.(discord.Channel)
	if !ok {
		return
	}

	channel := disgo.EntityBuilder().CreateChannel(deletedChannel, core.CacheStrategyNo)

	genericChannelEvent := &events.GenericChannelEvent{
		GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
		ChannelID:    channel.ID(),
		Channel:      channel,
	}

	var genericGuildChannelEvent *events.GenericGuildChannelEvent
	if channel.GuildID != nil {
		genericGuildChannelEvent = &events.GenericGuildChannelEvent{
			GenericChannelEvent: genericChannelEvent,
			GuildID:             *channel.GuildID,
			GuildChannel: &discord.GuildChannel{
				Channel: *channel,
			},
		}

		eventManager.Dispatch(&events.GuildChannelDeleteEvent{
			GenericGuildChannelEvent: genericGuildChannelEvent,
		})
	}

	switch channel.Type {
	case discord.ChannelTypeDM:
		disgo.Cache().UncacheDMChannel(channel.ID)

		eventManager.Dispatch(&events.DMChannelCreateEvent{
			GenericDMChannelEvent: &events.GenericDMChannelEvent{
				GenericChannelEvent: genericChannelEvent,
				DMChannel:           disgo.EntityBuilder().CreateDMChannel(channel, core.CacheStrategyNo),
			},
		})

	case discord.ChannelTypeGroupDM:
		disgo.Logger().Warnf("ChannelTypeGroupDM received what the hell discord")

	case discord.ChannelTypeText, discord.ChannelTypeNews:
		disgo.Cache().UncacheTextChannel(*channel.GuildID, channel.ID)

		eventManager.Dispatch(&events.TextChannelCreateEvent{
			GenericTextChannelEvent: &events.GenericTextChannelEvent{
				GenericGuildChannelEvent: genericGuildChannelEvent,
				TextChannel:              disgo.EntityBuilder().CreateTextChannel(channel, core.CacheStrategyNo),
			},
		})

	case discord.ChannelTypeStore:
		disgo.Cache().UncacheStoreChannel(*channel.GuildID, channel.ID)

		eventManager.Dispatch(&events.StoreChannelCreateEvent{
			GenericStoreChannelEvent: &events.GenericStoreChannelEvent{
				GenericGuildChannelEvent: genericGuildChannelEvent,
				StoreChannel:             disgo.EntityBuilder().CreateStoreChannel(channel, core.CacheStrategyNo),
			},
		})

	case discord.ChannelTypeCategory:
		disgo.Cache().UncacheCategory(*channel.GuildID, channel.ID)

		eventManager.Dispatch(&events.CategoryCreateEvent{
			GenericCategoryEvent: &events.GenericCategoryEvent{
				GenericGuildChannelEvent: genericGuildChannelEvent,
				Category:                 disgo.EntityBuilder().CreateCategory(channel, core.CacheStrategyNo),
			},
		})

	case discord.ChannelTypeVoice:
		disgo.Cache().UncacheVoiceChannel(*channel.GuildID, channel.ID)

		eventManager.Dispatch(&events.VoiceChannelCreateEvent{
			GenericVoiceChannelEvent: &events.GenericVoiceChannelEvent{
				GenericGuildChannelEvent: genericGuildChannelEvent,
				VoiceChannel:             disgo.EntityBuilder().CreateVoiceChannel(channel, core.CacheStrategyNo),
			},
		})

	default:
		disgo.Logger().Warnf("unknown channel type received: %d", channel.Type)
	}
}
