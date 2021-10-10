package core

type SlashCommandInteractionFilter func(slashCommandInteraction *SlashCommandInteraction) bool

type SlashCommandInteraction struct {
	*ApplicationCommandOptionsInteraction
	CreateInteractionResponses
}
