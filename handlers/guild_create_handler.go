package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerGuildCreate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildCreate) {
	wasUnready := client.Caches().Guilds().IsUnready(shardID, event.ID)
	wasUnavailable := client.Caches().Guilds().IsUnavailable(event.ID)

	client.Caches().Guilds().Put(event.ID, event.Guild)

	for _, channel := range event.Channels {
		channel = discord.ApplyGuildIDToChannel(channel, event.ID) // populate unset field
		client.Caches().Channels().Put(channel.ID(), discord.ApplyGuildIDToChannel(channel, event.ID))
	}

	for _, thread := range event.Threads {
		thread = discord.ApplyGuildIDToThread(thread, event.ID) // populate unset field
		client.Caches().Channels().Put(thread.ID(), discord.ApplyGuildIDToThread(thread, event.ID))
	}

	for _, role := range event.Roles {
		client.Caches().Roles().Put(event.ID, role.ID, role)
	}

	for _, member := range event.Members {
		member.GuildID = event.ID // populate unset field
		client.Caches().Members().Put(event.ID, member.User.ID, member)
	}

	for _, voiceState := range event.VoiceStates {
		voiceState.GuildID = event.ID // populate unset field
		client.Caches().VoiceStates().Put(voiceState.GuildID, voiceState.UserID, voiceState)
	}

	for _, emoji := range event.Emojis {
		client.Caches().Emojis().Put(event.ID, emoji.ID, emoji)
	}

	for _, sticker := range event.Stickers {
		client.Caches().Stickers().Put(event.ID, sticker.ID, sticker)
	}

	for _, stageInstance := range event.StageInstances {
		client.Caches().StageInstances().Put(event.ID, stageInstance.ID, stageInstance)
	}

	for _, guildScheduledEvent := range event.GuildScheduledEvents {
		client.Caches().GuildScheduledEvents().Put(event.ID, guildScheduledEvent.ID, guildScheduledEvent)
	}

	for _, presence := range event.Presences {
		presence.GuildID = event.ID // populate unset field
		client.Caches().Presences().Put(event.ID, presence.PresenceUser.ID, presence)
	}

	genericGuildEvent := &events.GenericGuild{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		GuildID:      event.ID,
		Guild:        event.Guild,
	}

	if wasUnready {
		client.Caches().Guilds().SetReady(shardID, event.ID)
		client.EventManager().DispatchEvent(&events.GuildReady{
			GenericGuild: genericGuildEvent,
		})
		if len(client.Caches().Guilds().UnreadyGuilds(shardID)) == 0 {
			client.EventManager().DispatchEvent(&events.GuildsReady{
				GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			})
		}
		if client.MemberChunkingManager().MemberChunkingFilter()(event.ID) {
			go func() {
				if _, err := client.MemberChunkingManager().RequestMembersWithQuery(event.ID, "", 0); err != nil {
					client.Logger().Error("failed to chunk guild on guild_create. error: ", err)
				}
			}()
		}

	}
	if wasUnavailable {
		client.Caches().Guilds().SetAvailable(event.ID)
		client.EventManager().DispatchEvent(&events.GuildAvailable{
			GenericGuild: genericGuildEvent,
		})
	} else {
		client.EventManager().DispatchEvent(&events.GuildJoin{
			GenericGuild: genericGuildEvent,
		})
	}
}
