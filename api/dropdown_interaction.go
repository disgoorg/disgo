package api

// DropdownInteraction is a specific Interaction when CLicked on Dropdown(s)
type DropdownInteraction struct {
	*ComponentInteraction
	Data *DropdownInteractionData `json:"data,omitempty"`
}

// DropdownInteractionData is the Dropdown data payload
type DropdownInteractionData struct {
	*ComponentInteractionData
	Values []string `json:"values"`
}

// Dropdown returns the Dropdown which issued this DropdownInteraction. nil for ephemeral Message(s)
func (i *DropdownInteraction) Dropdown() *Dropdown {
	return i.Message.DropdownByID(i.Data.CustomID)
}

