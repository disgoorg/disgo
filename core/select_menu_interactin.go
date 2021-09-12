package core

type SelectMenuInteraction struct {
	*ComponentInteraction
	SelectMenuInteractionData
}

// SelectMenu returns the SelectMenu which issued this SelectMenuInteraction
func (i *SelectMenuInteraction) SelectMenu() SelectMenu {
	// this should never be nil
	return *i.Message.SelectMenuByID(i.CustomID)
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

type SelectMenuInteractionData struct {
	Values []string
}
