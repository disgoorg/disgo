package core

import (
	"strings"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/disgo/rest/route"
)

var _ Mentionable = (*Member)(nil)

type Member struct {
	discord.Member
	*User
	Bot *Bot
}

// Permissions returns the calculated Permissions the Member has in the Guild
func (m *Member) Permissions() discord.Permissions {
	return GetMemberPermissions(m)
}

// InteractionPermissions returns the Permissions the Member has in this Channel for this Interaction
func (m *Member) InteractionPermissions() discord.Permissions {
	if m.Member.Permissions != nil {
		return *m.Member.Permissions
	}
	return discord.PermissionsNone
}

func (m *Member) ChannelPermissions(channel GuildChannel) discord.Permissions {
	return GetMemberPermissionsInChannel(channel, m)
}

// Roles return all Role(s)the Member has
func (m *Member) Roles() []*Role {
	var roles []*Role
	allRoles := m.Bot.Caches.RoleCache().GuildCache(m.GuildID)
	for _, roleID := range m.RoleIDs {
		roles = append(roles, allRoles[roleID])
	}
	return roles
}

// VoiceState returns the VoiceState for this Member from the Caches(requires CacheFlagVoiceState and GatewayIntentsGuildVoiceStates)
func (m *Member) VoiceState() *VoiceState {
	return m.Bot.Caches.VoiceStateCache().Get(m.GuildID, m.User.ID)
}

// EffectiveName returns either the nickname or username depending on if the user has a nickname
func (m *Member) EffectiveName() string {
	if m.Nick != nil {
		return *m.Nick
	}
	return m.User.Username
}

// Guild returns the members guild from the caches
func (m *Member) Guild() *Guild {
	return m.Bot.Caches.GuildCache().Get(m.GuildID)
}

// IsOwner returns whether the member is the owner of the guild_events that it belongs to
func (m *Member) IsOwner() bool {
	if guild := m.Guild(); guild != nil {
		return guild.OwnerID == m.ID
	}
	return false
}

// AvatarURL returns the Avatar URL of the Member for this guild
func (m *Member) AvatarURL(size int) *string {
	if m.Avatar == nil {
		return nil
	}
	format := route.PNG
	if strings.HasPrefix(*m.Avatar, "a_") {
		format = route.GIF
	}
	compiledRoute, err := route.MemberAvatar.Compile(nil, format, size, m.GuildID, m.ID, *m.Avatar)
	if err != nil {
		return nil
	}
	url := compiledRoute.URL()
	return &url
}

func (m *Member) EffectiveAvatarURL(size int) string {
	if m.Avatar == nil {
		return m.User.EffectiveAvatarURL(size)
	}
	return *m.AvatarURL(size)
}

// Update updates the Member
func (m *Member) Update(updateGuildMember discord.MemberUpdate, opts ...rest.RequestOpt) (*Member, error) {
	member, err := m.Bot.RestServices.GuildService().UpdateMember(m.GuildID, m.User.ID, updateGuildMember, opts...)
	if err != nil {
		return nil, err
	}
	return m.Bot.EntityBuilder.CreateMember(m.GuildID, *member, CacheStrategyNoWs), nil
}

// Move moves/kicks the member to/from a voice channel
func (m *Member) Move(channelID discord.Snowflake, opts ...rest.RequestOpt) (*Member, error) {
	return m.Update(discord.MemberUpdate{ChannelID: &channelID}, opts...)
}

// Kick kicks the Member from the Guild
func (m *Member) Kick(opts ...rest.RequestOpt) error {
	return m.Bot.RestServices.GuildService().RemoveMember(m.GuildID, m.User.ID, opts...)
}

// Ban bans the Member from the Guild
func (m *Member) Ban(deleteMessageDays int, opts ...rest.RequestOpt) error {
	return m.Bot.RestServices.GuildService().AddBan(m.GuildID, m.User.ID, deleteMessageDays, opts...)
}

// Unban unbans the Member from the Guild
func (m *Member) Unban(opts ...rest.RequestOpt) error {
	return m.Bot.RestServices.GuildService().DeleteBan(m.GuildID, m.User.ID, opts...)
}

// AddRole adds a specific role the member
func (m *Member) AddRole(roleID discord.Snowflake, opts ...rest.RequestOpt) error {
	return m.Bot.RestServices.GuildService().AddMemberRole(m.GuildID, m.User.ID, roleID, opts...)
}

// RemoveRole removes a specific role the member
func (m *Member) RemoveRole(roleID discord.Snowflake, opts ...rest.RequestOpt) error {
	return m.Bot.RestServices.GuildService().RemoveMemberRole(m.GuildID, m.User.ID, roleID, opts...)
}
