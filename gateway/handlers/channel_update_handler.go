package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

// ChannelUpdateHandler handles api.GatewayEventChannelUpdate
type ChannelUpdateHandler struct{}

// EventType returns the api.GatewayGatewayEventType
func (h *ChannelUpdateHandler) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeChannelUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *ChannelUpdateHandler) New() interface{} {
	return discord.Channel{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *ChannelUpdateHandler) HandleGatewayEvent(disgo core.Disgo, eventManager core.EventManager, sequenceNumber int, i interface{}) {
	channel, ok := i.(discord.Channel)
	if !ok {
		return
	}

	oldCoreChannel := disgo.Cache().ChannelCache().GetChannelCopy(channel.ID)

	coreChannel := disgo.EntityBuilder().CreateChannel(channel, core.CacheStrategyYes)

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

		eventManager.Dispatch(&events.GuildChannelUpdateEvent{
			GenericGuildChannelEvent: genericGuildChannelEvent,
			OldGuildChannel:          oldCoreChannel.(core.GuildChannel),
		})
	}

	switch channel.Type {
	case discord.ChannelTypeDM:
		eventManager.Dispatch(&events.DMChannelUpdateEvent{
			GenericDMChannelEvent: &events.GenericDMChannelEvent{
				GenericChannelEvent: genericChannelEvent,
				DMChannel:           coreChannel.(core.DMChannel),
			},
			OldDMChannel: oldCoreChannel.(core.DMChannel),
		})

	case discord.ChannelTypeGroupDM:
		disgo.Logger().Warnf("ChannelTypeGroupDM received what the hell discord")

	case discord.ChannelTypeText:
		eventManager.Dispatch(&events.TextChannelUpdateEvent{
			GenericTextChannelEvent: &events.GenericTextChannelEvent{
				GenericGuildChannelEvent: genericGuildChannelEvent,
				TextChannel:              coreChannel.(core.TextChannel),
			},
			OldTextChannel: oldCoreChannel.(core.TextChannel),
		})

	case discord.ChannelTypeNews:
		eventManager.Dispatch(&events.NewsChannelUpdateEvent{
			GenericNewsChannelEvent: &events.GenericNewsChannelEvent{
				GenericGuildChannelEvent: genericGuildChannelEvent,
				NewsChannel:              coreChannel.(core.NewsChannel),
			},
			OldNewsChannel: oldCoreChannel.(core.NewsChannel),
		})

	case discord.ChannelTypeStore:
		eventManager.Dispatch(&events.StoreChannelUpdateEvent{
			GenericStoreChannelEvent: &events.GenericStoreChannelEvent{
				GenericGuildChannelEvent: genericGuildChannelEvent,
				StoreChannel:             coreChannel.(core.StoreChannel),
			},
			OldStoreChannel: oldCoreChannel.(core.StoreChannel),
		})

	case discord.ChannelTypeCategory:
		eventManager.Dispatch(&events.CategoryUpdateEvent{
			GenericCategoryEvent: &events.GenericCategoryEvent{
				GenericGuildChannelEvent: genericGuildChannelEvent,
				Category:                 coreChannel.(core.Category),
			},
			OldCategory: oldCoreChannel.(core.Category),
		})

	case discord.ChannelTypeVoice:
		eventManager.Dispatch(&events.VoiceChannelUpdateEvent{
			GenericVoiceChannelEvent: &events.GenericVoiceChannelEvent{
				GenericGuildChannelEvent: genericGuildChannelEvent,
				VoiceChannel:             coreChannel.(core.VoiceChannel),
			},
			OldVoiceChannel: oldCoreChannel.(core.VoiceChannel),
		})

	case discord.ChannelTypeStage:
		eventManager.Dispatch(&events.StageChannelUpdateEvent{
			GenericStageChannelEvent: &events.GenericStageChannelEvent{
				GenericGuildChannelEvent: genericGuildChannelEvent,
				StageChannel:             coreChannel.(core.StageChannel),
			},
			OldStageChannel: oldCoreChannel.(core.StageChannel),
		})

	default:
		disgo.Logger().Warnf("unknown channel type received: %d", channel.Type)
	}
}
