package cache

import (
	"sync"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

// Caches combines all different entity caches into one with some utility methods.
type Caches interface {
	// CacheFlags returns the current configured FLags of the caches.
	CacheFlags() Flags

	// GetMemberPermissions returns the calculated permissions of the given member.
	// This requires the FlagRoles to be set.
	GetMemberPermissions(member discord.Member) discord.Permissions

	// GetMemberPermissionsInChannel returns the calculated permissions of the given member in the given channel.
	// This requires the FlagRoles and FlagChannels to be set.
	GetMemberPermissionsInChannel(channel discord.GuildChannel, member discord.Member) discord.Permissions

	// MemberRoles returns all roles of the given member.
	// This requires the FlagRoles to be set.
	MemberRoles(member discord.Member) []discord.Role

	// AudioChannelMembers returns all members which are in the given audio channel.
	// This requires the FlagVoiceStates to be set.
	AudioChannelMembers(channel discord.GuildAudioChannel) []discord.Member

	// GetSelfUser returns the current bot user.
	// This is only available after we received the gateway.EventTypeReady event.
	GetSelfUser() (discord.OAuth2User, bool)

	// PutSelfUser overrides the current bot user.
	PutSelfUser(user discord.OAuth2User)

	// GetSelfMember returns the current bot member from the given guildID.
	// This is only available after we received the gateway.EventTypeGuildCreate event for the given guildID.
	GetSelfMember(guildID snowflake.ID) (discord.Member, bool)

	// Roles returns the role cache.
	Roles() GroupedCache[discord.Role]

	// Members returns the member cache.
	Members() GroupedCache[discord.Member]

	// ThreadMembers returns the thread member cache.
	ThreadMembers() GroupedCache[discord.ThreadMember]

	// Presences returns the presence cache.
	Presences() GroupedCache[discord.Presence]

	// VoiceStates returns the voice state cache.
	VoiceStates() GroupedCache[discord.VoiceState]

	// Messages returns the message cache.
	Messages() GroupedCache[discord.Message]

	// Emojis returns the emoji cache.
	Emojis() GroupedCache[discord.Emoji]

	// Stickers returns the sticker cache.
	Stickers() GroupedCache[discord.Sticker]

	// Guilds returns the guild cache.
	Guilds() GuildCache

	// Channels returns the channel cache.
	Channels() ChannelCache

	// StageInstances returns the stage instance cache.
	StageInstances() GroupedCache[discord.StageInstance]

	// GuildScheduledEvents returns the guild scheduled event cache.
	GuildScheduledEvents() GroupedCache[discord.GuildScheduledEvent]
}

// New returns a new default Caches instance with the given ConfigOpt(s) applied.
func New(opts ...ConfigOpt) Caches {
	config := DefaultConfig()
	config.Apply(opts)

	return &cachesImpl{
		config: *config,

		guildCache:               NewGuildCache(config.CacheFlags, config.GuildCachePolicy),
		channelCache:             NewChannelCache(config.CacheFlags, config.ChannelCachePolicy),
		stageInstanceCache:       NewGroupedCache[discord.StageInstance](config.CacheFlags, FlagStageInstances, config.StageInstanceCachePolicy),
		guildScheduledEventCache: NewGroupedCache[discord.GuildScheduledEvent](config.CacheFlags, FlagGuildScheduledEvents, config.GuildScheduledEventCachePolicy),
		roleCache:                NewGroupedCache[discord.Role](config.CacheFlags, FlagRoles, config.RoleCachePolicy),
		memberCache:              NewGroupedCache[discord.Member](config.CacheFlags, FlagMembers, config.MemberCachePolicy),
		threadMemberCache:        NewGroupedCache[discord.ThreadMember](config.CacheFlags, FlagThreadMembers, config.ThreadMemberCachePolicy),
		presenceCache:            NewGroupedCache[discord.Presence](config.CacheFlags, FlagPresences, config.PresenceCachePolicy),
		voiceStateCache:          NewGroupedCache[discord.VoiceState](config.CacheFlags, FlagVoiceStates, config.VoiceStateCachePolicy),
		messageCache:             NewGroupedCache[discord.Message](config.CacheFlags, FlagMessages, config.MessageCachePolicy),
		emojiCache:               NewGroupedCache[discord.Emoji](config.CacheFlags, FlagEmojis, config.EmojiCachePolicy),
		stickerCache:             NewGroupedCache[discord.Sticker](config.CacheFlags, FlagStickers, config.StickerCachePolicy),
	}
}

type cachesImpl struct {
	config Config

	selfUser   *discord.OAuth2User
	selfUserMu sync.Mutex

	guildCache               GuildCache
	channelCache             ChannelCache
	stageInstanceCache       GroupedCache[discord.StageInstance]
	guildScheduledEventCache GroupedCache[discord.GuildScheduledEvent]
	roleCache                GroupedCache[discord.Role]
	memberCache              GroupedCache[discord.Member]
	threadMemberCache        GroupedCache[discord.ThreadMember]
	presenceCache            GroupedCache[discord.Presence]
	voiceStateCache          GroupedCache[discord.VoiceState]
	messageCache             GroupedCache[discord.Message]
	emojiCache               GroupedCache[discord.Emoji]
	stickerCache             GroupedCache[discord.Sticker]
}

func (c *cachesImpl) CacheFlags() Flags {
	return c.config.CacheFlags
}

func (c *cachesImpl) GetMemberPermissions(member discord.Member) discord.Permissions {
	if guild, ok := c.Guilds().Get(member.GuildID); ok && guild.OwnerID == member.User.ID {
		return discord.PermissionsAll
	}

	var permissions discord.Permissions
	if publicRole, ok := c.Roles().Get(member.GuildID, member.GuildID); ok {
		permissions = publicRole.Permissions
	}

	for _, role := range c.MemberRoles(member) {
		permissions = permissions.Add(role.Permissions)
		if permissions.Has(discord.PermissionAdministrator) {
			return discord.PermissionsAll
		}
	}
	if member.CommunicationDisabledUntil != nil {
		permissions &= discord.PermissionViewChannel | discord.PermissionReadMessageHistory
	}
	return permissions
}

func (c *cachesImpl) GetMemberPermissionsInChannel(channel discord.GuildChannel, member discord.Member) discord.Permissions {
	permissions := c.GetMemberPermissions(member)
	if permissions.Has(discord.PermissionAdministrator) {
		return discord.PermissionsAll
	}

	var (
		allowRaw discord.Permissions
		denyRaw  discord.Permissions
	)
	if overwrite, ok := channel.PermissionOverwrites().Role(channel.GuildID()); ok {
		allowRaw = overwrite.Allow
		denyRaw = overwrite.Deny
	}

	var (
		allowRole discord.Permissions
		denyRole  discord.Permissions
	)
	for _, roleID := range member.RoleIDs {
		if roleID == channel.GuildID() {
			continue
		}

		if overwrite, ok := channel.PermissionOverwrites().Role(roleID); ok {
			allowRole = allowRole.Add(overwrite.Allow)
			denyRole = denyRole.Add(overwrite.Deny)
		}
	}

	allowRaw = (allowRaw & (denyRole - 1)) | allowRole
	denyRaw = (denyRaw & (allowRole - 1)) | denyRole

	if overwrite, ok := channel.PermissionOverwrites().Member(member.User.ID); ok {
		allowRaw = (allowRaw & (overwrite.Deny - 1)) | overwrite.Allow
		denyRaw = (denyRaw & (overwrite.Allow - 1)) | overwrite.Deny
	}

	permissions &= denyRaw - 1
	permissions |= allowRaw

	if member.CommunicationDisabledUntil != nil {
		permissions &= discord.PermissionViewChannel | discord.PermissionReadMessageHistory
	}
	return permissions
}

func (c *cachesImpl) MemberRoles(member discord.Member) []discord.Role {
	return c.Roles().FindAll(func(groupID snowflake.ID, role discord.Role) bool {
		for _, roleID := range member.RoleIDs {
			if roleID == role.ID {
				return true
			}
		}
		return false
	})
}

func (c *cachesImpl) AudioChannelMembers(channel discord.GuildAudioChannel) []discord.Member {
	var members []discord.Member
	c.VoiceStates().GroupForEach(channel.GuildID(), func(state discord.VoiceState) {
		if member, ok := c.Members().Get(channel.GuildID(), state.UserID); ok {
			members = append(members, member)
		}
	})
	return members
}

func (c *cachesImpl) GetSelfUser() (discord.OAuth2User, bool) {
	c.selfUserMu.Lock()
	defer c.selfUserMu.Unlock()

	if c.selfUser == nil {
		return discord.OAuth2User{}, false
	}
	return *c.selfUser, true
}

func (c *cachesImpl) PutSelfUser(user discord.OAuth2User) {
	c.selfUserMu.Lock()
	defer c.selfUserMu.Unlock()

	c.selfUser = &user
}

func (c *cachesImpl) GetSelfMember(guildID snowflake.ID) (discord.Member, bool) {
	selfUser, ok := c.GetSelfUser()
	if !ok {
		return discord.Member{}, false
	}
	return c.Members().Get(guildID, selfUser.ID)
}

func (c *cachesImpl) Roles() GroupedCache[discord.Role] {
	return c.roleCache
}

func (c *cachesImpl) Members() GroupedCache[discord.Member] {
	return c.memberCache
}

func (c *cachesImpl) ThreadMembers() GroupedCache[discord.ThreadMember] {
	return c.threadMemberCache
}

func (c *cachesImpl) Presences() GroupedCache[discord.Presence] {
	return c.presenceCache
}

func (c *cachesImpl) VoiceStates() GroupedCache[discord.VoiceState] {
	return c.voiceStateCache
}

func (c *cachesImpl) Messages() GroupedCache[discord.Message] {
	return c.messageCache
}

func (c *cachesImpl) Emojis() GroupedCache[discord.Emoji] {
	return c.emojiCache
}

func (c *cachesImpl) Stickers() GroupedCache[discord.Sticker] {
	return c.stickerCache
}

func (c *cachesImpl) Guilds() GuildCache {
	return c.guildCache
}

func (c *cachesImpl) Channels() ChannelCache {
	return c.channelCache
}

func (c *cachesImpl) StageInstances() GroupedCache[discord.StageInstance] {
	return c.stageInstanceCache
}

func (c *cachesImpl) GuildScheduledEvents() GroupedCache[discord.GuildScheduledEvent] {
	return c.guildScheduledEventCache
}
