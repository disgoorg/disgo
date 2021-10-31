package discord

import "github.com/DisgoOrg/disgo/json"

// NewSelectMenu builds a new SelectMenuComponent from the provided values
func NewSelectMenu(customID string, placeholder string, minValues OptionalInt, maxValues int, disabled bool, options ...SelectMenuOption) SelectMenuComponent {
	return SelectMenuComponent{
		CustomID:    customID,
		Placeholder: placeholder,
		MinValues:   minValues,
		MaxValues:   maxValues,
		Disabled:    disabled,
		Options:     options,
	}
}

var _ Component = (*SelectMenuComponent)(nil)

type SelectMenuComponent struct {
	CustomID    string             `json:"custom_id"`
	Placeholder string             `json:"placeholder,omitempty"`
	MinValues   OptionalInt        `json:"min_values,omitempty"`
	MaxValues   int                `json:"max_values,omitempty"`
	Disabled    bool               `json:"disabled,omitempty"`
	Options     []SelectMenuOption `json:"options,omitempty"`
}

func (m SelectMenuComponent) MarshalJSON() ([]byte, error) {
	type selectMenu SelectMenuComponent
	v := struct {
		Type ComponentType `json:"type"`
		selectMenu
	}{
		Type:       m.Type(),
		selectMenu: selectMenu(m),
	}
	return json.Marshal(v)
}

func (_ SelectMenuComponent) Type() ComponentType {
	return ComponentTypeSelectMenu
}

// WithCustomID returns a new SelectMenuComponent with the provided customID
func (m SelectMenuComponent) WithCustomID(customID string) SelectMenuComponent {
	m.CustomID = customID
	return m
}

// WithPlaceholder returns a new SelectMenuComponent with the provided placeholder
func (m SelectMenuComponent) WithPlaceholder(placeholder string) SelectMenuComponent {
	m.Placeholder = placeholder
	return m
}

// WithMinValues returns a new SelectMenuComponent with the provided minValue
func (m SelectMenuComponent) WithMinValues(minValue OptionalInt) SelectMenuComponent {
	m.MinValues = minValue
	return m
}

// WithMaxValues returns a new SelectMenuComponent with the provided maxValue
func (m SelectMenuComponent) WithMaxValues(maxValue int) SelectMenuComponent {
	m.MaxValues = maxValue
	return m
}

// SetOptions returns a new SelectMenuComponent with the provided SelectMenuOption(s)
func (m SelectMenuComponent) SetOptions(options ...SelectMenuOption) SelectMenuComponent {
	m.Options = options
	return m
}

// AddOptions returns a new SelectMenuComponent with the provided SelectMenuOption(s) added
func (m SelectMenuComponent) AddOptions(options ...SelectMenuOption) SelectMenuComponent {
	m.Options = append(m.Options, options...)
	return m
}

// SetOption returns a new SelectMenuComponent with the SelectMenuOption which has the value replaced
func (m SelectMenuComponent) SetOption(value string, option SelectMenuOption) SelectMenuComponent {
	for i, o := range m.Options {
		if o.Value == value {
			m.Options[i] = option
			break
		}
	}
	return m
}

// RemoveOption returns a new SelectMenuComponent with the provided SelectMenuOption at the index removed
func (m SelectMenuComponent) RemoveOption(index int) SelectMenuComponent {
	if len(m.Options) > index {
		m.Options = append(m.Options[:index], m.Options[index+1:]...)
	}
	return m
}

// NewSelectMenuOption builds a new SelectMenuOption
func NewSelectMenuOption(label string, value string) SelectMenuOption {
	return SelectMenuOption{
		Label: label,
		Value: value,
	}
}

// SelectMenuOption represents an option in a SelectMenuComponent
type SelectMenuOption struct {
	Label       string          `json:"label"`
	Value       string          `json:"value"`
	Description string          `json:"description,omitempty"`
	Emoji       *ComponentEmoji `json:"emoji,omitempty"`
	Default     bool            `json:"default,omitempty"`
}

// WithLabel returns a new SelectMenuOption with the provided label
func (o SelectMenuOption) WithLabel(label string) SelectMenuOption {
	o.Label = label
	return o
}

// WithValue returns a new SelectMenuOption with the provided value
func (o SelectMenuOption) WithValue(value string) SelectMenuOption {
	o.Value = value
	return o
}

// WithDescription returns a new SelectMenuOption with the provided description
func (o SelectMenuOption) WithDescription(description string) SelectMenuOption {
	o.Description = description
	return o
}

// WithDefault returns a new SelectMenuOption as default/non-default
func (o SelectMenuOption) WithDefault(defaultOption bool) SelectMenuOption {
	o.Default = defaultOption
	return o
}

// WithEmoji returns a new SelectMenuOption with the provided Emoji
func (o SelectMenuOption) WithEmoji(emoji ComponentEmoji) SelectMenuOption {
	o.Emoji = &emoji
	return o
}
