package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

type ApplicationCommandAutocompleteInteraction struct {
	*SlashCommandInteraction
}

func (e *ApplicationCommandAutocompleteInteraction) Result(choices []discord.ApplicationCommandOptionChoice, opts ...rest.RequestOpt) rest.Error {
	return e.Respond(discord.InteractionCallbackTypeApplicationCommandAutoCompleteResult, discord.ApplicationCommandAutoCompleteResult{Choices: choices}, opts...)
}

func (e *ApplicationCommandAutocompleteInteraction) ResultMap(resultMap map[string]interface{}, opts ...rest.RequestOpt) rest.Error {
	choices := make([]discord.ApplicationCommandOptionChoice, len(resultMap))
	i := 0
	for name, value := range resultMap {
		choices[i] = discord.ApplicationCommandOptionChoice{
			Name:  name,
			Value: value,
		}
		i++
	}
	return e.Result(choices, opts...)
}

func (e *ApplicationCommandAutocompleteInteraction) GetFocusedOption() ApplicationCommandOption {
	return *e.Options.Find(func(option ApplicationCommandOption) bool {
		return option.Focused
	})
}
