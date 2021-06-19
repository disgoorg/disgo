package main

import "github.com/DisgoOrg/disgo/api"

var rawCmds = []api.CommandCreate{
	{
		Name:              "eval",
		Description:       "runs some go code",
		DefaultPermission: true,
		Options: []api.CommandOption{
			{
				Type:        api.CommandOptionTypeString,
				Name:        "code",
				Description: "the code to eval",
				Required:    true,
			},
		},
	},
	{
		Name:              "test",
		Description:       "test test test test test test",
		DefaultPermission: true,
	},
	{
		Name:              "say",
		Description:       "says what you say",
		DefaultPermission: true,
		Options: []api.CommandOption{
			{
				Type:        api.CommandOptionTypeString,
				Name:        "message",
				Description: "What to say",
				Required:    true,
			},
		},
	},
	{
		Name:              "addrole",
		Description:       "This command adds a role to a member",
		DefaultPermission: true,
		Options: []api.CommandOption{
			{
				Type:        api.CommandOptionTypeUser,
				Name:        "member",
				Description: "The member to add a role to",
				Required:    true,
			},
			{
				Type:        api.CommandOptionTypeRole,
				Name:        "role",
				Description: "The role to add to a member",
				Required:    true,
			},
		},
	},
	{
		Name:              "removerole",
		Description:       "This command removes a role from a member",
		DefaultPermission: true,
		Options: []api.CommandOption{
			{
				Type:        api.CommandOptionTypeUser,
				Name:        "member",
				Description: "The member to removes a role from",
				Required:    true,
			},
			{
				Type:        api.CommandOptionTypeRole,
				Name:        "role",
				Description: "The role to removes from a member",
				Required:    true,
			},
		},
	},
}
