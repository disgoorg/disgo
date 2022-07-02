package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerStageInstanceCreate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventStageInstanceCreate) {
	client.Caches().StageInstances().Put(event.GuildID, event.ID, event.StageInstance)

	client.EventManager().DispatchEvent(&events.StageInstanceCreate{
		GenericStageInstance: &events.GenericStageInstance{
			GenericEvent:    events.NewGenericEvent(client, sequenceNumber, shardID),
			StageInstanceID: event.ID,
			StageInstance:   event.StageInstance,
		},
	})
}

func gatewayHandlerStageInstanceUpdate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventStageInstanceUpdate) {
	oldStageInstance, _ := client.Caches().StageInstances().Get(event.GuildID, event.ID)
	client.Caches().StageInstances().Put(event.GuildID, event.ID, event.StageInstance)

	client.EventManager().DispatchEvent(&events.StageInstanceUpdate{
		GenericStageInstance: &events.GenericStageInstance{
			GenericEvent:    events.NewGenericEvent(client, sequenceNumber, shardID),
			StageInstanceID: event.ID,
			StageInstance:   event.StageInstance,
		},
		OldStageInstance: oldStageInstance,
	})
}

func gatewayHandlerStageInstanceDelete(client bot.Client, sequenceNumber int, shardID int, event gateway.EventStageInstanceDelete) {
	client.Caches().StageInstances().Remove(event.GuildID, event.ID)

	client.EventManager().DispatchEvent(&events.StageInstanceDelete{
		GenericStageInstance: &events.GenericStageInstance{
			GenericEvent:    events.NewGenericEvent(client, sequenceNumber, shardID),
			StageInstanceID: event.ID,
			StageInstance:   event.StageInstance,
		},
	})
}
