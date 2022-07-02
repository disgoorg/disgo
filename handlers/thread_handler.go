package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerThreadCreate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventThreadCreate) {
	client.Caches().Channels().Put(event.ID(), event.GuildThread)
	client.Caches().ThreadMembers().Put(event.ID(), event.ThreadMember.UserID, event.ThreadMember)

	client.EventManager().DispatchEvent(&events.ThreadCreate{
		GenericThread: &events.GenericThread{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			ThreadID:     event.ID(),
			GuildID:      event.GuildID(),
			Thread:       event.GuildThread,
		},
		ThreadMember: event.ThreadMember,
	})
}

func gatewayHandlerThreadUpdate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventThreadUpdate) {
	oldGuildThread, _ := client.Caches().Channels().GetGuildThread(event.ID())
	client.Caches().Channels().Put(event.ID(), event.GuildThread)

	client.EventManager().DispatchEvent(&events.ThreadUpdate{
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

func gatewayHandlerThreadDelete(client bot.Client, sequenceNumber int, shardID int, event gateway.EventThreadDelete) {
	channel, _ := client.Caches().Channels().Remove(event.ID)
	client.Caches().ThreadMembers().RemoveAll(event.ID)

	client.EventManager().DispatchEvent(&events.ThreadDelete{
		GenericThread: &events.GenericThread{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			ThreadID:     event.ID,
			GuildID:      event.GuildID,
			ParentID:     event.ParentID,
			Thread:       channel.(discord.GuildThread),
		},
	})
}

func gatewayHandlerThreadListSync(client bot.Client, sequenceNumber int, shardID int, event gateway.EventThreadListSync) {
	for _, thread := range event.Threads {
		client.Caches().Channels().Put(thread.ID(), thread)
		client.EventManager().DispatchEvent(&events.ThreadShow{
			GenericThread: &events.GenericThread{
				GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
				Thread:       thread,
				ThreadID:     thread.ID(),
				GuildID:      event.GuildID,
			},
		})
	}
}

func gatewayHandlerThreadMemberUpdate(_ bot.Client, _ int, _ int, _ gateway.EventData) {
	// ThreadMembersUpdate kinda handles this already?
}

func gatewayHandlerThreadMembersUpdate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventThreadMembersUpdate) {
	genericEvent := events.NewGenericEvent(client, sequenceNumber, shardID)

	if thread, ok := client.Caches().Channels().GetGuildThread(event.ID); ok {
		thread.MemberCount = event.MemberCount
		client.Caches().Channels().Put(thread.ID(), thread)
	}

	for _, addedMember := range event.AddedMembers {
		addedMember.Member.GuildID = event.ID
		client.Caches().ThreadMembers().Put(event.ID, addedMember.UserID, addedMember.ThreadMember)
		client.Caches().Members().Put(event.GuildID, addedMember.UserID, addedMember.Member)

		if addedMember.Presence != nil {
			client.Caches().Presences().Put(event.GuildID, addedMember.UserID, *addedMember.Presence)
		}

		client.EventManager().DispatchEvent(&events.ThreadMemberAdd{
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
		threadMember, _ := client.Caches().ThreadMembers().Remove(event.ID, removedMemberID)

		client.EventManager().DispatchEvent(&events.ThreadMemberRemove{
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
