package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

var _ Component = (*SelectMenu)(nil)

// NewSelectMenu builds a new SelectMenu from the provided values
func NewSelectMenu(customID string, placeholder string, minValues int, maxValues int, options ...SelectMenuOption) SelectMenu {
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

type SelectMenu struct {
	discord.Component
	Options []SelectMenuOption `json:"options"`
}

func (m SelectMenu) Type() discord.ComponentType {
	return m.Component.Type
}

// WithCustomID returns a new SelectMenu with the provided customID
func (m SelectMenu) WithCustomID(customID string) SelectMenu {
	m.Component.CustomID = customID
	return m
}

// WithPlaceholder returns a new SelectMenu with the provided placeholder
func (m SelectMenu) WithPlaceholder(placeholder string) SelectMenu {
	m.Component.Placeholder = placeholder
	return m
}

// WithMinValues returns a new SelectMenu with the provided minValue
func (m SelectMenu) WithMinValues(minValue int) SelectMenu {
	m.Component.MinValues = minValue
	return m
}

// WithMaxValues returns a new SelectMenu with the provided maxValue
func (m SelectMenu) WithMaxValues(maxValue int) SelectMenu {
	m.Component.MaxValues = maxValue
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
		SelectOption: discord.SelectOption{
			Label: label,
			Value: value,
		},
	}
}

type SelectMenuOption struct {
	discord.SelectOption
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
func (o SelectMenuOption) WithEmoji(emoji *Emoji) SelectMenuOption {
	o.Emoji = &emoji.Emoji
	return o
}
