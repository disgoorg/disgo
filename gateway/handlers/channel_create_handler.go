package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
)

// ChannelCreateHandler handles api.GatewayEventChannelCreate
type ChannelCreateHandler struct{}

// EventType returns the api.GatewayEventType
func (h *ChannelCreateHandler) EventType() gateway.EventType {
	return gateway.EventTypeChannelCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h *ChannelCreateHandler) New() interface{} {
	return discord.Channel{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *ChannelCreateHandler) HandleGatewayEvent(disgo core.Disgo, eventManager core.EventManager, sequenceNumber int, i interface{}) {
	channel, ok := i.(discord.Channel)
	if !ok {
		return
	}

	genericChannelEvent := &events.GenericChannelEvent{
		GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
		ChannelID:    channel.ID,
	}

	var genericGuildChannelEvent *events.GenericGuildChannelEvent
	if channel.GuildID != nil {
		genericGuildChannelEvent = &events.GenericGuildChannelEvent{
			GuildID:             *channel.GuildID,
			GenericChannelEvent: genericChannelEvent,
		}

		eventManager.Dispatch(&events.GuildChannelCreateEvent{
			GenericGuildChannelEvent: genericGuildChannelEvent,
		})
	}

	switch channel.Type {
	case discord.ChannelTypeDM:
		eventManager.Dispatch(&events.DMChannelCreateEvent{
			GenericDMChannelEvent: &events.GenericDMChannelEvent{
				GenericChannelEvent: genericChannelEvent,
				DMChannel:           disgo.EntityBuilder().CreateChannel(channel, core.CacheStrategyYes).(core.DMChannel),
			},
		})

	case discord.ChannelTypeGroupDM:
		disgo.Logger().Warnf("ChannelTypeGroupDM received what the hell discord")

	case discord.ChannelTypeText, discord.ChannelTypeNews:
		eventManager.Dispatch(&events.TextChannelCreateEvent{
			GenericTextChannelEvent: &events.GenericTextChannelEvent{
				GenericGuildChannelEvent: genericGuildChannelEvent,
				TextChannel:              disgo.EntityBuilder().CreateChannel(channel, core.CacheStrategyYes).(core.TextChannel),
			},
		})

	case discord.ChannelTypeStore:
		eventManager.Dispatch(&events.StoreChannelCreateEvent{
			GenericStoreChannelEvent: &events.GenericStoreChannelEvent{
				GenericGuildChannelEvent: genericGuildChannelEvent,
				StoreChannel:             disgo.EntityBuilder().CreateChannel(channel, core.CacheStrategyYes).(core.StageChannel),
			},
		})

	case discord.ChannelTypeCategory:
		eventManager.Dispatch(&events.CategoryCreateEvent{
			GenericCategoryEvent: &events.GenericCategoryEvent{
				GenericGuildChannelEvent: genericGuildChannelEvent,
				Category:                 disgo.EntityBuilder().CreateChannel(channel, core.CacheStrategyYes).(core.Category),
			},
		})

	case discord.ChannelTypeVoice:
		eventManager.Dispatch(&events.VoiceChannelCreateEvent{
			GenericVoiceChannelEvent: &events.GenericVoiceChannelEvent{
				GenericGuildChannelEvent: genericGuildChannelEvent,
				VoiceChannel:             disgo.EntityBuilder().CreateChannel(channel, core.CacheStrategyYes).(core.VoiceChannel),
			},
		})

	default:
		disgo.Logger().Warnf("unknown channel type received: %d", channel.Type)
	}

}
