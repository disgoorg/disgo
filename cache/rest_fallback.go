package cache

import (
	"sync"

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

type sfCall struct {
	wg  sync.WaitGroup
	val any
	err error
}

type singleflight struct {
	mu sync.Mutex
	m  map[string]*sfCall
}

func (g *singleflight) Do(key string, fn func() (any, error)) (any, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*sfCall)
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}
	c := new(sfCall)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	defer func() {
		g.mu.Lock()
		delete(g.m, key)
		g.mu.Unlock()
		c.wg.Done()
	}()

	c.val, c.err = fn()
	return c.val, c.err
}

var _ Caches = (*restFallbackCachesImpl)(nil)

// NewRestFallbackCaches wraps a Caches implementation and provides automatic REST fallback
// for cache misses on enabled entity types. In-flight requests for the same entity are
// automatically deduplicated to prevent request storms.
//
// If the provided caches instance is already wrapped with REST fallback, the existing wrapper
// is updated with the merged fallback flags and rest client rather than double-wrapped.
func NewRestFallbackCaches(caches Caches, rest rest.Rest, fallbackFlags RestFallbackFlags) Caches {
	if fallbackFlags == RestFallbackFlagsNone || rest == nil {
		return caches
	}
	if existing, ok := caches.(*restFallbackCachesImpl); ok {
		existing.flags = existing.flags.Add(fallbackFlags)
		if rest != nil {
			existing.rest = rest
		}
		return existing
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
	sf    singleflight
}

func (c *restFallbackCachesImpl) Member(guildID snowflake.ID, userID snowflake.ID) (discord.Member, bool) {
	if member, ok := c.Caches.Member(guildID, userID); ok {
		return member, true
	}
	if !c.flags.Has(RestFallbackFlagMembers) {
		return discord.Member{}, false
	}
	key := "member:" + guildID.String() + ":" + userID.String()
	val, err := c.sf.Do(key, func() (any, error) {
		if member, ok := c.Caches.Member(guildID, userID); ok {
			return member, nil
		}
		member, err := c.rest.GetMember(guildID, userID)
		if err != nil {
			return discord.Member{}, err
		}
		c.AddMember(*member)
		return *member, nil
	})
	if err != nil {
		return discord.Member{}, false
	}
	return val.(discord.Member), true
}

func (c *restFallbackCachesImpl) SelfMember(guildID snowflake.ID) (discord.Member, bool) {
	if member, ok := c.Caches.SelfMember(guildID); ok {
		return member, true
	}
	if !c.flags.Has(RestFallbackFlagMembers) {
		return discord.Member{}, false
	}
	selfUser, ok := c.SelfUser()
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
	key := "guild:" + guildID.String()
	val, err := c.sf.Do(key, func() (any, error) {
		if guild, ok := c.Caches.Guild(guildID); ok {
			return guild, nil
		}
		restGuild, err := c.rest.GetGuild(guildID, true)
		if err != nil {
			return discord.CacheGuild{}, err
		}
		cacheGuild := discord.CacheGuild{
			Guild:       restGuild.Guild,
			MemberCount: restGuild.ApproximateMemberCount,
		}
		c.GuildCache().Put(guildID, cacheGuild)
		return cacheGuild, nil
	})
	if err != nil {
		return discord.CacheGuild{}, false
	}
	return val.(discord.CacheGuild), true
}

// Channel returns the discord.GuildChannel for the given channelID.
//
// Note: As defined by the Caches interface, Disgo's channel cache only stores guild channels.
// Non-guild channels (such as DM channels) returned by the REST fallback do not implement
// discord.GuildChannel and will return (nil, false).
func (c *restFallbackCachesImpl) Channel(channelID snowflake.ID) (discord.GuildChannel, bool) {
	if ch, ok := c.Caches.Channel(channelID); ok {
		return ch, true
	}
	if !c.flags.Has(RestFallbackFlagChannels) {
		return nil, false
	}
	key := "channel:" + channelID.String()
	val, err := c.sf.Do(key, func() (any, error) {
		if ch, ok := c.Caches.Channel(channelID); ok {
			return ch, nil
		}
		ch, err := c.rest.GetChannel(channelID)
		if err != nil {
			return nil, err
		}
		if guildCh, ok := ch.(discord.GuildChannel); ok {
			c.AddChannel(guildCh)
			return guildCh, nil
		}
		return nil, nil
	})
	if err != nil || val == nil {
		return nil, false
	}
	return val.(discord.GuildChannel), true
}

func (c *restFallbackCachesImpl) Role(guildID snowflake.ID, roleID snowflake.ID) (discord.Role, bool) {
	if role, ok := c.Caches.Role(guildID, roleID); ok {
		return role, true
	}
	if !c.flags.Has(RestFallbackFlagRoles) {
		return discord.Role{}, false
	}
	key := "roles:" + guildID.String()
	_, err := c.sf.Do(key, func() (any, error) {
		if _, ok := c.Caches.Role(guildID, roleID); ok {
			return nil, nil
		}
		roles, err := c.rest.GetRoles(guildID)
		if err != nil {
			return nil, err
		}
		for _, r := range roles {
			c.AddRole(r)
		}
		return nil, nil
	})
	if err != nil {
		return discord.Role{}, false
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
	key := "message:" + channelID.String() + ":" + messageID.String()
	val, err := c.sf.Do(key, func() (any, error) {
		if msg, ok := c.Caches.Message(channelID, messageID); ok {
			return msg, nil
		}
		msg, err := c.rest.GetMessage(channelID, messageID)
		if err != nil {
			return discord.Message{}, err
		}
		c.AddMessage(*msg)
		return *msg, nil
	})
	if err != nil {
		return discord.Message{}, false
	}
	return val.(discord.Message), true
}

func (c *restFallbackCachesImpl) Emoji(guildID snowflake.ID, emojiID snowflake.ID) (discord.Emoji, bool) {
	if emoji, ok := c.Caches.Emoji(guildID, emojiID); ok {
		return emoji, true
	}
	if !c.flags.Has(RestFallbackFlagEmojis) {
		return discord.Emoji{}, false
	}
	key := "emoji:" + guildID.String() + ":" + emojiID.String()
	val, err := c.sf.Do(key, func() (any, error) {
		if emoji, ok := c.Caches.Emoji(guildID, emojiID); ok {
			return emoji, nil
		}
		emoji, err := c.rest.GetEmoji(guildID, emojiID)
		if err != nil {
			return discord.Emoji{}, err
		}
		c.AddEmoji(*emoji)
		return *emoji, nil
	})
	if err != nil {
		return discord.Emoji{}, false
	}
	return val.(discord.Emoji), true
}

func (c *restFallbackCachesImpl) Sticker(guildID snowflake.ID, stickerID snowflake.ID) (discord.Sticker, bool) {
	if sticker, ok := c.Caches.Sticker(guildID, stickerID); ok {
		return sticker, true
	}
	if !c.flags.Has(RestFallbackFlagStickers) {
		return discord.Sticker{}, false
	}
	key := "sticker:" + stickerID.String()
	val, err := c.sf.Do(key, func() (any, error) {
		if sticker, ok := c.Caches.Sticker(guildID, stickerID); ok {
			return sticker, nil
		}
		sticker, err := c.rest.GetSticker(stickerID)
		if err != nil {
			return discord.Sticker{}, err
		}
		c.AddSticker(*sticker)
		return *sticker, nil
	})
	if err != nil {
		return discord.Sticker{}, false
	}
	return val.(discord.Sticker), true
}

func (c *restFallbackCachesImpl) StageInstance(guildID snowflake.ID, stageInstanceID snowflake.ID) (discord.StageInstance, bool) {
	if stage, ok := c.Caches.StageInstance(guildID, stageInstanceID); ok {
		return stage, true
	}
	if !c.flags.Has(RestFallbackFlagStageInstances) {
		return discord.StageInstance{}, false
	}
	key := "stage:" + stageInstanceID.String()
	val, err := c.sf.Do(key, func() (any, error) {
		if stage, ok := c.Caches.StageInstance(guildID, stageInstanceID); ok {
			return stage, nil
		}
		stage, err := c.rest.GetStageInstance(stageInstanceID)
		if err != nil {
			return discord.StageInstance{}, err
		}
		c.AddStageInstance(*stage)
		return *stage, nil
	})
	if err != nil {
		return discord.StageInstance{}, false
	}
	return val.(discord.StageInstance), true
}

func (c *restFallbackCachesImpl) GuildScheduledEvent(guildID snowflake.ID, guildScheduledEventID snowflake.ID) (discord.GuildScheduledEvent, bool) {
	if event, ok := c.Caches.GuildScheduledEvent(guildID, guildScheduledEventID); ok {
		return event, true
	}
	if !c.flags.Has(RestFallbackFlagGuildScheduledEvents) {
		return discord.GuildScheduledEvent{}, false
	}
	key := "event:" + guildID.String() + ":" + guildScheduledEventID.String()
	val, err := c.sf.Do(key, func() (any, error) {
		if event, ok := c.Caches.GuildScheduledEvent(guildID, guildScheduledEventID); ok {
			return event, nil
		}
		event, err := c.rest.GetGuildScheduledEvent(guildID, guildScheduledEventID, false)
		if err != nil {
			return discord.GuildScheduledEvent{}, err
		}
		c.AddGuildScheduledEvent(*event)
		return *event, nil
	})
	if err != nil {
		return discord.GuildScheduledEvent{}, false
	}
	return val.(discord.GuildScheduledEvent), true
}
