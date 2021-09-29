package core

import (
	"strings"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/disgo/rest/route"
)

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

// ChannelPermissions returns the Permissions the Member has in the provided Channel
func (m *Member) ChannelPermissions(channel *Channel) discord.Permissions {
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

// VoiceState returns the VoiceState for this Member.
// This will only check cached voice states! (requires core.CacheFlagVoiceStates and discord.GatewayIntentGuildVoiceStates)
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

// Guild returns the Guild this Member is tied to.
// This will only check cached guilds!
func (m *Member) Guild() *Guild {
	return m.Bot.Caches.GuildCache().Get(m.GuildID)
}

// IsOwner returns whether this Member is the owner of the Guild
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

// EffectiveAvatarURL returns either the server avatar or global avatar depending on if the user has one
func (m *Member) EffectiveAvatarURL(size int) string {
	if m.Avatar == nil {
		return m.User.EffectiveAvatarURL(size)
	}
	return *m.AvatarURL(size)
}

// Update updates the Member with the properties provided in discord.MemberUpdate
func (m *Member) Update(updateGuildMember discord.MemberUpdate, opts ...rest.RequestOpt) (*Member, rest.Error) {
	member, err := m.Bot.RestServices.GuildService().UpdateMember(m.GuildID, m.User.ID, updateGuildMember, opts...)
	if err != nil {
		return nil, err
	}
	return m.Bot.EntityBuilder.CreateMember(m.GuildID, *member, CacheStrategyNoWs), nil
}

// Move moves/kicks the member to/from a voice channel
func (m *Member) Move(channelID discord.Snowflake, opts ...rest.RequestOpt) (*Member, rest.Error) {
	return m.Update(discord.MemberUpdate{ChannelID: &channelID}, opts...)
}

// Kick kicks this Member from the Guild
func (m *Member) Kick(opts ...rest.RequestOpt) rest.Error {
	return m.Bot.RestServices.GuildService().RemoveMember(m.GuildID, m.User.ID, opts...)
}

// Ban bans this Member from the Guild
func (m *Member) Ban(deleteMessageDays int, opts ...rest.RequestOpt) rest.Error {
	return m.Bot.RestServices.GuildService().AddBan(m.GuildID, m.User.ID, deleteMessageDays, opts...)
}

// Unban unbans this Member from the Guild
func (m *Member) Unban(opts ...rest.RequestOpt) rest.Error {
	return m.Bot.RestServices.GuildService().DeleteBan(m.GuildID, m.User.ID, opts...)
}

// AddRole adds a specific role this Member
func (m *Member) AddRole(roleID discord.Snowflake, opts ...rest.RequestOpt) rest.Error {
	return m.Bot.RestServices.GuildService().AddMemberRole(m.GuildID, m.User.ID, roleID, opts...)
}

// RemoveRole removes a specific role this Member
func (m *Member) RemoveRole(roleID discord.Snowflake, opts ...rest.RequestOpt) rest.Error {
	return m.Bot.RestServices.GuildService().RemoveMemberRole(m.GuildID, m.User.ID, roleID, opts...)
}
