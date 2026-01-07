package handlers

import (
	"errors"
	"log/slog"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerThreadCreate(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventThreadCreate) {
	if err := client.Caches.AddChannel(event.GuildThread); err != nil {
		client.Logger.Error("failed to add channel to cache", slog.Any("err", err), slog.String("thread_id", event.ID().String()))
	}
	if err := client.Caches.AddThreadMember(event.ThreadMember); err != nil {
		client.Logger.Error("failed to add thread member to cache", slog.Any("err", err), slog.String("thread_id", event.ID().String()), slog.String("user_id", event.ThreadMember.UserID.String()))
	}

	client.EventManager.DispatchEvent(&events.ThreadCreate{
		GenericThread: &events.GenericThread{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			ThreadID:     event.ID(),
			GuildID:      event.GuildID(),
			Thread:       event.GuildThread,
		},
		ThreadMember: event.ThreadMember,
		NewlyCreated: event.NewlyCreated,
	})
}

func gatewayHandlerThreadUpdate(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventThreadUpdate) {
	oldGuildThread, err := client.Caches.GuildThread(event.ID())
	if err != nil && !errors.Is(err, cache.ErrNotFound) {
		client.Logger.Error("failed to get thread from cache", slog.Any("err", err), slog.String("thread_id", event.ID().String()))
	}
	if err := client.Caches.AddChannel(event.GuildThread); err != nil {
		client.Logger.Error("failed to add channel to cache", slog.Any("err", err), slog.String("thread_id", event.ID().String()))
	}

	client.EventManager.DispatchEvent(&events.ThreadUpdate{
		GenericThread: &events.GenericThread{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			Thread:       event.GuildThread,
			ThreadID:     event.ID(),
			GuildID:      event.GuildID(),
			ParentID:     *event.ParentID(),
		},
		OldThread: oldGuildThread,
	})
}

func gatewayHandlerThreadDelete(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventThreadDelete) {
	var thread discord.GuildThread
	channel, err := client.Caches.RemoveChannel(event.ID)
	if err != nil && !errors.Is(err, cache.ErrNotFound) {
		client.Logger.Error("failed to remove channel from cache", slog.Any("err", err), slog.String("channel_id", event.ID.String()))
	}
	if err == nil {
		thread, _ = channel.(discord.GuildThread)
	}
	if err := client.Caches.RemoveThreadMembersByThreadID(event.ID); err != nil {
		client.Logger.Error("failed to remove thread members from cache", slog.Any("err", err), slog.String("thread_id", event.ID.String()))
	}

	client.EventManager.DispatchEvent(&events.ThreadDelete{
		GenericThread: &events.GenericThread{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			ThreadID:     event.ID,
			GuildID:      event.GuildID,
			ParentID:     event.ParentID,
			Thread:       thread,
		},
	})
}

func gatewayHandlerThreadListSync(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventThreadListSync) {
	for _, thread := range event.Threads {
		if err := client.Caches.AddChannel(thread); err != nil {
			client.Logger.Error("failed to add channel to cache", slog.Any("err", err), slog.String("thread_id", thread.ID().String()))
		}
		client.EventManager.DispatchEvent(&events.ThreadShow{
			GenericThread: &events.GenericThread{
				GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
				Thread:       thread,
				ThreadID:     thread.ID(),
				GuildID:      event.GuildID,
			},
		})
	}
}

func gatewayHandlerThreadMemberUpdate(_ *bot.Client, _ int, _ int, _ gateway.EventData) {
	// ThreadMembersUpdate kinda handles this already?
}

func gatewayHandlerThreadMembersUpdate(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventThreadMembersUpdate) {
	genericEvent := events.NewGenericEvent(client, sequenceNumber, shardID)

	thread, err := client.Caches.GuildThread(event.ID)
	if err != nil && !errors.Is(err, cache.ErrNotFound) {
		client.Logger.Error("failed to get thread from cache", slog.Any("err", err), slog.String("thread_id", event.ID.String()))
	}
	if err == nil {
		thread.MemberCount = event.MemberCount
		if err := client.Caches.AddChannel(thread); err != nil {
			client.Logger.Error("failed to add channel to cache", slog.Any("err", err), slog.String("thread_id", event.ID.String()))
		}
	}

	for _, addedMember := range event.AddedMembers {
		addedMember.Member.GuildID = event.ID
		if err := client.Caches.AddThreadMember(addedMember.ThreadMember); err != nil {
			client.Logger.Error("failed to add thread member to cache", slog.Any("err", err), slog.String("thread_id", event.ID.String()), slog.String("user_id", addedMember.UserID.String()))
		}
		if err := client.Caches.AddMember(addedMember.Member); err != nil {
			client.Logger.Error("failed to add member to cache", slog.Any("err", err), slog.String("guild_id", event.GuildID.String()), slog.String("user_id", addedMember.UserID.String()))
		}

		if addedMember.Presence != nil {
			if err := client.Caches.AddPresence(*addedMember.Presence); err != nil {
				client.Logger.Error("failed to add presence to cache", slog.Any("err", err), slog.String("guild_id", event.GuildID.String()), slog.String("user_id", addedMember.UserID.String()))
			}
		}

		client.EventManager.DispatchEvent(&events.ThreadMemberAdd{
			GenericThreadMember: &events.GenericThreadMember{
				GenericEvent:   genericEvent,
				GuildID:        event.GuildID,
				ThreadID:       event.ID,
				ThreadMemberID: addedMember.UserID,
				ThreadMember:   addedMember.ThreadMember,
			},
			Member:   addedMember.Member,
			Presence: addedMember.Presence,
		})
	}

	for _, removedMemberID := range event.RemovedMemberIDs {
		threadMember, err := client.Caches.RemoveThreadMember(event.ID, removedMemberID)
		if err != nil && !errors.Is(err, cache.ErrNotFound) {
			client.Logger.Error("failed to remove thread member from cache", slog.Any("err", err), slog.String("thread_id", event.ID.String()), slog.String("user_id", removedMemberID.String()))
		}

		client.EventManager.DispatchEvent(&events.ThreadMemberRemove{
			GenericThreadMember: &events.GenericThreadMember{
				GenericEvent:   genericEvent,
				GuildID:        event.GuildID,
				ThreadID:       event.ID,
				ThreadMemberID: removedMemberID,
				ThreadMember:   threadMember,
			},
		})
	}
}
