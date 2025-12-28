package main

import (
	"github.com/disgoorg/omit"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var commands = []discord.ApplicationCommandCreate{
	discord.SlashCommandCreate{
		Name:        "test",
		Description: "test",
	},
}

func handleTestCommand(_ discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	if err := e.CreateMessage(discord.MessageCreate{
		Content: "Test command received!",
	}); err != nil {
		return err
	}

	_, err := e.UpdateInteractionResponse(discord.MessageUpdate{
		Content: omit.Ptr("Test command response updated!"),
	})
	return err
}
