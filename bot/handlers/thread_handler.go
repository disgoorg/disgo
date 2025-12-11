package handlers

import (
	"log/slog"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerThreadCreate(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventThreadCreate) {
	client.Caches.AddChannel(event.GuildThread)
	client.Caches.AddThreadMember(event.ThreadMember)

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
	if err != nil && err != cache.ErrNotFound {
		client.Logger.Error("failed to get thread from cache", slog.Any("err", err), slog.String("thread_id", event.ID().String()))
	}
	client.Caches.AddChannel(event.GuildThread)

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
	if err != nil && err != cache.ErrNotFound {
		client.Logger.Error("failed to remove channel from cache", slog.Any("err", err), slog.String("channel_id", event.ID.String()))
	}
	if err == nil {
		thread, _ = channel.(discord.GuildThread)
	}
	client.Caches.RemoveThreadMembersByThreadID(event.ID)

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
		client.Caches.AddChannel(thread)
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
	if err != nil && err != cache.ErrNotFound {
		client.Logger.Error("failed to get thread from cache", slog.Any("err", err), slog.String("thread_id", event.ID.String()))
	}
	if err == nil {
		thread.MemberCount = event.MemberCount
		client.Caches.AddChannel(thread)
	}

	for _, addedMember := range event.AddedMembers {
		addedMember.Member.GuildID = event.ID
		client.Caches.AddThreadMember(addedMember.ThreadMember)
		client.Caches.AddMember(addedMember.Member)

		if addedMember.Presence != nil {
			client.Caches.AddPresence(*addedMember.Presence)
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
		if err != nil && err != cache.ErrNotFound {
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
