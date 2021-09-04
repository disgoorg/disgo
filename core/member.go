package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

type Member struct {
	discord.Member
	Disgo Disgo
	User  *User
}

// Permissions returns the Permissions the Member has in the Guild
func (m *Member) Permissions() discord.Permissions {
	if m.Member.Permissions != nil {
		return *m.Member.Permissions
	}
	return GetMemberPermissions(m)
}

// Roles return all Role(s)the Member has
func (m *Member) Roles() []*Role {
	var roles []*Role
	allRoles := m.Disgo.Caches().RoleCache().RoleCache(m.GuildID)
	for _, roleID := range m.RoleIDs {
		roles = append(roles, allRoles[roleID])
	}
	return roles
}

// VoiceState returns the VoiceState for this Member from the Caches(requires CacheFlagVoiceState and GatewayIntentsGuildVoiceStates)
func (m *Member) VoiceState() *VoiceState {
	return m.Disgo.Caches().VoiceStateCache().Get(m.GuildID, m.User.ID)
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
	return m.Disgo.Caches().GuildCache().Get(m.GuildID)
}

// IsOwner returns whether the member is the owner of the guild_events that it belongs to
func (m *Member) IsOwner() bool {
	if guild := m.Guild(); guild != nil {
		return guild.OwnerID == m.User.ID
	}
	return false
}

// Update updates the Member
func (m *Member) Update(updateGuildMember discord.MemberUpdate, opts ...rest.RequestOpt) (*Member, rest.Error) {
	member, err := m.Disgo.RestServices().GuildService().UpdateMember(m.GuildID, m.User.ID, updateGuildMember, opts...)
	if err != nil {
		return nil, err
	}
	return m.Disgo.EntityBuilder().CreateMember(m.GuildID, *member, CacheStrategyNoWs), nil
}

// Kick kicks the Member from the Guild
func (m *Member) Kick(opts ...rest.RequestOpt) rest.Error {
	return m.Disgo.RestServices().GuildService().RemoveMember(m.GuildID, m.User.ID, opts...)
}

// Ban bans the Member from the Guild
func (m *Member) Ban(deleteMessageDays int, opts ...rest.RequestOpt) rest.Error {
	return m.Disgo.RestServices().GuildService().AddBan(m.GuildID, m.User.ID, deleteMessageDays, opts...)
}

// Unban unbans the Member from the Guild
func (m *Member) Unban(opts ...rest.RequestOpt) rest.Error {
	return m.Disgo.RestServices().GuildService().DeleteBan(m.GuildID, m.User.ID, opts...)
}

// Move moves/kicks the member to/from a voice channel
func (m *Member) Move(channelID discord.Snowflake, opts ...rest.RequestOpt) (*Member, rest.Error) {
	member, err := m.Disgo.RestServices().GuildService().UpdateMember(m.GuildID, m.User.ID, discord.MemberUpdate{ChannelID: &channelID}, opts...)
	if err != nil {
		return nil, err
	}
	return m.Disgo.EntityBuilder().CreateMember(m.GuildID, *member, CacheStrategyNoWs), nil
}

// AddRole adds a specific role the member
func (m *Member) AddRole(roleID discord.Snowflake, opts ...rest.RequestOpt) rest.Error {
	return m.Disgo.RestServices().GuildService().AddMemberRole(m.GuildID, m.User.ID, roleID, opts...)
}

// RemoveRole removes a specific role the member
func (m *Member) RemoveRole(roleID discord.Snowflake, opts ...rest.RequestOpt) rest.Error {
	return m.Disgo.RestServices().GuildService().RemoveMemberRole(m.GuildID, m.User.ID, roleID, opts...)
}

// GetMemberPermissions returns all Permissions from the provided Member
func GetMemberPermissions(member *Member) discord.Permissions {
	if member.IsOwner() {
		return discord.PermissionsAll
	}
	if guild := member.Guild(); guild != nil {
		var permissions discord.Permissions
		for _, role := range member.Roles() {
			permissions = permissions.Add(role.Permissions)
			if permissions.Has(discord.PermissionAdministrator) {
				return discord.PermissionsAll
			}
		}
		return permissions
	}
	return discord.PermissionsNone
}
