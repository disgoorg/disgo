package api

// NewDropdown builds a new Dropdown from the provided values
func NewDropdown(customID string, placeholder string, minValues int, maxValues int, options ...DropdownOption) *Dropdown {
	return &Dropdown{
		ComponentImpl: newComponentImpl(ComponentTypeDropdown),
		CustomID:      customID,
		Placeholder:   placeholder,
		MinValues:     minValues,
		MaxValues:     maxValues,
		Options:       options,
	}
}

// Dropdown is a Component which lets the User select from various options
type Dropdown struct {
	ComponentImpl
	CustomID    string           `json:"custom_id"`
	Placeholder string           `json:"placeholder"`
	MinValues   int              `json:"min_values,omitempty"`
	MaxValues   int              `json:"max_values,omitempty"`
	Options     []DropdownOption `json:"options"`
}

// NewDropdownOption builds a new DropdownOption
func NewDropdownOption(label string, value string) DropdownOption {
	return DropdownOption{
		Label: label,
		Value: value,
	}
}

// DropdownOption represents an option in a Dropdown
type DropdownOption struct {
	Label string `json:"label"`
	Value string `json:"value"`
}
