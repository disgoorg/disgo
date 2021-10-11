package discord

import "github.com/DisgoOrg/disgo/json"

// NewSelectMenu builds a new SelectMenu from the provided values
func NewSelectMenu(customID string, placeholder string, minValues OptionalInt, maxValues int, disabled bool, options ...SelectMenuOption) SelectMenu {
	return SelectMenu{
		CustomID:    customID,
		Placeholder: placeholder,
		MinValues:   minValues,
		MaxValues:   maxValues,
		Disabled:    disabled,
		Options:     options,
	}
}

var _ Component = (*SelectMenu)(nil)

type SelectMenu struct {
	CustomID    string             `json:"custom_id"`
	Placeholder string             `json:"placeholder,omitempty"`
	MinValues   OptionalInt        `json:"min_values,omitempty"`
	MaxValues   int                `json:"max_values,omitempty"`
	Disabled    bool               `json:"disabled,omitempty"`
	Options     []SelectMenuOption `json:"options,omitempty"`
}

func (m SelectMenu) MarshalJSON() ([]byte, error) {
	v := struct {
		Type ComponentType `json:"type"`
		SelectMenu
	}{
		Type:       m.Type(),
		SelectMenu: m,
	}
	return json.Marshal(v)
}

func (_ SelectMenu) Type() ComponentType {
	return ComponentTypeSelectMenu
}

// WithCustomID returns a new SelectMenu with the provided customID
func (m SelectMenu) WithCustomID(customID string) SelectMenu {
	m.CustomID = customID
	return m
}

// WithPlaceholder returns a new SelectMenu with the provided placeholder
func (m SelectMenu) WithPlaceholder(placeholder string) SelectMenu {
	m.Placeholder = placeholder
	return m
}

// WithMinValues returns a new SelectMenu with the provided minValue
func (m SelectMenu) WithMinValues(minValue OptionalInt) SelectMenu {
	m.MinValues = minValue
	return m
}

// WithMaxValues returns a new SelectMenu with the provided maxValue
func (m SelectMenu) WithMaxValues(maxValue int) SelectMenu {
	m.MaxValues = maxValue
	return m
}

// SetOptions returns a new SelectMenu with the provided SelectMenuOption(s)
func (m SelectMenu) SetOptions(options ...SelectMenuOption) SelectMenu {
	m.Options = options
	return m
}

// AddOptions returns a new SelectMenu with the provided SelectMenuOption(s) added
func (m SelectMenu) AddOptions(options ...SelectMenuOption) SelectMenu {
	m.Options = append(m.Options, options...)
	return m
}

// SetOption returns a new SelectMenu with the SelectMenuOption which has the value replaced
func (m SelectMenu) SetOption(value string, option SelectMenuOption) SelectMenu {
	for i, o := range m.Options {
		if o.Value == value {
			m.Options[i] = option
			break
		}
	}
	return m
}

// RemoveOption returns a new SelectMenu with the provided SelectMenuOption at the index removed
func (m SelectMenu) RemoveOption(index int) SelectMenu {
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

// SelectMenuOption represents an option in a SelectMenu
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
