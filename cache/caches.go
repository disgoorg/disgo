package cache

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/snowflake"
)

type Caches interface {
	Config() Config

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

func NewCaches(config Config) Caches {
	return &cachesImpl{
		config: config,

		guildCache:               NewGuildCache(config.CacheFlags, FlagGuilds, nil),
		channelCache:             NewChannelCache(config.CacheFlags, FlagsAllChannels, nil),
		stageInstanceCache:       NewGroupedCache[discord.StageInstance](config.CacheFlags, FlagStageInstances, nil),
		guildScheduledEventCache: NewGroupedCache[discord.GuildScheduledEvent](config.CacheFlags, FlagGuildScheduledEvents, nil),
		roleCache:                NewGroupedCache[discord.Role](config.CacheFlags, FlagRoles, nil),
		memberCache:              NewGroupedCache[discord.Member](config.CacheFlags, FlagMembers, config.MemberCachePolicy),
		threadMemberCache:        NewGroupedCache[discord.ThreadMember](config.CacheFlags, FlagThreadMembers, nil),
		presenceCache:            NewGroupedCache[discord.Presence](config.CacheFlags, FlagPresences, nil),
		voiceStateCache:          NewGroupedCache[discord.VoiceState](config.CacheFlags, FlagVoiceStates, nil),
		messageCache:             NewGroupedCache[discord.Message](config.CacheFlags, FlagMessages, config.MessageCachePolicy),
		emojiCache:               NewGroupedCache[discord.Emoji](config.CacheFlags, FlagEmojis, nil),
		stickerCache:             NewGroupedCache[discord.Sticker](config.CacheFlags, FlagStickers, nil),
	}
}

type cachesImpl struct {
	config Config

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

func (c *cachesImpl) Config() Config {
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

	if guild, ok := c.Guilds().Option(channel.GuildID()); ok && guild.OwnerID == member.User.CommandID {
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

	if overwrite := channel.MemberPermissionOverwrite(member.User.CommandID); overwrite != nil {
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
