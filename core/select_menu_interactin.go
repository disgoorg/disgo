package core

import "github.com/DisgoOrg/disgo/rest"

type SelectMenuInteraction struct {
	*ComponentInteraction
	SelectMenuInteractionData
}

// SelectMenu returns the SelectMenu which issued this SelectMenuInteraction
func (i *SelectMenuInteraction) SelectMenu() SelectMenu {
	// this should never be nil
	return *i.Message.SelectMenuByID(i.CustomID)
}

// UpdateSelectMenu updates the used SelectMenu with a new SelectMenu
func (i *SelectMenuInteraction) UpdateSelectMenu(selectMenu SelectMenu, opts ...rest.RequestOpt) rest.Error {
	return i.UpdateComponent(selectMenu, opts...)
}

// SelectedOptions returns the selected SelectMenuOption(s)
func (i *SelectMenuInteraction) SelectedOptions() []SelectMenuOption {
	options := make([]SelectMenuOption, len(i.Values))
	for ii, option := range i.SelectMenu().Options {
		for _, value := range i.Values {
			if value == option.Value {
				options[ii] = option
				break
			}
		}
	}
	return options
}

// SelectMenuInteractionData is data specifically from the SelectMenuInteraction
type SelectMenuInteractionData struct {
	Values []string
}
