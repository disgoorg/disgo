package main

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/log"
)

var commands = []discord.ApplicationCommandCreate{
	discord.SlashCommandCreate{
		CommandName: "locale",
		Description: "return the guild & your locale",
	},
	discord.SlashCommandCreate{
		CommandName: "test",
		Description: "test",
	},
	discord.SlashCommandCreate{
		CommandName: "say",
		Description: "says what you say",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionString{
				OptionName:  "message",
				Description: "What to say",
				Required:    true,
			},
			discord.ApplicationCommandOptionBool{
				OptionName:  "ephemeral",
				Description: "ephemeral",
				Required:    true,
			},
		},
	},
}

func registerCommands(client bot.Client) {
	if _, err := client.Rest().SetGuildCommands(client.ApplicationID(), guildID, commands); err != nil {
		log.Fatalf("error while registering guild commands: %s", err)
	}
}
