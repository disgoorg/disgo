package cache

import (
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
)

func defaultConfig() config {
	return config{
		GuildCachePolicy:                PolicyAll[discord.Guild],
		ChannelCachePolicy:              PolicyAll[discord.GuildChannel],
		StageInstanceCachePolicy:        PolicyAll[discord.StageInstance],
		GuildScheduledEventCachePolicy:  PolicyAll[discord.GuildScheduledEvent],
		GuildSoundboardSoundCachePolicy: PolicyAll[discord.SoundboardSound],
		RoleCachePolicy:                 PolicyAll[discord.Role],
		MemberCachePolicy:               PolicyAll[discord.Member],
		ThreadMemberCachePolicy:         PolicyAll[discord.ThreadMember],
		PresenceCachePolicy:             PolicyAll[discord.Presence],
		VoiceStateCachePolicy:           PolicyAll[discord.VoiceState],
		MessageCachePolicy:              PolicyAll[discord.Message],
		EmojiCachePolicy:                PolicyAll[discord.Emoji],
		StickerCachePolicy:              PolicyAll[discord.Sticker],
	}
}

type config struct {
	CacheFlags Flags

	SelfUserCache SelfUserCache

	GuildCache       GuildCache
	GuildCachePolicy Policy[discord.Guild]

	ChannelCache       ChannelCache
	ChannelCachePolicy Policy[discord.GuildChannel]

	StageInstanceCache       StageInstanceCache
	StageInstanceCachePolicy Policy[discord.StageInstance]

	GuildScheduledEventCache       GuildScheduledEventCache
	GuildScheduledEventCachePolicy Policy[discord.GuildScheduledEvent]

	GuildSoundboardSoundCache       GuildSoundboardSoundCache
	GuildSoundboardSoundCachePolicy Policy[discord.SoundboardSound]

	RoleCache       RoleCache
	RoleCachePolicy Policy[discord.Role]

	MemberCache       MemberCache
	MemberCachePolicy Policy[discord.Member]

	ThreadMemberCache       ThreadMemberCache
	ThreadMemberCachePolicy Policy[discord.ThreadMember]

	PresenceCache       PresenceCache
	PresenceCachePolicy Policy[discord.Presence]

	VoiceStateCache       VoiceStateCache
	VoiceStateCachePolicy Policy[discord.VoiceState]

	MessageCache       MessageCache
	MessageCachePolicy Policy[discord.Message]

	EmojiCache       EmojiCache
	EmojiCachePolicy Policy[discord.Emoji]

	StickerCache       StickerCache
	StickerCachePolicy Policy[discord.Sticker]
}

// ConfigOpt is a type alias for a function that takes a config and is used to configure your Caches.
type ConfigOpt func(config *config)

