package gatewayhandlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
)

// gatewayHandlerGuildCreate handles core.GuildCreateGatewayEvent
type gatewayHandlerGuildCreate struct{}

// EventType returns the core.GatewayGatewayEventType
func (h *gatewayHandlerGuildCreate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildCreate) New() interface{} {
	return &discord.GatewayGuild{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildCreate) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*discord.GatewayGuild)

	wasUnready := bot.Caches.GuildCache().IsUnready(payload.ID)
	wasUnavailable := bot.Caches.GuildCache().IsUnavailable(payload.ID)

	guild := bot.EntityBuilder.CreateGuild(payload.Guild, core.CacheStrategyYes)

	for _, channel := range payload.Channels {
		channel.GuildID = &payload.ID
		bot.EntityBuilder.CreateChannel(channel, core.CacheStrategyYes)
	}

	for _, role := range payload.Roles {
		role.GuildID = payload.ID
		bot.EntityBuilder.CreateRole(payload.ID, role, core.CacheStrategyYes)
	}

	for _, member := range payload.Members {
		bot.EntityBuilder.CreateMember(payload.ID, member, core.CacheStrategyYes)
	}

	for _, voiceState := range payload.VoiceStates {
		voiceState.GuildID = payload.ID // populate unset field
		vs := bot.EntityBuilder.CreateVoiceState(voiceState, core.CacheStrategyYes)
		if channel := vs.Channel(); channel != nil {
			channel.ConnectedMemberIDs[voiceState.UserID] = struct{}{}
		}
	}

	for _, emote := range payload.Emojis {
		bot.EntityBuilder.CreateEmoji(payload.ID, emote, core.CacheStrategyYes)
	}

	for _, stageInstance := range payload.StageInstances {
		bot.EntityBuilder.CreateStageInstance(stageInstance, core.CacheStrategyYes)
	}

	for _, presence := range payload.Presences {
		bot.EntityBuilder.CreatePresence(presence, core.CacheStrategyYes)
	}

	genericGuildEvent := &events.GenericGuildEvent{
		GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
		GuildID:      payload.ID,
		Guild:        guild,
	}

	if wasUnready {
		bot.Caches.GuildCache().SetReady(payload.ID)
		bot.EventManager.Dispatch(&events.GuildReadyEvent{
			GenericGuildEvent: genericGuildEvent,
		})
		if len(bot.Caches.GuildCache().UnreadyGuilds()) == 0 {
			bot.EventManager.Dispatch(&events.GuildsReadyEvent{
				GenericEvent: events.NewGenericEvent(bot, -1),
			})
		}
		if bot.MemberChunkingManager.MemberChunkingFilter()(payload.ID) {
			go func() {
				if _, err := bot.MemberChunkingManager.LoadAllMembers(payload.ID); err != nil {
					bot.Logger.Error("failed to chunk guild on guild_create. error: ", err)
				}
			}()
		}

	} else if wasUnavailable {
		bot.Caches.GuildCache().SetAvailable(payload.ID)
		bot.EventManager.Dispatch(&events.GuildAvailableEvent{
			GenericGuildEvent: genericGuildEvent,
		})
	} else {
		bot.EventManager.Dispatch(&events.GuildJoinEvent{
			GenericGuildEvent: genericGuildEvent,
		})
	}
}
