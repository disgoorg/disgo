package core

type ApplicationCommandAutocompleteInteraction struct {
	*ApplicationCommandOptionsInteraction
	AutoCompleteInteractionResponses
}

func (i *ApplicationCommandAutocompleteInteraction) FocusedOption() ApplicationCommandOption {
	return *i.Options.Find(func(option ApplicationCommandOption) bool {
		return option.Focused
	})
}
