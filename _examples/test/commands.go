package main

import "github.com/DisgoOrg/disgo/discord"

var commands = []discord.ApplicationCommandCreate{
	{
		Type:              discord.ApplicationCommandTypeSlash,
		Name:              "eval",
		Description:       "runs some go code",
		DefaultPermission: true,
		Options: []discord.SlashCommandOption{
			{
				Type:        discord.CommandOptionTypeString,
				Name:        "code",
				Description: "the code to eval",
				Required:    true,
			},
		},
	},
	{
		Type:              discord.ApplicationCommandTypeSlash,
		Name:              "test",
		Description:       "test test test test test test",
		DefaultPermission: true,
	},
	{
		Type:              discord.ApplicationCommandTypeSlash,
		Name:              "say",
		Description:       "says what you say",
		DefaultPermission: true,
		Options: []discord.SlashCommandOption{
			{
				Type:        discord.CommandOptionTypeString,
				Name:        "message",
				Description: "What to say",
				Required:    true,
			},
		},
	},
	{
		Type:              discord.ApplicationCommandTypeSlash,
		Name:              "addrole",
		Description:       "This command adds a role to a member",
		DefaultPermission: true,
		Options: []discord.SlashCommandOption{
			{
				Type:        discord.CommandOptionTypeUser,
				Name:        "member",
				Description: "The member to add a role to",
				Required:    true,
			},
			{
				Type:        discord.CommandOptionTypeRole,
				Name:        "role",
				Description: "The role to add to a member",
				Required:    true,
			},
		},
	},
	{
		Type:              discord.ApplicationCommandTypeSlash,
		Name:              "removerole",
		Description:       "This command removes a role from a member",
		DefaultPermission: true,
		Options: []discord.SlashCommandOption{
			{
				Type:        discord.CommandOptionTypeUser,
				Name:        "member",
				Description: "The member to removes a role from",
				Required:    true,
			},
			{
				Type:        discord.CommandOptionTypeRole,
				Name:        "role",
				Description: "The role to removes from a member",
				Required:    true,
			},
		},
	},
}
