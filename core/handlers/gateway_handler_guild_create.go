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
func (h *gatewayHandlerGuildCreate) HandleGatewayEvent(bot core.Bot, sequenceNumber discord.GatewaySequence, v interface{}) {
	payload := *v.(*discord.GatewayGuild)

	shard, _ := bot.Shard(payload.ID)
	shardID := shard.ShardID()

	wasUnready := bot.Caches().Guilds().IsUnready(shardID, payload.ID)
	wasUnavailable := bot.Caches().Guilds().IsUnavailable(payload.ID)

	for _, channel := range payload.Channels {

	}

	for _, thread := range payload.Threads {

	}

	for _, role := range payload.Roles {
		role.GuildID = payload.ID

	}

	for _, member := range payload.Members {

	}

	for _, voiceState := range payload.VoiceStates {
		voiceState.GuildID = payload.ID // populate unset field

	}

	for _, emoji := range payload.Emojis {
		emoji.GuildID = payload.ID

	}

	for _, sticker := range payload.Stickers {

	}

	for _, stageInstance := range payload.StageInstances {

	}

	for _, guildScheduledEvent := range payload.GuildScheduledEvents {

	}

	for _, presence := range payload.Presences {

	}

	genericGuildEvent := &events.GenericGuildEvent{
		GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
		GuildID:      payload.ID,
		Guild:        guild,
	}

	if wasUnready {
		bot.Caches().Guilds().SetReady(shardID, payload.ID)
		bot.EventManager().Dispatch(&events.GuildReadyEvent{
			GenericGuildEvent: genericGuildEvent,
		})
		if len(bot.Caches().Guilds().UnreadyGuilds(shardID)) == 0 {
			bot.EventManager().Dispatch(&events.GuildsReadyEvent{
				GenericEvent: events.NewGenericEvent(bot, -1),
				ShardID:      shardID,
			})
		}
		if bot.MemberChunkingManager().MemberChunkingFilter()(payload.ID) {
			go func() {
				if _, err := bot.MemberChunkingManager().RequestMembersWithQuery(payload.ID, "", 0); err != nil {
					bot.Logger().Error("failed to chunk guild on guild_create. error: ", err)
				}
			}()
		}

	} else if wasUnavailable {
		bot.Caches().Guilds().SetAvailable(payload.ID)
		bot.EventManager().Dispatch(&events.GuildAvailableEvent{
			GenericGuildEvent: genericGuildEvent,
		})
	} else {
		bot.EventManager().Dispatch(&events.GuildJoinEvent{
			GenericGuildEvent: genericGuildEvent,
		})
	}
}
