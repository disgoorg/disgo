package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/snowflake"
)

type Caches interface {
	Config() CacheConfig

	GetMemberPermissions(member discord.Member) discord.Permissions
	GetMemberPermissionsInChannel(channel discord.GuildChannel, member discord.Member) discord.Permissions
	MemberRoles(member discord.Member) []discord.Role

	Roles() GroupedCache[discord.Role]
	Members() GroupedCache[discord.Member]
	ThreadMembers() GroupedCache[discord.ThreadMember]
	Presences() GroupedCache[discord.Presence]
	VoiceStates() GroupedCache[discord.VoiceState]
	Messages() GroupedCache[discord.Message]
	Emojis() GroupedCache[discord.Emoji]
	Stickers() GroupedCache[discord.Sticker]
	Guilds() GuildCache
	Channels() ChannelCache
	StageInstances() GroupedCache[discord.StageInstance]
	GuildScheduledEvents() GroupedCache[discord.GuildScheduledEvent]
}

func NewCaches(config CacheConfig) Caches {
	return &cachesImpl{
		config: config,

		guildCache:               NewGuildCache(config.CacheFlags, CacheFlagGuilds, nil),
		channelCache:             NewChannelCache(config.CacheFlags, CacheFlagsAllChannels, nil),
		stageInstanceCache:       NewGroupedCache[discord.StageInstance](config.CacheFlags, CacheFlagStageInstances, nil),
		guildScheduledEventCache: NewGroupedCache[discord.GuildScheduledEvent](config.CacheFlags, CacheFlagGuildScheduledEvents, nil),
		roleCache:                NewGroupedCache[discord.Role](config.CacheFlags, CacheFlagRoles, nil),
		memberCache:              NewGroupedCache[discord.Member](config.CacheFlags, CacheFlagMembers, config.MemberCachePolicy),
		threadMemberCache:        NewGroupedCache[discord.ThreadMember](config.CacheFlags, CacheFlagThreadMembers, nil),
		presenceCache:            NewGroupedCache[discord.Presence](config.CacheFlags, CacheFlagPresences, nil),
		voiceStateCache:          NewGroupedCache[discord.VoiceState](config.CacheFlags, CacheFlagVoiceStates, nil),
		messageCache:             NewGroupedCache[discord.Message](config.CacheFlags, CacheFlagMessages, config.MessageCachePolicy),
		emojiCache:               NewGroupedCache[discord.Emoji](config.CacheFlags, CacheFlagEmojis, nil),
		stickerCache:             NewGroupedCache[discord.Sticker](config.CacheFlags, CacheFlagStickers, nil),
	}
}

type cachesImpl struct {
	config CacheConfig

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

func (c *cachesImpl) Config() CacheConfig {
	return c.config
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
	return 0
	/*channel, ok := c.Channels().GetGuildChannel(channelID)
	if !ok {
		return discord.PermissionsNone
	}

	if guild, ok := c.Guilds().Get(channel.GuildID()); ok && guild.OwnerID == member.User.ID {
		return discord.PermissionsAll
	}

	permissions := c.GetMemberPermissions(member)
	if permissions.Has(discord.PermissionAdministrator) {
		return discord.PermissionsAll
	}

	var (
		allowRaw discord.Permissions
		denyRaw  discord.Permissions
	)
	if overwrite := channel.RolePermissionOverwrite(channel.GuildID()); overwrite != nil {
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

		overwrite := channel.RolePermissionOverwrite(roleID)
		if overwrite == nil {
			break
		}
		allowRole = allowRole.Add(overwrite.Allow)
		denyRole = denyRole.Add(overwrite.Deny)
	}

	allowRaw = (allowRaw & (denyRole - 1)) | allowRole
	denyRaw = (denyRaw & (allowRole - 1)) | denyRole

	if overwrite := channel.MemberPermissionOverwrite(member.User.ID); overwrite != nil {
		allowRaw = (allowRaw & (overwrite.Deny - 1)) | overwrite.Allow
		denyRaw = (denyRaw & (overwrite.Allow - 1)) | overwrite.Deny
	}

	permissions &= denyRaw - 1
	permissions |= allowRaw

	if member.CommunicationDisabledUntil != nil {
		permissions &= discord.PermissionViewChannel | discord.PermissionReadMessageHistory
	}
	return permissions*/
}

func (c *cachesImpl) MemberRoles(member discord.Member) []discord.Role {
	return c.Roles().FindAll(func(groupID snowflake.Snowflake, role discord.Role) bool {
		if groupID != member.GuildID {
			return false
		}
		for _, roleID := range member.RoleIDs {
			if roleID == role.ID {
				return true
			}
		}
		return false
	})
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
