package api

func NewSelect(customID string, placeholder string, minValues int, maxValues int, options ...*SelectOption) *Select {
	return &Select{
		ComponentImpl: ComponentImpl{
			ComponentType: ComponentTypeSelect,
		},
		CustomID:    customID,
		Placeholder: placeholder,
		Options:     options,
		MinValues:   minValues,
		MaxValues:   maxValues,
	}
}

type Select struct {
	ComponentImpl
	CustomID    string          `json:"custom_id,omitempty"`
	Placeholder string          `json:"placeholder,omitempty"`
	Options     []*SelectOption `json:"options,omitempty"`
	MinValues   int             `json:"min_values,omitempty"`
	MaxValues   int             `json:"max_values,omitempty"`
}

type SelectOption struct {
	Label       string      `json:"label"`
	Value       interface{} `json:"value"`
	Default     bool        `json:"default,omitempty"`
	Description string      `json:"description"`
	Emoji       *Emoji      `json:"emoji,omitempty"`
}
