package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

type SlashCommandOption interface {
	discord.SlashCommandOption
}

type SlashCommandOptionSubCommand struct {
	discord.SlashCommandOptionSubCommand
	Options []SlashCommandOption
}

type SlashCommandOptionSubCommandGroup struct {
	discord.SlashCommandOptionSubCommandGroup
	Options []SlashCommandOptionSubCommand
}

type SlashCommandOptionString struct {
	discord.SlashCommandOptionString
	Resolved *SlashCommandResolved
}

func (o *SlashCommandOptionString) MentionedUsers() []*User {
	matches := discord.MentionTypeUser.FindAllStringSubmatch(o.Value, -1)
	users := make([]*User, len(matches))
	if matches == nil {
		return nil
	}
	for i := range matches {
		users[i] = o.Resolved.Users[discord.Snowflake(matches[i][1])]
	}
	return users
}

func (o *SlashCommandOptionString) MentionedMembers() []*Member {
	matches := discord.MentionTypeUser.FindAllStringSubmatch(o.Value, -1)
	members := make([]*Member, len(matches))
	if matches == nil {
		return nil
	}
	for i := range matches {
		if member, ok := o.Resolved.Members[discord.Snowflake(matches[i][1])]; ok {
			members[i] = member
		}
	}
	return members
}

func (o *SlashCommandOptionString) MentionedChannels() []Channel {
	matches := discord.MentionTypeChannel.FindAllStringSubmatch(o.Value, -1)
	channels := make([]Channel, len(matches))
	if matches == nil {
		return nil
	}
	for i := range matches {
		if channel, ok := o.Resolved.Channels[discord.Snowflake(matches[i][1])]; ok {
			channels[i] = channel
		}
	}
	return channels
}

func (o *SlashCommandOptionString) MentionedRoles() []*Role {
	matches := discord.MentionTypeRole.FindAllStringSubmatch(o.Value, -1)
	roles := make([]*Role, len(matches))
	if matches == nil {
		return nil
	}
	for i := range matches {
		if role, ok := o.Resolved.Roles[discord.Snowflake(matches[i][1])]; ok {
			roles[i] = role
		}
	}
	return roles
}

type SlashCommandOptionInt struct {
	discord.SlashCommandOptionInt
}

type SlashCommandOptionBool struct {
	discord.SlashCommandOptionBool
}

type SlashCommandOptionUser struct {
	discord.SlashCommandOptionUser
	Resolved *SlashCommandResolved
}

func (o SlashCommandOptionUser) User() *User {
	return o.Resolved.Users[o.Value]
}

func (o SlashCommandOptionUser) Member() *Member {
	return o.Resolved.Members[o.Value]
}

type SlashCommandOptionChannel struct {
	discord.SlashCommandOptionChannel
	Resolved *SlashCommandResolved
}

func (o SlashCommandOptionChannel) Channel() Channel {
	return o.Resolved.Channels[o.Value]
}

type SlashCommandOptionRole struct {
	discord.SlashCommandOptionRole
	Resolved *SlashCommandResolved
}

func (o SlashCommandOptionRole) Role() *Role {
	return o.Resolved.Roles[o.Value]
}

type SlashCommandOptionMentionable struct {
	discord.SlashCommandOptionMentionable
	Resolved *SlashCommandResolved
}

func (o SlashCommandOptionMentionable) User() *User {
	return o.Resolved.Users[o.Value]
}

func (o SlashCommandOptionMentionable) Member() *Member {
	return o.Resolved.Members[o.Value]
}

func (o SlashCommandOptionMentionable) Role() *Role {
	return o.Resolved.Roles[o.Value]
}

type SlashCommandOptionFloat struct {
	discord.SlashCommandOptionFloat
}
