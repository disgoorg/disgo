package cache

import (
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/internal/flags"
	"github.com/disgoorg/disgo/rest"
)

// RestFallbackFlags is used to configure which entity types will fallback to REST requests on cache misses.
type RestFallbackFlags int

const (
	RestFallbackFlagGuilds RestFallbackFlags = 1 << iota
	RestFallbackFlagChannels
	RestFallbackFlagRoles
	RestFallbackFlagMembers
	RestFallbackFlagMessages
	RestFallbackFlagEmojis
	RestFallbackFlagStickers
	RestFallbackFlagStageInstances
	RestFallbackFlagGuildScheduledEvents

	RestFallbackFlagsAll = RestFallbackFlagGuilds |
		RestFallbackFlagChannels |
		RestFallbackFlagRoles |
		RestFallbackFlagMembers |
		RestFallbackFlagMessages |
		RestFallbackFlagEmojis |
		RestFallbackFlagStickers |
		RestFallbackFlagStageInstances |
		RestFallbackFlagGuildScheduledEvents

	RestFallbackFlagsNone RestFallbackFlags = 0
)

func (f RestFallbackFlags) Add(bits ...RestFallbackFlags) RestFallbackFlags {
	return flags.Add(f, bits...)
}

func (f RestFallbackFlags) Remove(bits ...RestFallbackFlags) RestFallbackFlags {
	return flags.Remove(f, bits...)
}

func (f RestFallbackFlags) Has(bits ...RestFallbackFlags) bool {
	return flags.Has(f, bits...)
}

func (f RestFallbackFlags) Missing(bits ...RestFallbackFlags) bool {
	return flags.Missing(f, bits...)
}

var _ Caches = (*restFallbackCachesImpl)(nil)

// NewRestFallbackCaches wraps a Caches implementation and provides automatic REST fallback
// for cache misses on enabled entity types.
func NewRestFallbackCaches(caches Caches, rest rest.Rest, fallbackFlags RestFallbackFlags) Caches {
	if fallbackFlags == RestFallbackFlagsNone || rest == nil {
		return caches
	}
	return &restFallbackCachesImpl{
		Caches: caches,
		rest:   rest,
		flags:  fallbackFlags,
	}
}

type restFallbackCachesImpl struct {
	Caches
	rest  rest.Rest
	flags RestFallbackFlags
}

func (c *restFallbackCachesImpl) Member(guildID snowflake.ID, userID snowflake.ID) (discord.Member, bool) {
	if member, ok := c.Caches.Member(guildID, userID); ok {
		return member, true
	}
	if !c.flags.Has(RestFallbackFlagMembers) {
		return discord.Member{}, false
	}
	member, err := c.rest.GetMember(guildID, userID)
	if err != nil {
		return discord.Member{}, false
	}
	c.Caches.AddMember(*member)
	return *member, true
}

func (c *restFallbackCachesImpl) SelfMember(guildID snowflake.ID) (discord.Member, bool) {
	if member, ok := c.Caches.SelfMember(guildID); ok {
		return member, true
	}
	if !c.flags.Has(RestFallbackFlagMembers) {
		return discord.Member{}, false
	}
	selfUser, ok := c.Caches.SelfUser()
	if !ok {
		return discord.Member{}, false
	}
	return c.Member(guildID, selfUser.ID)
}

func (c *restFallbackCachesImpl) Guild(guildID snowflake.ID) (discord.CacheGuild, bool) {
	if guild, ok := c.Caches.Guild(guildID); ok {
		return guild, true
	}
	if !c.flags.Has(RestFallbackFlagGuilds) {
		return discord.CacheGuild{}, false
	}
	restGuild, err := c.rest.GetGuild(guildID, true)
	if err != nil {
		return discord.CacheGuild{}, false
	}
	cacheGuild := discord.CacheGuild{
		Guild:       restGuild.Guild,
		MemberCount: restGuild.ApproximateMemberCount,
	}
	c.Caches.GuildCache().Put(guildID, cacheGuild)
	return cacheGuild, true
}

