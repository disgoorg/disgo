package main

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/log"
)

var commands = []discord.ApplicationCommandCreate{
	discord.SlashCommandCreate{
		CommandName:       "locale",
		Description:       "return the guild & your locale",
		DefaultPermission: true,
	},
	discord.SlashCommandCreate{
		CommandName:       "test",
		Description:       "test",
		DefaultPermission: true,
	},
	discord.SlashCommandCreate{
		CommandName:       "say",
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
}

func registerCommands(bot core.Bot) {
	if _, err := bot.SetGuildCommands(guildID, commands); err != nil {
		log.Fatalf("error while registering guild commands: %s", err)
	}
}
