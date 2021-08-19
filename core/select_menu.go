package core

import "github.com/DisgoOrg/disgo/discord"

// NewSelectMenu builds a new SelectMenu from the provided values
func NewSelectMenu(customID string, placeholder string, minValues int, maxValues int, options ...SelectOption) SelectMenu {
	return SelectMenu{
		Component: discord.Component{
			Type:        discord.ComponentTypeSelectMenu,
			CustomID:    customID,
			Placeholder: placeholder,
			MinValues:   minValues,
			MaxValues:   maxValues,
		},
		Options: options,
	}
}

// SelectMenu is a Component which lets the User select from various options
type SelectMenu struct {
	discord.Component
	Options []SelectOption `json:"options"`
}

// Type returns the ComponentType of this Component
func (m SelectMenu) Type() discord.ComponentType {
	return m.Component.Type
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
func (m SelectMenu) WithMinValues(minValue int) SelectMenu {
	m.MinValues = minValue
	return m
}

// WithMaxValues returns a new SelectMenu with the provided maxValue
func (m SelectMenu) WithMaxValues(maxValue int) SelectMenu {
	m.MaxValues = maxValue
	return m
}

// SetOptions returns a new SelectMenu with the provided SelectOption(s)
func (m SelectMenu) SetOptions(options ...SelectOption) SelectMenu {
	m.Options = options
	return m
}

// AddOptions returns a new SelectMenu with the provided SelectOption(s) added
func (m SelectMenu) AddOptions(options ...SelectOption) SelectMenu {
	m.Options = append(m.Options, options...)
	return m
}

// SetOption returns a new SelectMenu with the SelectOption which has the value replaced
func (m SelectMenu) SetOption(value string, option SelectOption) SelectMenu {
	for i, o := range m.Options {
		if o.Value == value {
			m.Options[i] = option
			break
		}
	}
	return m
}

// RemoveOption returns a new SelectMenu with the provided SelectOption at the index removed
func (m SelectMenu) RemoveOption(index int) SelectMenu {
	if len(m.Options) > index {
		m.Options = append(m.Options[:index], m.Options[index+1:]...)
	}
	return m
}

// NewSelectOption builds a new SelectOption
func NewSelectOption(label string, value string) SelectOption {
	return SelectOption{
		SelectOption: discord.SelectOption{
			Label: label,
			Value: value,
		},
	}
}

type SelectOption struct {
	discord.SelectOption
}

// WithLabel returns a new SelectOption with the provided label
func (o SelectOption) WithLabel(label string) SelectOption {
	o.Label = label
	return o
}

// WithValue returns a new SelectOption with the provided value
func (o SelectOption) WithValue(value string) SelectOption {
	o.Value = value
	return o
}

// WithDescription returns a new SelectOption with the provided description
func (o SelectOption) WithDescription(description string) SelectOption {
	o.Description = description
	return o
}

// WithDefault returns a new SelectOption as default/non-default
func (o SelectOption) WithDefault(defaultOption bool) SelectOption {
	o.Default = defaultOption
	return o
}

// WithEmoji returns a new SelectOption with the provided Emoji
func (o SelectOption) WithEmoji(emoji *Emoji) SelectOption {
	o.Emoji = &emoji.Emoji
	return o
}
