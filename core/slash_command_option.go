package core

import "github.com/DisgoOrg/disgo/discord"

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

// TODO: add methods for resolving options

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
