package api

// SelectMenuInteraction is a specific Interaction when CLicked on SelectMenu(s)
type SelectMenuInteraction struct {
	*ComponentInteraction
	Data *SelectMenuInteractionData `json:"data,omitempty"`
}

// SelectMenu returns the SelectMenu which issued this SelectMenuInteraction. nil for ephemeral Message(s)
func (i *SelectMenuInteraction) SelectMenu() *SelectMenu {
	if i.Message.IsEphemeral() {
		return nil
	}
	return i.Message.SelectMenuByID(i.CustomID())
}

// Values returns the selected values
func (i *SelectMenuInteraction) Values() []string {
	return i.Data.Values
}

// SelectedOptions returns the selected SelectedOption(s)
func (i *SelectMenuInteraction) SelectedOptions() []SelectOption {
	if i.Message.IsEphemeral() {
		return nil
	}
	options := make([]SelectOption, len(i.Values()))
	for ii, option := range i.SelectMenu().Options {
		for _, value := range i.Values() {
			if value == option.Value {
				options[ii] = option
				break
			}
		}
	}
	return options
}

// SelectMenuInteractionData is the SelectMenu data payload
type SelectMenuInteractionData struct {
	*ComponentInteractionData
	Values []string `json:"values"`
}
