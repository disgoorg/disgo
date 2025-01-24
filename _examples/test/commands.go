package main

import (
	"log/slog"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
)

var commands = []discord.ApplicationCommandCreate{
	discord.SlashCommandCreate{
		Name:        "locale",
		Description: "return the guild & your locale",
	},
	discord.SlashCommandCreate{
		Name:        "test",
		Description: "test",
	},
	discord.SlashCommandCreate{
		Name:        "test2",
		Description: "test",
	},
	discord.SlashCommandCreate{
		Name:        "say",
		Description: "says what you say",
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
	if _, err := client.Rest().SetGuildCommands(client.ApplicationID(), guildID, commands); err != nil {
		slog.Error("error while registering guild commands", slog.Any("err", err))
	}
}
