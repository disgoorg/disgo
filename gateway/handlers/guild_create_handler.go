package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
)

// gatewayHandlerGuildCreate handles core.GuildCreateGatewayEvent
type gatewayHandlerGuildCreate struct{}

// EventType returns the core.GatewayGatewayEventType
func (h *gatewayHandlerGuildCreate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildCreate) New() any {
	return &discord.GatewayGuild{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildCreate) HandleGatewayEvent(client bot.Client, sequenceNumber discord.GatewaySequence, v any) {
	/*payload := *v.(*discord.GatewayGuild)

	shard, _ := bot.Shard(payload.ID)
	shardID := shard.ShardID()

	wasUnready := client.Caches().Guilds().IsUnready(shardID, payload.ID)
	wasUnavailable := client.Caches().Guilds().IsUnavailable(payload.ID)

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
		GenericEvent: events.NewGenericEvent(client, sequenceNumber),
		GuildID:      payload.ID,
		Guild:        guild,
	}

	if wasUnready {
		client.Caches().Guilds().SetReady(shardID, payload.ID)
		client.EventManager().Dispatch(&events.GuildReadyEvent{
			GenericGuildEvent: genericGuildEvent,
		})
		if len(client.Caches().Guilds().UnreadyGuilds(shardID)) == 0 {
			client.EventManager().Dispatch(&events.GuildsReadyEvent{
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
		client.Caches().Guilds().SetAvailable(payload.ID)
		client.EventManager().Dispatch(&events.GuildAvailableEvent{
			GenericGuildEvent: genericGuildEvent,
		})
	} else {
		client.EventManager().Dispatch(&events.GuildJoinEvent{
			GenericGuildEvent: genericGuildEvent,
		})
	}*/
}
