package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
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

	shard, _ := bot.Shard(payload.ID)
	shardID := shard.ShardID()

	wasUnready := bot.Caches.Guilds().IsUnready(shardID, payload.ID)
	wasUnavailable := bot.Caches.Guilds().IsUnavailable(payload.ID)

	guild := bot.EntityBuilder.CreateGuild(payload.Guild, core.CacheStrategyYes)

	for _, channel := range payload.Channels {
		bot.EntityBuilder.CreateChannel(setGuildID(channel, payload.ID), core.CacheStrategyYes)
	}

	for _, thread := range payload.Threads {
		bot.EntityBuilder.CreateChannel(setGuildID(thread, payload.ID), core.CacheStrategyYes)
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
			if ch, ok := channel.(*core.GuildVoiceChannel); ok {
				ch.ConnectedMemberIDs[voiceState.UserID] = struct{}{}
			} else if ch, ok := channel.(*core.GuildStageVoiceChannel); ok {
				ch.ConnectedMemberIDs[voiceState.UserID] = struct{}{}
			}
		}
	}

	for _, emote := range payload.Emojis {
		bot.EntityBuilder.CreateEmoji(payload.ID, emote, core.CacheStrategyYes)
	}

	for _, sticker := range payload.Stickers {
		bot.EntityBuilder.CreateSticker(sticker, core.CacheStrategyYes)
	}

	for _, stageInstance := range payload.StageInstances {
		bot.EntityBuilder.CreateStageInstance(stageInstance, core.CacheStrategyYes)
	}

	for _, guildScheduledEvent := range payload.GuildScheduledEvents {
		bot.EntityBuilder.CreateGuildScheduledEvent(guildScheduledEvent, core.CacheStrategyYes)
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
		bot.Caches.Guilds().SetReady(shardID, payload.ID)
		bot.EventManager.Dispatch(&events.GuildReadyEvent{
			GenericGuildEvent: genericGuildEvent,
		})
		if len(bot.Caches.Guilds().UnreadyGuilds(shardID)) == 0 {
			bot.EventManager.Dispatch(&events.GuildsReadyEvent{
				GenericEvent: events.NewGenericEvent(bot, -1),
				ShardID:      shardID,
			})
		}
		if bot.MemberChunkingManager.MemberChunkingFilter()(payload.ID) {
			go func() {
				if _, err := bot.MemberChunkingManager.RequestMembersWithQuery(payload.ID, "", 0); err != nil {
					bot.Logger.Error("failed to chunk guild on guild_create. error: ", err)
				}
			}()
		}

	} else if wasUnavailable {
		bot.Caches.Guilds().SetAvailable(payload.ID)
		bot.EventManager.Dispatch(&events.GuildAvailableEvent{
			GenericGuildEvent: genericGuildEvent,
		})
	} else {
		bot.EventManager.Dispatch(&events.GuildJoinEvent{
			GenericGuildEvent: genericGuildEvent,
		})
	}
}

func setGuildID(channel discord.GuildChannel, guildID discord.Snowflake) discord.GuildChannel {
	switch ch := channel.(type) {
	case discord.GuildTextChannel:
		ch.ChannelGuildID = guildID
		return ch

	case discord.GuildVoiceChannel:
		ch.ChannelGuildID = guildID
		return ch

	case discord.GuildCategoryChannel:
		ch.ChannelGuildID = guildID
		return ch

	case discord.GuildNewsChannel:
		ch.ChannelGuildID = guildID
		return ch

	case discord.GuildStoreChannel:
		ch.ChannelGuildID = guildID
		return ch

	case discord.GuildNewsThread:
		ch.ChannelGuildID = guildID
		return ch

	case discord.GuildPrivateThread:
		ch.ChannelGuildID = guildID
		return ch

	case discord.GuildPublicThread:
		ch.ChannelGuildID = guildID
		return ch

	case discord.GuildStageVoiceChannel:
		ch.ChannelGuildID = guildID
		return ch
	}
	return nil
}
