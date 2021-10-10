package core

import "github.com/DisgoOrg/disgo/discord"

type SlashCommandInteractionFilter func(slashCommandInteraction *SlashCommandInteraction) bool

type SlashCommandInteraction struct {
	*ApplicationCommandOptionsInteraction
	CreateInteractionResponses
}
