package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/snowflake"
)

// gatewayHandlerGuildDelete handles discord.GatewayEventTypeGuildDelete
type gatewayHandlerGuildDelete struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerGuildDelete) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildDelete
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildDelete) New() any {
	return &discord.UnavailableGuild{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildDelete) HandleGatewayEvent(bot core.Bot, sequenceNumber discord.GatewaySequence, v any) {
	unavailableGuild := *v.(*discord.UnavailableGuild)

	guild, _ := bot.Caches().Guilds().Remove(unavailableGuild.ID)
	bot.Caches().VoiceStates().RemoveAll(unavailableGuild.ID)
	bot.Caches().Presences().RemoveAll(unavailableGuild.ID)
	bot.Caches().Channels().RemoveIf(func(channel discord.Channel) bool {
		if guildChannel, ok := channel.(discord.GuildChannel); ok {
			return guildChannel.GuildID() == unavailableGuild.ID
		}
		return false
	})
	bot.Caches().Emojis().RemoveAll(unavailableGuild.ID)
	bot.Caches().Stickers().RemoveAll(unavailableGuild.ID)
	bot.Caches().Roles().RemoveAll(unavailableGuild.ID)
	bot.Caches().StageInstances().RemoveAll(unavailableGuild.ID)
	bot.Caches().ThreadMembers().RemoveIf(func(groupID snowflake.Snowflake, threadMember discord.ThreadMember) bool {
		// TODO: figure out how to remove thread members from cache via guild id
		return false
	})
	bot.Caches().Messages().RemoveIf(func(channelID snowflake.Snowflake, message discord.Message) bool {
		return message.GuildID != nil && *message.GuildID == unavailableGuild.ID
	})

	if unavailableGuild.Unavailable {
		bot.Caches().Guilds().SetUnavailable(unavailableGuild.ID)
	}

	genericGuildEvent := &events.GenericGuildEvent{
		GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
		GuildID:      unavailableGuild.ID,
		Guild:        guild,
	}

	if unavailableGuild.Unavailable {
		bot.EventManager().Dispatch(&events.GuildUnavailableEvent{
			GenericGuildEvent: genericGuildEvent,
		})
	} else {
		bot.EventManager().Dispatch(&events.GuildLeaveEvent{
			GenericGuildEvent: genericGuildEvent,
		})
	}
}
