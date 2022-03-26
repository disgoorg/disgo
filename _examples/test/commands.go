package main

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/log"
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

func registerCommands(client bot.Client) {
	if _, err := client.Rest().Applications().SetGuildCommands(client.ApplicationID(), guildID, commands); err != nil {
		log.Fatalf("error while registering guild commands: %s", err)
	}
}
