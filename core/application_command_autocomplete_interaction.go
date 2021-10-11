package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

type ApplicationCommandAutocompleteInteractionFilter func(applicationCommandAutocompleteInteraction *ApplicationCommandAutocompleteInteraction) bool

type ApplicationCommandAutocompleteInteraction struct {
	*Interaction
	CommandID           discord.Snowflake
	CommandName         string
	SubCommandName      *string
	SubCommandGroupName *string
	Options             ApplicationCommandAutocompleteOptionsMap
	Resolved            *Resolved
}

// CommandPath returns the ApplicationCommand path
func (i *ApplicationCommandAutocompleteInteraction) CommandPath() string {
	path := i.CommandName
	if name := i.SubCommandName; name != nil {
		path += "/" + *name
	}
	if name := i.SubCommandGroupName; name != nil {
		path += "/" + *name
	}
	return path
}

func (i *ApplicationCommandAutocompleteInteraction) Result(choices []discord.ApplicationCommandOptionChoice, opts ...rest.RequestOpt) error {
	return i.Respond(discord.InteractionCallbackTypeApplicationCommandAutoCompleteResult, discord.ApplicationCommandAutoCompleteResult{Choices: choices}, opts...)
}

func (i *ApplicationCommandAutocompleteInteraction) ResultMap(resultMap map[string]string, opts ...rest.RequestOpt) error {
	choices := make([]discord.ApplicationCommandOptionChoice, len(resultMap))
	ii := 0
	for name, value := range resultMap {
		choices[ii] = discord.ApplicationCommandOptionChoice{
			Name:  name,
			Value: value,
		}
		ii++
	}
	return i.Result(choices, opts...)
}

func (i *ApplicationCommandAutocompleteInteraction) FocusedOption() ApplicationCommandAutocompleteOption {
	return *i.Options.Find(func(option ApplicationCommandAutocompleteOption) bool {
		return option.Focused
	})
}

type ApplicationCommandAutocompleteOptionsMap map[string]ApplicationCommandAutocompleteOption

func (m ApplicationCommandAutocompleteOptionsMap) Get(name string) *ApplicationCommandAutocompleteOption {
	if option, ok := m[name]; ok {
		return &option
	}
	return nil
}

func (m ApplicationCommandAutocompleteOptionsMap) GetAll() []ApplicationCommandAutocompleteOption {
	options := make([]ApplicationCommandAutocompleteOption, len(m))
	i := 0
	for _, option := range m {
		options[i] = option
		i++
	}
	return options
}

func (m ApplicationCommandAutocompleteOptionsMap) GetByType(optionType discord.ApplicationCommandOptionType) []ApplicationCommandAutocompleteOption {
	return m.FindAll(func(option ApplicationCommandAutocompleteOption) bool {
		return option.Type == optionType
	})
}

func (m ApplicationCommandAutocompleteOptionsMap) Find(optionFindFunc func(option ApplicationCommandAutocompleteOption) bool) *ApplicationCommandAutocompleteOption {
	for _, option := range m {
		if optionFindFunc(option) {
			return &option
		}
	}
	return nil
}

func (m ApplicationCommandAutocompleteOptionsMap) FindAll(optionFindFunc func(option ApplicationCommandAutocompleteOption) bool) []ApplicationCommandAutocompleteOption {
	var options []ApplicationCommandAutocompleteOption
	for _, option := range m {
		if optionFindFunc(option) {
			options = append(options, option)
		}
	}
	return options
}

func (m ApplicationCommandAutocompleteOptionsMap) FocusedOption() ApplicationCommandAutocompleteOption {
	return *m.Find(func(option ApplicationCommandAutocompleteOption) bool {
		return option.Focused
	})
}

// ApplicationCommandAutocompleteOption holds info about an ApplicationCommandAutocompleteOption.Value
type ApplicationCommandAutocompleteOption struct {
	ApplicationCommandOption
	Focused  bool
}
