package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerGuildCreate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildCreate) {
	wasUnready := client.Caches().IsGuildUnready(event.ID)
	wasUnavailable := client.Caches().IsGuildUnavailable(event.ID)

	client.Caches().AddGuild(event.Guild)

	for _, channel := range event.Channels {
		channel = discord.ApplyGuildIDToChannel(channel, event.ID) // populate unset field
		client.Caches().AddChannel(channel)
	}

	for _, thread := range event.Threads {
		thread = discord.ApplyGuildIDToThread(thread, event.ID) // populate unset field
		client.Caches().AddChannel(thread)
	}

	for _, role := range event.Roles {
		role.GuildID = event.ID // populate unset field
		client.Caches().AddRole(role)
	}

	for _, member := range event.Members {
		member.GuildID = event.ID // populate unset field
		client.Caches().AddMember(member)
	}

	for _, voiceState := range event.VoiceStates {
		voiceState.GuildID = event.ID // populate unset field
		client.Caches().AddVoiceState(voiceState)
	}

	for _, emoji := range event.Emojis {
		emoji.GuildID = event.ID // populate unset field
		client.Caches().AddEmoji(emoji)
	}

	for _, sticker := range event.Stickers {
		sticker.GuildID = &event.ID // populate unset field
		client.Caches().AddSticker(sticker)
	}

	for _, stageInstance := range event.StageInstances {
		client.Caches().AddStageInstance(stageInstance)
	}

	for _, guildScheduledEvent := range event.GuildScheduledEvents {
		client.Caches().AddGuildScheduledEvent(guildScheduledEvent)
	}

	for _, presence := range event.Presences {
		presence.GuildID = event.ID // populate unset field
		client.Caches().AddPresence(presence)
	}

	genericGuildEvent := &events.GenericGuild{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		GuildID:      event.ID,
		Guild:        event.Guild,
	}

	if wasUnready {
		client.Caches().SetGuildUnready(event.ID, false)
		client.EventManager().DispatchEvent(&events.GuildReady{
			GenericGuild: genericGuildEvent,
		})
		if len(client.Caches().UnreadyGuildIDs()) == 0 {
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
		client.Caches().SetGuildUnavailable(event.ID, false)
		client.EventManager().DispatchEvent(&events.GuildAvailable{
			GenericGuild: genericGuildEvent,
		})
	} else {
		client.EventManager().DispatchEvent(&events.GuildJoin{
			GenericGuild: genericGuildEvent,
		})
	}
}