func (c *config) apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
	if c.SelfUserCache == nil {
		c.SelfUserCache = NewSelfUserCache()
	}
	if c.GuildCache == nil {
		c.GuildCache = NewGuildCache(NewCache[discord.Guild](c.CacheFlags, FlagGuilds, c.GuildCachePolicy), NewSet[snowflake.ID](), NewSet[snowflake.ID]())
	}
	if c.ChannelCache == nil {
		c.ChannelCache = NewChannelCache(NewCache[discord.GuildChannel](c.CacheFlags, FlagChannels, c.ChannelCachePolicy))
	}
	if c.StageInstanceCache == nil {
		c.StageInstanceCache = NewStageInstanceCache(NewGroupedCache[discord.StageInstance](c.CacheFlags, FlagStageInstances, c.StageInstanceCachePolicy))
	}
	if c.GuildScheduledEventCache == nil {
		c.GuildScheduledEventCache = NewGuildScheduledEventCache(NewGroupedCache[discord.GuildScheduledEvent](c.CacheFlags, FlagGuildScheduledEvents, c.GuildScheduledEventCachePolicy))
	}
	if c.GuildSoundboardSoundCache == nil {
		c.GuildSoundboardSoundCache = NewGuildSoundboardSoundCache(NewGroupedCache[discord.SoundboardSound](c.CacheFlags, FlagGuildSoundboardSounds, c.GuildSoundboardSoundCachePolicy))
	}
	if c.RoleCache == nil {
		c.RoleCache = NewRoleCache(NewGroupedCache[discord.Role](c.CacheFlags, FlagRoles, c.RoleCachePolicy))
	}
	if c.MemberCache == nil {
		c.MemberCache = NewMemberCache(NewGroupedCache[discord.Member](c.CacheFlags, FlagMembers, c.MemberCachePolicy))
	}
	if c.ThreadMemberCache == nil {
		c.ThreadMemberCache = NewThreadMemberCache(NewGroupedCache[discord.ThreadMember](c.CacheFlags, FlagThreadMembers, c.ThreadMemberCachePolicy))
	}
	if c.PresenceCache == nil {
		c.PresenceCache = NewPresenceCache(NewGroupedCache[discord.Presence](c.CacheFlags, FlagPresences, c.PresenceCachePolicy))
	}
	if c.VoiceStateCache == nil {
		c.VoiceStateCache = NewVoiceStateCache(NewGroupedCache[discord.VoiceState](c.CacheFlags, FlagVoiceStates, c.VoiceStateCachePolicy))
	}
	if c.MessageCache == nil {
		c.MessageCache = NewMessageCache(NewGroupedCache[discord.Message](c.CacheFlags, FlagMessages, c.MessageCachePolicy))
	}
	if c.EmojiCache == nil {
		c.EmojiCache = NewEmojiCache(NewGroupedCache[discord.Emoji](c.CacheFlags, FlagEmojis, c.EmojiCachePolicy))
	}
	if c.StickerCache == nil {
		c.StickerCache = NewStickerCache(NewGroupedCache[discord.Sticker](c.CacheFlags, FlagStickers, c.StickerCachePolicy))
	}
}

// WithCaches sets the Flags of the config.
func WithCaches(flags ...Flags) ConfigOpt {
	return func(config *config) {
		config.CacheFlags = config.CacheFlags.Add(flags...)
	}
}

// WithGuildCachePolicy sets the Policy[discord.Guild] of the config.
func WithGuildCachePolicy(policy Policy[discord.Guild]) ConfigOpt {
	return func(config *config) {
		config.GuildCachePolicy = policy
	}
}

// WithGuildCache sets the GuildCache of the config.
func WithGuildCache(guildCache GuildCache) ConfigOpt {
	return func(config *config) {
		config.GuildCache = guildCache
	}
}

// WithChannelCachePolicy sets the Policy[discord.Channel] of the config.
func WithChannelCachePolicy(policy Policy[discord.GuildChannel]) ConfigOpt {
	return func(config *config) {
		config.ChannelCachePolicy = policy
	}
}

// WithChannelCache sets the ChannelCache of the config.
func WithChannelCache(channelCache ChannelCache) ConfigOpt {
	return func(config *config) {
		config.ChannelCache = channelCache
	}
}

// WithStageInstanceCachePolicy sets the Policy[discord.Guild] of the config.
func WithStageInstanceCachePolicy(policy Policy[discord.StageInstance]) ConfigOpt {
	return func(config *config) {
		config.StageInstanceCachePolicy = policy
	}
}

// WithStageInstanceCache sets the StageInstanceCache of the config.
func WithStageInstanceCache(stageInstanceCache StageInstanceCache) ConfigOpt {
	return func(config *config) {
		config.StageInstanceCache = stageInstanceCache
	}
}

// WithGuildScheduledEventCachePolicy sets the Policy[discord.GuildScheduledEvent] of the config.
func WithGuildScheduledEventCachePolicy(policy Policy[discord.GuildScheduledEvent]) ConfigOpt {
	return func(config *config) {
		config.GuildScheduledEventCachePolicy = policy
	}
}

// WithGuildScheduledEventCache sets the GuildScheduledEventCache of the config.
func WithGuildScheduledEventCache(guildScheduledEventCache GuildScheduledEventCache) ConfigOpt {
	return func(config *config) {
		config.GuildScheduledEventCache = guildScheduledEventCache
	}
}

