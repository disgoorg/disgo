package handlers

import (
	"log/slog"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerVoiceChannelEffectSend(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventVoiceChannelEffectSend) {
	client.EventManager.DispatchEvent(&events.GuildVoiceChannelEffectSend{
		GenericEvent:                events.NewGenericEvent(client, sequenceNumber, shardID),
		EventVoiceChannelEffectSend: event,
	})
}

func gatewayHandlerVoiceStateUpdate(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventVoiceStateUpdate) {
	member := event.Member

	oldVoiceState, err := client.Caches.VoiceState(event.GuildID, event.UserID)
	oldOk := err == nil
	if event.ChannelID == nil {
		if _, err := client.Caches.RemoveVoiceState(event.GuildID, event.UserID); err != nil {
			client.Logger.Error("failed to remove voice state from cache", slog.Any("err", err), slog.String("guild_id", event.GuildID.String()), slog.String("user_id", event.UserID.String()))
		}
	} else {
		if err := client.Caches.AddVoiceState(event.VoiceState); err != nil {
			client.Logger.Error("failed to add voice state to cache", slog.Any("err", err), slog.String("guild_id", event.GuildID.String()), slog.String("user_id", event.UserID.String()))
		}
	}
	if err := client.Caches.AddMember(member); err != nil {
		client.Logger.Error("failed to add member to cache", slog.Any("err", err), slog.String("guild_id", event.GuildID.String()), slog.String("user_id", event.UserID.String()))
	}

	if event.UserID == client.ID() && client.VoiceManager != nil {
		client.VoiceManager.HandleVoiceStateUpdate(event)
	}

	genericGuildVoiceEvent := &events.GenericGuildVoiceState{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		VoiceState:   event.VoiceState,
		Member:       member,
	}

	client.EventManager.DispatchEvent(&events.GuildVoiceStateUpdate{
		GenericGuildVoiceState: genericGuildVoiceEvent,
		OldVoiceState:          oldVoiceState,
	})

	if oldOk && oldVoiceState.ChannelID != nil && event.ChannelID != nil {
		client.EventManager.DispatchEvent(&events.GuildVoiceMove{
			GenericGuildVoiceState: genericGuildVoiceEvent,
			OldVoiceState:          oldVoiceState,
		})
	} else if (oldOk || oldVoiceState.ChannelID == nil) && event.ChannelID != nil {
		client.EventManager.DispatchEvent(&events.GuildVoiceJoin{
			GenericGuildVoiceState: genericGuildVoiceEvent,
		})
	} else if event.ChannelID == nil {
		client.EventManager.DispatchEvent(&events.GuildVoiceLeave{
			GenericGuildVoiceState: genericGuildVoiceEvent,
			OldVoiceState:          oldVoiceState,
		})
	} else {
		client.Logger.Warn("could not decide which GuildVoice to fire")
	}
}

func gatewayHandlerVoiceServerUpdate(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventVoiceServerUpdate) {
	if client.VoiceManager != nil {
		client.VoiceManager.HandleVoiceServerUpdate(event)
	}

	client.EventManager.DispatchEvent(&events.VoiceServerUpdate{
		GenericEvent:           events.NewGenericEvent(client, sequenceNumber, shardID),
		EventVoiceServerUpdate: event,
	})
}
