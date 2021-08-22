package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

// ChannelDeleteHandler handles api.GatewayEventChannelDelete
type ChannelDeleteHandler struct{}

// EventType returns the api.GatewayGatewayEventType
func (h *ChannelDeleteHandler) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeChannelDelete
}

// New constructs a new payload receiver for the raw gateway event
func (h *ChannelDeleteHandler) New() interface{} {
	return discord.Channel{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *ChannelDeleteHandler) HandleGatewayEvent(disgo core.Disgo, eventManager core.EventManager, sequenceNumber int, i interface{}) {
	channel, ok := i.(discord.Channel)
	if !ok {
		return
	}

	coreChannel := disgo.EntityBuilder().CreateChannel(channel, core.CacheStrategyNo)

	genericChannelEvent := &events.GenericChannelEvent{
		GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
		ChannelID:    channel.ID,
		Channel:      coreChannel,
	}

	var genericGuildChannelEvent *events.GenericGuildChannelEvent
	if channel.GuildID != nil {
		genericGuildChannelEvent = &events.GenericGuildChannelEvent{
			GenericChannelEvent: genericChannelEvent,
			GuildID:             *channel.GuildID,
			GuildChannel:        coreChannel.(core.GuildChannel),
		}

		eventManager.Dispatch(&events.GuildChannelDeleteEvent{
			GenericGuildChannelEvent: genericGuildChannelEvent,
		})
	}

	switch channel.Type {
	case discord.ChannelTypeDM:
		disgo.Cache().DMChannelCache().Uncache(channel.ID)

		eventManager.Dispatch(&events.DMChannelCreateEvent{
			GenericDMChannelEvent: &events.GenericDMChannelEvent{
				GenericChannelEvent: genericChannelEvent,
				DMChannel:           coreChannel.(core.DMChannel),
			},
		})

	case discord.ChannelTypeGroupDM:
		disgo.Logger().Warnf("ChannelTypeGroupDM received what the hell discord")

	case discord.ChannelTypeText:
		disgo.Cache().TextChannelCache().Uncache(*channel.GuildID, channel.ID)

		eventManager.Dispatch(&events.TextChannelCreateEvent{
			GenericTextChannelEvent: &events.GenericTextChannelEvent{
				GenericGuildChannelEvent: genericGuildChannelEvent,
				TextChannel:              coreChannel.(core.TextChannel),
			},
		})

	case discord.ChannelTypeNews:
		disgo.Cache().NewsChannelCache().Uncache(*channel.GuildID, channel.ID)

		eventManager.Dispatch(&events.NewsChannelCreateEvent{
			GenericNewsChannelEvent: &events.GenericNewsChannelEvent{
				GenericGuildChannelEvent: genericGuildChannelEvent,
				NewsChannel:              coreChannel.(core.NewsChannel),
			},
		})

	case discord.ChannelTypeStore:
		disgo.Cache().StoreChannelCache().Uncache(*channel.GuildID, channel.ID)

		eventManager.Dispatch(&events.StoreChannelCreateEvent{
			GenericStoreChannelEvent: &events.GenericStoreChannelEvent{
				GenericGuildChannelEvent: genericGuildChannelEvent,
				StoreChannel:             coreChannel.(core.StoreChannel),
			},
		})

	case discord.ChannelTypeCategory:
		disgo.Cache().CategoryCache().Uncache(*channel.GuildID, channel.ID)

		eventManager.Dispatch(&events.CategoryCreateEvent{
			GenericCategoryEvent: &events.GenericCategoryEvent{
				GenericGuildChannelEvent: genericGuildChannelEvent,
				Category:                 coreChannel.(core.Category),
			},
		})

	case discord.ChannelTypeVoice:
		disgo.Cache().VoiceChannelCache().Uncache(*channel.GuildID, channel.ID)

		eventManager.Dispatch(&events.VoiceChannelCreateEvent{
			GenericVoiceChannelEvent: &events.GenericVoiceChannelEvent{
				GenericGuildChannelEvent: genericGuildChannelEvent,
				VoiceChannel:             coreChannel.(core.VoiceChannel),
			},
		})

	default:
		disgo.Logger().Warnf("unknown channel type received: %d", channel.Type)
	}
}