// WithGuildSoundboardSoundCache sets the GuildSoundboardSoundCache of the config.
func WithGuildSoundboardSoundCache(guildSoundboardSoundCache GuildSoundboardSoundCache) ConfigOpt {
	return func(config *config) {
		config.GuildSoundboardSoundCache = guildSoundboardSoundCache
	}
}

// WithRoleCachePolicy sets the Policy[discord.Role] of the config.
func WithRoleCachePolicy(policy Policy[discord.Role]) ConfigOpt {
	return func(config *config) {
		config.RoleCachePolicy = policy
	}
}

// WithRoleCache sets the RoleCache of the config.
func WithRoleCache(roleCache RoleCache) ConfigOpt {
	return func(config *config) {
		config.RoleCache = roleCache
	}
}

// WithMemberCachePolicy sets the Policy[discord.Member] of the config.
func WithMemberCachePolicy(policy Policy[discord.Member]) ConfigOpt {
	return func(config *config) {
		config.MemberCachePolicy = policy
	}
}

// WithMemberCache sets the MemberCache of the config.
func WithMemberCache(memberCache MemberCache) ConfigOpt {
	return func(config *config) {
		config.MemberCache = memberCache
	}
}

// WithThreadMemberCachePolicy sets the Policy[discord.ThreadMember] of the config.
func WithThreadMemberCachePolicy(policy Policy[discord.ThreadMember]) ConfigOpt {
	return func(config *config) {
		config.ThreadMemberCachePolicy = policy
	}
}

// WithThreadMemberCache sets the ThreadMemberCache of the config.
func WithThreadMemberCache(threadMemberCache ThreadMemberCache) ConfigOpt {
	return func(config *config) {
		config.ThreadMemberCache = threadMemberCache
	}
}

// WithPresenceCachePolicy sets the Policy[discord.Presence] of the config.
func WithPresenceCachePolicy(policy Policy[discord.Presence]) ConfigOpt {
	return func(config *config) {
		config.PresenceCachePolicy = policy
	}
}

// WithPresenceCache sets the PresenceCache of the config.
func WithPresenceCache(presenceCache PresenceCache) ConfigOpt {
	return func(config *config) {
		config.PresenceCache = presenceCache
	}
}

// WithVoiceStateCachePolicy sets the Policy[discord.VoiceState] of the config.
func WithVoiceStateCachePolicy(policy Policy[discord.VoiceState]) ConfigOpt {
	return func(config *config) {
		config.VoiceStateCachePolicy = policy
	}
}

// WithVoiceStateCache sets the VoiceStateCache of the config.
func WithVoiceStateCache(voiceStateCache VoiceStateCache) ConfigOpt {
	return func(config *config) {
		config.VoiceStateCache = voiceStateCache
	}
}

// WithMessageCachePolicy sets the Policy[discord.Message] of the config.
func WithMessageCachePolicy(policy Policy[discord.Message]) ConfigOpt {
	return func(config *config) {
		config.MessageCachePolicy = policy
	}
}

// WithMessageCache sets the MessageCache of the config.
func WithMessageCache(messageCache MessageCache) ConfigOpt {
	return func(config *config) {
		config.MessageCache = messageCache
	}
}

// WithEmojiCachePolicy sets the Policy[discord.Emoji] of the config.
func WithEmojiCachePolicy(policy Policy[discord.Emoji]) ConfigOpt {
	return func(config *config) {
		config.EmojiCachePolicy = policy
	}
}

// WithEmojiCache sets the EmojiCache of the config.
func WithEmojiCache(emojiCache EmojiCache) ConfigOpt {
	return func(config *config) {
		config.EmojiCache = emojiCache
	}
}

// WithStickerCachePolicy sets the Policy[discord.Sticker] of the config.
func WithStickerCachePolicy(policy Policy[discord.Sticker]) ConfigOpt {
	return func(config *config) {
		config.StickerCachePolicy = policy
	}
}

// WithStickerCache sets the StickerCache of the config.
func WithStickerCache(stickerCache StickerCache) ConfigOpt {
	return func(config *config) {
		config.StickerCache = stickerCache
	}
}
