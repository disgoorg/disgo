package main

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/log"
)

var commands = []discord.ApplicationCommandCreate{
	{
		Type:              discord.ApplicationCommandTypeSlash,
		Name:              "autocomplete",
		Description:       "autocomplete",
		DefaultPermission: true,
		Options: []discord.ApplicationCommandOption{
			{
				Type:         discord.ApplicationCommandOptionTypeString,
				Name:         "group",
				Description:  "group",
				Required:     true,
				Autocomplete: true,
			},
			{
				Type:         discord.ApplicationCommandOptionTypeString,
				Name:         "value",
				Description:  "value",
				Required:     true,
				Autocomplete: true,
			},
		},
	},
	{
		Type:              discord.ApplicationCommandTypeSlash,
		Name:              "eval",
		Description:       "runs some go code",
		DefaultPermission: true,
		Options: []discord.ApplicationCommandOption{
			{
				Type:        discord.ApplicationCommandOptionTypeString,
				Name:        "code",
				Description: "the code to eval",
				Required:    true,
			},
		},
	},
	{
		Type:              discord.ApplicationCommandTypeSlash,
		Name:              "test",
		Description:       "test",
		DefaultPermission: true,
		Options: []discord.ApplicationCommandOption{
			{
				Type:        discord.ApplicationCommandOptionTypeBoolean,
				Name:        "ephemeral",
				Description: "ephemeral",
				Required:    true,
			},
		},
	},
	{
		Type:              discord.ApplicationCommandTypeSlash,
		Name:              "say",
		Description:       "says what you say",
		DefaultPermission: true,
		Options: []discord.ApplicationCommandOption{
			{
				Type:        discord.ApplicationCommandOptionTypeString,
				Name:        "message",
				Description: "What to say",
				Required:    true,
			},
			{
				Type:        discord.ApplicationCommandOptionTypeBoolean,
				Name:        "ephemeral",
				Description: "ephemeral",
				Required:    true,
			},
		},
	},
	{
		Type:              discord.ApplicationCommandTypeSlash,
		Name:              "addrole",
		Description:       "This command adds a role to a member",
		DefaultPermission: true,
		Options: []discord.ApplicationCommandOption{
			{
				Type:        discord.ApplicationCommandOptionTypeUser,
				Name:        "member",
				Description: "The member to add a role to",
				Required:    true,
			},
			{
				Type:        discord.ApplicationCommandOptionTypeRole,
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
		Options: []discord.ApplicationCommandOption{
			{
				Type:        discord.ApplicationCommandOptionTypeUser,
				Name:        "member",
				Description: "The member to removes a role from",
				Required:    true,
			},
			{
				Type:        discord.ApplicationCommandOptionTypeRole,
				Name:        "role",
				Description: "The role to removes from a member",
				Required:    true,
			},
		},
	},
}

func registerCommands(bot *core.Bot) {
	cmds, err := bot.SetGuildCommands(guildID, commands)
	if err != nil {
		log.Fatalf("error while registering guild commands: %s", err)
	}

	var cmdsPermissions []discord.ApplicationCommandPermissionsSet
	for _, cmd := range cmds {
		var perms discord.ApplicationCommandPermission
		if cmd.Name == "eval" {
			perms = discord.ApplicationCommandPermission{
				ID:         adminRoleID,
				Type:       discord.ApplicationCommandPermissionTypeRole,
				Permission: true,
			}
		} else {
			perms = discord.ApplicationCommandPermission{
				ID:         testRoleID,
				Type:       discord.ApplicationCommandPermissionTypeRole,
				Permission: true,
			}
			cmdsPermissions = append(cmdsPermissions, discord.ApplicationCommandPermissionsSet{
				ID:          cmd.ID,
				Permissions: []discord.ApplicationCommandPermission{perms},
			})
		}
		cmdsPermissions = append(cmdsPermissions, discord.ApplicationCommandPermissionsSet{
			ID:          cmd.ID,
			Permissions: []discord.ApplicationCommandPermission{perms},
		})
	}
	if _, err = bot.SetGuildCommandsPermissions(guildID, cmdsPermissions); err != nil {
		log.Fatalf("error while setting command permissions: %s", err)
	}
}
