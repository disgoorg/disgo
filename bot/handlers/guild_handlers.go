package handlers

import (
	"context"
	"errors"
	"log/slog"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerGuildCreate(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildCreate) {
	wasUnready, err := client.Caches.IsGuildUnready(event.ID)
	if err != nil {
		client.Logger.Error("failed to check if guild is unready", slog.Any("err", err), slog.String("guild_id", event.ID.String()))
	}
	wasUnavailable, err := client.Caches.IsGuildUnavailable(event.ID)
	if err != nil {
		client.Logger.Error("failed to check if guild is unavailable", slog.Any("err", err), slog.String("guild_id", event.ID.String()))
	}

	if err := client.Caches.AddGuild(event.Guild); err != nil {
		client.Logger.Error("failed to add guild to cache", slog.Any("err", err), slog.String("guild_id", event.ID.String()))
	}

	for _, channel := range event.Channels {
		channel = discord.ApplyGuildIDToChannel(channel, event.ID) // populate unset field
		if err := client.Caches.AddChannel(channel); err != nil {
			client.Logger.Error("failed to add channel to cache", slog.Any("err", err), slog.String("channel_id", channel.ID().String()))
		}
	}

	for _, thread := range event.Threads {
		thread = discord.ApplyGuildIDToThread(thread, event.ID) // populate unset field
		if err := client.Caches.AddChannel(thread); err != nil {
			client.Logger.Error("failed to add thread to cache", slog.Any("err", err), slog.String("thread_id", thread.ID().String()))
		}
	}

	for _, role := range event.Roles {
		role.GuildID = event.ID // populate unset field
		if err := client.Caches.AddRole(role); err != nil {
			client.Logger.Error("failed to add role to cache", slog.Any("err", err), slog.String("role_id", role.ID.String()))
		}
	}

	for _, member := range event.Members {
		member.GuildID = event.ID // populate unset field
		if err := client.Caches.AddMember(member); err != nil {
			client.Logger.Error("failed to add member to cache", slog.Any("err", err), slog.String("member_id", member.User.ID.String()))
		}
	}

	for _, voiceState := range event.VoiceStates {
		voiceState.GuildID = event.ID // populate unset field
		if err := client.Caches.AddVoiceState(voiceState); err != nil {
			client.Logger.Error("failed to add voice state to cache", slog.Any("err", err), slog.String("voice_state_id", voiceState.UserID.String()))
		}
	}

	for _, emoji := range event.Emojis {
		emoji.GuildID = event.ID // populate unset field
		if err := client.Caches.AddEmoji(emoji); err != nil {
			client.Logger.Error("failed to add emoji to cache", slog.Any("err", err), slog.String("emoji_id", emoji.ID.String()))
		}
	}

	for _, sticker := range event.Stickers {
		sticker.GuildID = &event.ID // populate unset field
		if err := client.Caches.AddSticker(sticker); err != nil {
			client.Logger.Error("failed to add sticker to cache", slog.Any("err", err), slog.String("sticker_id", sticker.ID.String()))
		}
	}

	for _, stageInstance := range event.StageInstances {
		if err := client.Caches.AddStageInstance(stageInstance); err != nil {
			client.Logger.Error("failed to add stage instance to cache", slog.Any("err", err), slog.String("stage_instance_id", stageInstance.ID.String()))
		}
	}

	for _, guildScheduledEvent := range event.GuildScheduledEvents {
		if err := client.Caches.AddGuildScheduledEvent(guildScheduledEvent); err != nil {
			client.Logger.Error("failed to add guild scheduled event to cache", slog.Any("err", err), slog.String("guild_scheduled_event_id", guildScheduledEvent.ID.String()))
		}
	}

	for _, soundboardSound := range event.SoundboardSounds {
		if err := client.Caches.AddGuildSoundboardSound(soundboardSound); err != nil {
			client.Logger.Error("failed to add guild soundboard sound to cache", slog.Any("err", err), slog.String("soundboard_sound_id", soundboardSound.SoundID.String()))
		}
	}

	for _, presence := range event.Presences {
		presence.GuildID = event.ID // populate unset field
		if err := client.Caches.AddPresence(presence); err != nil {
			client.Logger.Error("failed to add presence to cache", slog.Any("err", err), slog.String("presence_id", presence.PresenceUser.ID.String()))
		}
	}

	genericGuildEvent := &events.GenericGuild{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		GuildID:      event.ID,
	}

	if wasUnready {
		if err := client.Caches.SetGuildUnready(event.ID, false); err != nil {
			client.Logger.Error("failed to set guild unready", slog.Any("err", err), slog.String("guild_id", event.ID.String()))
		}

		client.EventManager.DispatchEvent(&events.GuildReady{
			GenericGuild: genericGuildEvent,
			Guild:        event.GatewayGuild,
		})
		unreadyGuildIDs, err := client.Caches.UnreadyGuildIDs()
		if err != nil {
			client.Logger.Error("failed to get unready guild ids", slog.Any("err", err))
		}

		if len(unreadyGuildIDs) == 0 {
			client.EventManager.DispatchEvent(&events.GuildsReady{
				GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			})
		}
		if client.MemberChunkingManager.MemberChunkingFilter()(event.ID) {
			go func() {
				if _, err := client.MemberChunkingManager.RequestMembersWithQuery(context.Background(), event.ID, "", 0); err != nil {
					client.Logger.Error("failed to chunk guild on guild_create", slog.Any("err", err))
				}
			}()
		}

		return
	}
	if wasUnavailable {
		if err := client.Caches.SetGuildUnavailable(event.ID, false); err != nil {
			client.Logger.Error("failed to set guild unavailable", slog.Any("err", err), slog.String("guild_id", event.ID.String()))
		}
		client.EventManager.DispatchEvent(&events.GuildAvailable{
			GenericGuild: genericGuildEvent,
			Guild:        event.GatewayGuild,
		})
	} else {
		client.EventManager.DispatchEvent(&events.GuildJoin{
			GenericGuild: genericGuildEvent,
			Guild:        event.GatewayGuild,
		})
	}
}

func gatewayHandlerGuildUpdate(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildUpdate) {
	oldGuild, err := client.Caches.Guild(event.ID)
	if err != nil && !errors.Is(err, cache.ErrNotFound) {
		client.Logger.Error("failed to get guild from cache", slog.Any("err", err), slog.String("guild_id", event.ID.String()))
	}
	if err := client.Caches.AddGuild(event.Guild); err != nil {
		client.Logger.Error("failed to add guild to cache", slog.Any("err", err), slog.String("guild_id", event.ID.String()))
	}

	client.EventManager.DispatchEvent(&events.GuildUpdate{
		GenericGuild: &events.GenericGuild{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			GuildID:      event.ID,
		},
		Guild:    event.Guild,
		OldGuild: oldGuild,
	})
}

func gatewayHandlerGuildDelete(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildDelete) {
	if event.Unavailable {
		if err := client.Caches.SetGuildUnavailable(event.ID, true); err != nil {
			client.Logger.Error("failed to set guild unavailable", slog.Any("err", err), slog.String("guild_id", event.ID.String()))
		}
	}

	guild, err := client.Caches.RemoveGuild(event.ID)
	if err != nil && !errors.Is(err, cache.ErrNotFound) {
		client.Logger.Error("failed to remove guild from cache", slog.Any("err", err), slog.String("guild_id", event.ID.String()))
	}
	if err := client.Caches.RemoveVoiceStatesByGuildID(event.ID); err != nil {
		client.Logger.Error("failed to remove voice states from cache", slog.Any("err", err), slog.String("guild_id", event.ID.String()))
	}
	if err := client.Caches.RemovePresencesByGuildID(event.ID); err != nil {
		client.Logger.Error("failed to remove presences from cache", slog.Any("err", err), slog.String("guild_id", event.ID.String()))
	}
	// TODO: figure out a better way to remove thread members from cache via guild id without requiring cached GuildThreads
	channels, err := client.Caches.Channels()
	if err == nil {
		for channel, err := range channels {
			if err != nil {
				client.Logger.Error("failed to get channel from cache", slog.Any("err", err), slog.String("channel_id", channel.ID().String()))
				break
			}

			if guildThread, ok := channel.(discord.GuildThread); ok && guildThread.GuildID() == event.ID {
				if err := client.Caches.RemoveThreadMembersByThreadID(guildThread.ID()); err != nil {
					client.Logger.Error("failed to remove thread members from cache", slog.Any("err", err), slog.String("thread_id", guildThread.ID().String()))
				}
			}
		}
	}
	if err := client.Caches.RemoveChannelsByGuildID(event.ID); err != nil {
		client.Logger.Error("failed to remove channels from cache", slog.Any("err", err), slog.String("guild_id", event.ID.String()))
	}
	if err := client.Caches.RemoveEmojisByGuildID(event.ID); err != nil {
		client.Logger.Error("failed to remove emojis from cache", slog.Any("err", err), slog.String("guild_id", event.ID.String()))
	}
	if err := client.Caches.RemoveStickersByGuildID(event.ID); err != nil {
		client.Logger.Error("failed to remove stickers from cache", slog.Any("err", err), slog.String("guild_id", event.ID.String()))
	}
	if err := client.Caches.RemoveRolesByGuildID(event.ID); err != nil {
		client.Logger.Error("failed to remove roles from cache", slog.Any("err", err), slog.String("guild_id", event.ID.String()))
	}
	if err := client.Caches.RemoveMembersByGuildID(event.ID); err != nil {
		client.Logger.Error("failed to remove members from cache", slog.Any("err", err), slog.String("guild_id", event.ID.String()))
	}
	if err := client.Caches.RemoveStageInstancesByGuildID(event.ID); err != nil {
		client.Logger.Error("failed to remove stage instances from cache", slog.Any("err", err), slog.String("guild_id", event.ID.String()))
	}
	if err := client.Caches.RemoveGuildScheduledEventsByGuildID(event.ID); err != nil {
		client.Logger.Error("failed to remove guild scheduled events from cache", slog.Any("err", err), slog.String("guild_id", event.ID.String()))
	}
	if err := client.Caches.RemoveGuildSoundboardSoundsByGuildID(event.ID); err != nil {
		client.Logger.Error("failed to remove guild soundboard sounds from cache", slog.Any("err", err), slog.String("guild_id", event.ID.String()))
	}
	if err := client.Caches.RemoveMessagesByGuildID(event.ID); err != nil {
		client.Logger.Error("failed to remove messages from cache", slog.Any("err", err), slog.String("guild_id", event.ID.String()))
	}

	genericGuildEvent := &events.GenericGuild{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		GuildID:      event.ID,
	}

	if event.Unavailable {
		client.EventManager.DispatchEvent(&events.GuildUnavailable{
			GenericGuild: genericGuildEvent,
			Guild:        guild,
		})
	} else {
		client.EventManager.DispatchEvent(&events.GuildLeave{
			GenericGuild: genericGuildEvent,
			Guild:        guild,
		})
	}
}

func gatewayHandlerGuildAuditLogEntryCreate(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildAuditLogEntryCreate) {
	genericGuildEvent := &events.GenericGuild{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		GuildID:      event.GuildID,
	}

	client.EventManager.DispatchEvent(&events.GuildAuditLogEntryCreate{
		GenericGuild:  genericGuildEvent,
		AuditLogEntry: event.AuditLogEntry,
	})
}
