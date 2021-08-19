package core

import "github.com/DisgoOrg/disgo/discord"

type SelectMenuInteraction struct {
	*ComponentInteraction
	Data *SelectMenuInteractionData
}

// SelectMenu returns the SelectMenu which issued this SelectMenuInteraction. nil for ephemeral Message(s)
func (i *SelectMenuInteraction) SelectMenu() *SelectMenu {
	return i.Message.SelectMenuByID(i.CustomID())
}

// Values returns the selected values
func (i *SelectMenuInteraction) Values() []string {
	return i.Data.Values
}

// SelectedOptions returns the selected SelectedOption(s)
func (i *SelectMenuInteraction) SelectedOptions() []discord.SelectOption {
	options := make([]discord.SelectOption, len(i.Values()))
	for ii, option := range i.SelectMenu().Component.Options {
		for _, value := range i.Values() {
			if value == option.Value {
				options[ii] = option
				break
			}
		}
	}
	return options
}

type SelectMenuInteractionData struct {
	*ComponentInteractionData
}
