package api

// SelectMenuInteraction is a specific Interaction when CLicked on SelectMenu(s)
type SelectMenuInteraction struct {
	*ComponentInteraction
	Data *SelectMenuInteractionData `json:"data,omitempty"`
}

// SelectMenuInteractionData is the SelectMenu data payload
type SelectMenuInteractionData struct {
	*ComponentInteractionData
	Values []string `json:"values"`
}

// SelectMenu returns the SelectMenu which issued this SelectMenuInteraction. nil for ephemeral Message(s)
func (i *SelectMenuInteraction) SelectMenu() *SelectMenu {
	return i.Message.SelectMenuByID(i.Data.CustomID)
}

