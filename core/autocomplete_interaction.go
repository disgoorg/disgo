package core

import "github.com/DisgoOrg/disgo/discord"

type AutocompleteInteraction struct {
	discord.AutocompleteInteraction
	ResultInteraction
	CommandID           discord.Snowflake
	CommandName         string
	SubCommandName      *string
	SubCommandGroupName *string
	Options             AutocompleteOptionsMap
}

type AutocompleteOptionsMap map[string]discord.AutocompleteOption

func (m AutocompleteOptionsMap) Get(name string) discord.AutocompleteOption {
	if option, ok := m[name]; ok {
		return option
	}
	return nil
}

func (m AutocompleteOptionsMap) GetAll() []discord.AutocompleteOption {
	options := make([]discord.AutocompleteOption, len(m))
	i := 0
	for _, option := range m {
		options[i] = option
		i++
	}
	return options
}

func (m AutocompleteOptionsMap) GetByType(optionType discord.ApplicationCommandOptionType) []discord.AutocompleteOption {
	return m.FindAll(func(option discord.AutocompleteOption) bool {
		return option.Type() == optionType
	})
}

func (m AutocompleteOptionsMap) Find(optionFindFunc func(option discord.AutocompleteOption) bool) discord.AutocompleteOption {
	for _, option := range m {
		if optionFindFunc(option) {
			return option
		}
	}
	return nil
}

func (m AutocompleteOptionsMap) FindAll(optionFindFunc func(option discord.AutocompleteOption) bool) []discord.AutocompleteOption {
	var options []discord.AutocompleteOption
	for _, option := range m {
		if optionFindFunc(option) {
			options = append(options, option)
		}
	}
	return options
}

func (m AutocompleteOptionsMap) FocusedOption() discord.AutocompleteOption {
	return m.Find(func(option discord.AutocompleteOption) bool {
		switch o := option.(type) {
		case discord.AutocompleteOptionString:
			return o.Focused
		case discord.AutocompleteOptionInt:
			return o.Focused
		case discord.AutocompleteOptionBool:
			return o.Focused
		case discord.AutocompleteOptionUser:
			return o.Focused
		case discord.AutocompleteOptionChannel:
			return o.Focused
		case discord.AutocompleteOptionRole:
			return o.Focused
		case discord.AutocompleteOptionMentionable:
			return o.Focused
		case discord.AutocompleteOptionFloat:
			return o.Focused
		default:
			return false
		}
	})
}
