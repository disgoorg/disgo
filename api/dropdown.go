package api

func NewDropdown(placeholder string, minValues int, maxValues int, options ...*DropdownOption) *Dropdown {
	return &Dropdown{
		ComponentImpl: newComponentImpl(ComponentTypeDropdown),
		Placeholder:   placeholder,
		MinValues:     minValues,
		MaxValues:     maxValues,
		Options:       options,
	}
}

// Dropdown ...
type Dropdown struct {
	ComponentImpl
	Placeholder string            `json:"placeholder"`
	MinValues   int               `json:"min_values,omitempty"`
	MaxValues   int               `json:"max_values,omitempty"`
	Options     []*DropdownOption `json:"options"`
}

func NewDropdownOption(label string, value string) *DropdownOption {
	return &DropdownOption{
		Label: label,
		Value: value,
	}
}

type DropdownOption struct {
	Label string `json:"label"`
	Value string `json:"value"`
}
