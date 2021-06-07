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

