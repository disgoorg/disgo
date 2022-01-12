package core

import (
	"strings"
	"time"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/json"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/disgo/rest/route"
)

type Member struct {
	discord.Member
	User User
	Bot  Bot
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
func (m *Member) ChannelPermissions(channel GuildChannel) discord.Permissions {
	return GetMemberPermissionsInChannel(channel, m)
}

// Roles return all Role(s)the Member has
func (m *Member) Roles() []Role {
	var roles []Role
	allRoles := m.Bot.Caches().Roles().GroupCache(m.GuildID)
	for _, roleID := range m.RoleIDs {
		roles = append(roles, allRoles[roleID])
	}
	return roles
}

// VoiceState returns the VoiceState for this Member.
// This will only check cached voice states! (requires core.CacheFlagVoiceStates and discord.GatewayIntentGuildVoiceStates)
func (m *Member) VoiceState() (VoiceState, bool) {
	return m.Bot.Caches().VoiceStates().Get(m.GuildID, m.User.ID)
}

// Guild returns the Guild this Member is tied to.
// This will only check cached guilds!
func (m *Member) Guild() (Guild, bool) {
	return m.Bot.Caches().Guilds().Get(m.GuildID)
}

// IsOwner returns whether this Member is the owner of the Guild
func (m *Member) IsOwner() bool {
	if guild, ok := m.Guild(); ok {
		return guild.OwnerID == m.User.ID
	}
	return false
}

// IsTimedOut returns whether this Member is timed out
func (m *Member) IsTimedOut() bool {
	return m.CommunicationDisabledUntil != nil && m.CommunicationDisabledUntil.After(time.Now())
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
	compiledRoute, err := route.MemberAvatar.Compile(nil, format, size, m.GuildID, m.User.ID, *m.Avatar)
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
func (m *Member) Update(updateGuildMember discord.MemberUpdate, opts ...rest.RequestOpt) (Member, error) {
	member, err := m.Bot.RestServices().GuildService().UpdateMember(m.GuildID, m.User.ID, updateGuildMember, opts...)
	if err != nil {
		return Member{}, err
	}
	return m.Bot.EntityBuilder().CreateMember(m.GuildID, *member, CacheStrategyNoWs), nil
}

// Move moves/kicks the member to/from a voice channel
func (m *Member) Move(channelID discord.Snowflake, opts ...rest.RequestOpt) (Member, error) {
	return m.Update(discord.MemberUpdate{ChannelID: json.NewNullablePtr(channelID)}, opts...)
}

// Kick kicks this Member from the Guild
func (m *Member) Kick(opts ...rest.RequestOpt) error {
	return m.Bot.RestServices().GuildService().RemoveMember(m.GuildID, m.User.ID, opts...)
}

// Ban bans this Member from the Guild
func (m *Member) Ban(deleteMessageDays int, opts ...rest.RequestOpt) error {
	return m.Bot.RestServices().GuildService().AddBan(m.GuildID, m.User.ID, deleteMessageDays, opts...)
}

// Unban unbans this Member from the Guild
func (m *Member) Unban(opts ...rest.RequestOpt) error {
	return m.Bot.RestServices().GuildService().DeleteBan(m.GuildID, m.User.ID, opts...)
}

// AddRole adds a specific Role the Member
func (m *Member) AddRole(roleID discord.Snowflake, opts ...rest.RequestOpt) error {
	return m.Bot.RestServices().GuildService().AddMemberRole(m.GuildID, m.User.ID, roleID, opts...)
}

// RemoveRole removes a specific Role this Member
func (m *Member) RemoveRole(roleID discord.Snowflake, opts ...rest.RequestOpt) error {
	return m.Bot.RestServices().GuildService().RemoveMemberRole(m.GuildID, m.User.ID, roleID, opts...)
}
