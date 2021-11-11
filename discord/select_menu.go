package discord

import "github.com/DisgoOrg/disgo/json"

// NewSelectMenu builds a new SelectMenuComponent from the provided values
func NewSelectMenu(customID CustomID, placeholder string, minValues json.NullInt, maxValues json.NullInt, disabled bool, options ...SelectMenuOption) SelectMenuComponent {
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
	CustomID    CustomID           `json:"custom_id"`
	Placeholder string             `json:"placeholder,omitempty"`
	MinValues   json.NullInt       `json:"min_values,omitempty"`
	MaxValues   json.NullInt       `json:"max_values,omitempty"`
	Disabled    bool               `json:"disabled,omitempty"`
	Options     []SelectMenuOption `json:"options,omitempty"`
}

func (c SelectMenuComponent) MarshalJSON() ([]byte, error) {
	type selectMenuComponent SelectMenuComponent
	return json.Marshal(struct {
		Type ComponentType `json:"type"`
		selectMenuComponent
	}{
		Type:                c.Type(),
		selectMenuComponent: selectMenuComponent(c),
	})
}

func (c SelectMenuComponent) Type() ComponentType {
	return ComponentTypeSelectMenu
}

func (c SelectMenuComponent) ID() CustomID {
	return c.CustomID
}

func (c SelectMenuComponent) component()            {}
func (c SelectMenuComponent) interactiveComponent() {}

// WithCustomID returns a new SelectMenuComponent with the provided customID
func (c SelectMenuComponent) WithCustomID(customID CustomID) SelectMenuComponent {
	c.CustomID = customID
	return c
}

// WithPlaceholder returns a new SelectMenuComponent with the provided placeholder
func (c SelectMenuComponent) WithPlaceholder(placeholder string) SelectMenuComponent {
	c.Placeholder = placeholder
	return c
}

// WithMinValues returns a new SelectMenuComponent with the provided minValue
func (c SelectMenuComponent) WithMinValues(minValue int) SelectMenuComponent {
	c.MinValues = *json.NewInt(minValue)
	return c
}

// WithMaxValues returns a new SelectMenuComponent with the provided maxValue
func (c SelectMenuComponent) WithMaxValues(maxValue int) SelectMenuComponent {
	c.MaxValues = *json.NewInt(maxValue)
	return c
}

// SetOptions returns a new SelectMenuComponent with the provided SelectMenuOption(s)
func (c SelectMenuComponent) SetOptions(options ...SelectMenuOption) SelectMenuComponent {
	c.Options = options
	return c
}

// AddOptions returns a new SelectMenuComponent with the provided SelectMenuOption(s) added
func (c SelectMenuComponent) AddOptions(options ...SelectMenuOption) SelectMenuComponent {
	c.Options = append(c.Options, options...)
	return c
}

// SetOption returns a new SelectMenuComponent with the SelectMenuOption which has the value replaced
func (c SelectMenuComponent) SetOption(value string, option SelectMenuOption) SelectMenuComponent {
	for i, o := range c.Options {
		if o.Value == value {
			c.Options[i] = option
			break
		}
	}
	return c
}

// RemoveOption returns a new SelectMenuComponent with the provided SelectMenuOption at the index removed
func (c SelectMenuComponent) RemoveOption(index int) SelectMenuComponent {
	if len(c.Options) > index {
		c.Options = append(c.Options[:index], c.Options[index+1:]...)
	}
	return c
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