func (c *restFallbackCachesImpl) Channel(channelID snowflake.ID) (discord.GuildChannel, bool) {
	if ch, ok := c.Caches.Channel(channelID); ok {
		return ch, true
	}
	if !c.flags.Has(RestFallbackFlagChannels) {
		return nil, false
	}
	ch, err := c.rest.GetChannel(channelID)
	if err != nil {
		return nil, false
	}
	if guildCh, ok := ch.(discord.GuildChannel); ok {
		c.Caches.AddChannel(guildCh)
		return guildCh, true
	}
	return nil, false
}

func (c *restFallbackCachesImpl) Role(guildID snowflake.ID, roleID snowflake.ID) (discord.Role, bool) {
	if role, ok := c.Caches.Role(guildID, roleID); ok {
		return role, true
	}
	if !c.flags.Has(RestFallbackFlagRoles) {
		return discord.Role{}, false
	}
	roles, err := c.rest.GetRoles(guildID)
	if err != nil {
		return discord.Role{}, false
	}
	for _, r := range roles {
		c.Caches.AddRole(r)
	}
	return c.Caches.Role(guildID, roleID)
}

func (c *restFallbackCachesImpl) Message(channelID snowflake.ID, messageID snowflake.ID) (discord.Message, bool) {
	if msg, ok := c.Caches.Message(channelID, messageID); ok {
		return msg, true
	}
	if !c.flags.Has(RestFallbackFlagMessages) {
		return discord.Message{}, false
	}
	msg, err := c.rest.GetMessage(channelID, messageID)
	if err != nil {
		return discord.Message{}, false
	}
	c.Caches.AddMessage(*msg)
	return *msg, true
}

func (c *restFallbackCachesImpl) Emoji(guildID snowflake.ID, emojiID snowflake.ID) (discord.Emoji, bool) {
	if emoji, ok := c.Caches.Emoji(guildID, emojiID); ok {
		return emoji, true
	}
	if !c.flags.Has(RestFallbackFlagEmojis) {
		return discord.Emoji{}, false
	}
	emoji, err := c.rest.GetEmoji(guildID, emojiID)
	if err != nil {
		return discord.Emoji{}, false
	}
	c.Caches.AddEmoji(*emoji)
	return *emoji, true
}

func (c *restFallbackCachesImpl) Sticker(guildID snowflake.ID, stickerID snowflake.ID) (discord.Sticker, bool) {
	if sticker, ok := c.Caches.Sticker(guildID, stickerID); ok {
		return sticker, true
	}
	if !c.flags.Has(RestFallbackFlagStickers) {
		return discord.Sticker{}, false
	}
	sticker, err := c.rest.GetSticker(stickerID)
	if err != nil {
		return discord.Sticker{}, false
	}
	c.Caches.AddSticker(*sticker)
	return *sticker, true
}

func (c *restFallbackCachesImpl) StageInstance(guildID snowflake.ID, stageInstanceID snowflake.ID) (discord.StageInstance, bool) {
	if stage, ok := c.Caches.StageInstance(guildID, stageInstanceID); ok {
		return stage, true
	}
	if !c.flags.Has(RestFallbackFlagStageInstances) {
		return discord.StageInstance{}, false
	}
	stage, err := c.rest.GetStageInstance(stageInstanceID)
	if err != nil {
		return discord.StageInstance{}, false
	}
	c.Caches.AddStageInstance(*stage)
	return *stage, true
}

func (c *restFallbackCachesImpl) GuildScheduledEvent(guildID snowflake.ID, guildScheduledEventID snowflake.ID) (discord.GuildScheduledEvent, bool) {
	if event, ok := c.Caches.GuildScheduledEvent(guildID, guildScheduledEventID); ok {
		return event, true
	}
	if !c.flags.Has(RestFallbackFlagGuildScheduledEvents) {
		return discord.GuildScheduledEvent{}, false
	}
	event, err := c.rest.GetGuildScheduledEvent(guildID, guildScheduledEventID, false)
	if err != nil {
		return discord.GuildScheduledEvent{}, false
	}
	c.Caches.AddGuildScheduledEvent(*event)
	return *event, true
}
