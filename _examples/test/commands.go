package main

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/log"
)

var commands = []discord.ApplicationCommandCreate{
	discord.SlashCommandCreate{
		Name:              "eval",
		Description:       "runs some go code",
		DefaultPermission: true,
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionString{
				Name:        "code",
				Description: "the code to eval",
				Required:    true,
			},
		},
	},
	discord.SlashCommandCreate{
		Name:              "test",
		Description:       "test",
		DefaultPermission: true,
	},
	discord.SlashCommandCreate{
		Name:              "say",
		Description:       "says what you say",
		DefaultPermission: true,
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionString{
				Name:        "message",
				Description: "What to say",
				Required:    true,
			},
			discord.ApplicationCommandOptionBool{
				Name:        "ephemeral",
				Description: "ephemeral",
				Required:    true,
			},
		},
	},
	discord.SlashCommandCreate{
		Name:              "addrole",
		Description:       "This command adds a role to a member",
		DefaultPermission: true,
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionUser{
				Name:        "member",
				Description: "The member to add a role to",
				Required:    true,
			},
			discord.ApplicationCommandOptionRole{
				Name:        "role",
				Description: "The role to add to a member",
				Required:    true,
			},
		},
	},
	discord.SlashCommandCreate{
		Name:              "removerole",
		Description:       "This command removes a role from a member",
		DefaultPermission: true,
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionUser{
				Name:        "member",
				Description: "The member to removes a role from",
				Required:    true,
			},
			discord.ApplicationCommandOptionRole{
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
		if c, ok := cmd.(core.SlashCommand); ok {
			if c.Name == "eval" {
				perms = discord.ApplicationCommandPermissionRole{
					ID:         adminRoleID,
					Permission: true,
				}
			} else {
				perms = discord.ApplicationCommandPermissionUser{
					ID:         testRoleID,
					Permission: true,
				}
				cmdsPermissions = append(cmdsPermissions, discord.ApplicationCommandPermissionsSet{
					ID:          c.ID,
					Permissions: []discord.ApplicationCommandPermission{perms},
				})
			}
			cmdsPermissions = append(cmdsPermissions, discord.ApplicationCommandPermissionsSet{
				ID:          c.ID,
				Permissions: []discord.ApplicationCommandPermission{perms},
			})
		}
	}
	if _, err = bot.SetGuildCommandsPermissions(guildID, cmdsPermissions); err != nil {
		log.Fatalf("error while setting command permissions: %s", err)
	}
}
