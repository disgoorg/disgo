package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

// gatewayHandlerGuildCreate handles core.GuildCreateGatewayEvent
type gatewayHandlerGuildCreate struct{}

// EventType returns the core.GatewayGatewayEventType
func (h *gatewayHandlerGuildCreate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildCreate) New() any {
	return &discord.GatewayGuild{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildCreate) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	gatewayGuild := *v.(*discord.GatewayGuild)

	wasUnready := client.Caches().Guilds().IsUnready(shardID, gatewayGuild.ID)
	wasUnavailable := client.Caches().Guilds().IsUnavailable(gatewayGuild.ID)

	client.Caches().Guilds().Put(gatewayGuild.ID, gatewayGuild.Guild)

	for _, channel := range gatewayGuild.Channels {
		channel = discord.ApplyGuildIDToChannel(channel, gatewayGuild.ID) // populate unset field
		client.Caches().Channels().Put(channel.ID(), discord.ApplyGuildIDToChannel(channel, gatewayGuild.ID))
	}

	for _, thread := range gatewayGuild.Threads {
		thread = discord.ApplyGuildIDToThread(thread, gatewayGuild.ID) // populate unset field
		client.Caches().Channels().Put(thread.ID(), discord.ApplyGuildIDToThread(thread, gatewayGuild.ID))
	}

	for _, role := range gatewayGuild.Roles {
		client.Caches().Roles().Put(gatewayGuild.ID, role.ID, role)
	}

	for _, member := range gatewayGuild.Members {
		member.GuildID = gatewayGuild.ID // populate unset field
		client.Caches().Members().Put(gatewayGuild.ID, member.User.ID, member)
	}

	for _, voiceState := range gatewayGuild.VoiceStates {
		voiceState.GuildID = gatewayGuild.ID // populate unset field
		client.Caches().VoiceStates().Put(voiceState.GuildID, voiceState.UserID, voiceState)
	}

	for _, emoji := range gatewayGuild.Emojis {
		client.Caches().Emojis().Put(gatewayGuild.ID, emoji.ID, emoji)
	}

	for _, sticker := range gatewayGuild.Stickers {
		client.Caches().Stickers().Put(gatewayGuild.ID, sticker.ID, sticker)
	}

	for _, stageInstance := range gatewayGuild.StageInstances {
		client.Caches().StageInstances().Put(gatewayGuild.ID, stageInstance.ID, stageInstance)
	}

	for _, guildScheduledEvent := range gatewayGuild.GuildScheduledEvents {
		client.Caches().GuildScheduledEvents().Put(gatewayGuild.ID, guildScheduledEvent.ID, guildScheduledEvent)
	}

	for _, presence := range gatewayGuild.Presences {
		presence.GuildID = gatewayGuild.ID // populate unset field
		client.Caches().Presences().Put(gatewayGuild.ID, presence.PresenceUser.ID, presence)
	}

	genericGuildEvent := &events.GenericGuild{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		GuildID:      gatewayGuild.ID,
		Guild:        gatewayGuild.Guild,
	}

	if wasUnready {
		client.Caches().Guilds().SetReady(shardID, gatewayGuild.ID)
		client.EventManager().DispatchEvent(&events.GuildReady{
			GenericGuild: genericGuildEvent,
		})
		if len(client.Caches().Guilds().UnreadyGuilds(shardID)) == 0 {
			client.EventManager().DispatchEvent(&events.GuildsReady{
				GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			})
		}
		if client.MemberChunkingManager().MemberChunkingFilter()(gatewayGuild.ID) {
			go func() {
				if _, err := client.MemberChunkingManager().RequestMembersWithQuery(gatewayGuild.ID, "", 0); err != nil {
					client.Logger().Error("failed to chunk guild on guild_create. error: ", err)
				}
			}()
		}

	}
	if wasUnavailable {
		client.Caches().Guilds().SetAvailable(gatewayGuild.ID)
		client.EventManager().DispatchEvent(&events.GuildAvailable{
			GenericGuild: genericGuildEvent,
		})
	} else {
		client.EventManager().DispatchEvent(&events.GuildJoin{
			GenericGuild: genericGuildEvent,
		})
	}
}
