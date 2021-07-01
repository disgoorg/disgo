package api

// NewSelectMenu builds a new SelectMenu from the provided values
func NewSelectMenu(customID string, placeholder string, minValues int, maxValues int, options ...SelectMenuOption) SelectMenu {
	return SelectMenu{
		ComponentImpl: newComponentImpl(ComponentTypeSelectMenu),
		CustomID:      customID,
		Placeholder:   placeholder,
		MinValues:     minValues,
		MaxValues:     maxValues,
		Options:       options,
	}
}

// SelectMenu is a Component which lets the User select from various options
type SelectMenu struct {
	ComponentImpl
	CustomID    string           `json:"custom_id"`
	Placeholder string           `json:"placeholder"`
	MinValues   int              `json:"min_values,omitempty"`
	MaxValues   int              `json:"max_values,omitempty"`
	Options     []SelectMenuOption `json:"options"`
}

// NewSelectMenuOption builds a new SelectMenuOption
func NewSelectMenuOption(label string, value string) SelectMenuOption {
	return SelectMenuOption{
		Label: label,
		Value: value,
	}
}

// SelectMenuOption represents an option in a SelectMenu
type SelectMenuOption struct {
	Label       string `json:"label"`
	Value       string `json:"value"`
	Description string `json:"description,omitempty"`
	Default     bool   `json:"default"`
	Emoji       *Emoji `json:"emoji"`
}
