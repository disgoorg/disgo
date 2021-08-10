package discord

// ComponentType defines different Component(s)
type ComponentType int

// Supported ComponentType(s)
//goland:noinspection GoUnusedConst
const (
	ComponentTypeActionRow = iota + 1
	ComponentTypeButton
	ComponentTypeSelectMenu
)

// ButtonStyle defines how the Button looks like (https://discord.com/assets/7bb017ce52cfd6575e21c058feb3883b.png)
type ButtonStyle int

// Supported ButtonStyle(s)
const (
	ButtonStylePrimary = iota + 1
	ButtonStyleSecondary
	ButtonStyleSuccess
	ButtonStyleDanger
	ButtonStyleLink
)

// UnmarshalComponent is used for easier unmarshalling of different Component(s)
type UnmarshalComponent struct {
	Type ComponentType `json:"type"`

	// Button && SelectMenu
	CustomID string `json:"custom_id"`

	// Button
	Style    ButtonStyle `json:"style"`
	Label    string      `json:"label"`
	Emoji    *Emoji      `json:"emoji"`
	URL      string      `json:"url"`
	Disabled bool        `json:"disabled"`

	// ActionRow
	Components []UnmarshalComponent `json:"components"`

	// SelectMenu
	Placeholder string         `json:"placeholder"`
	MinValues   int            `json:"min_values"`
	MaxValues   int            `json:"max_values"`
	Options     []SelectOption `json:"options"`
}

// SelectOption represents an option in a SelectMenu
type SelectOption struct {
	Label       string `json:"label"`
	Value       string `json:"value"`
	Description string `json:"description,omitempty"`
	Default     bool   `json:"default,omitempty"`
	Emoji       *Emoji `json:"emoji,omitempty"`
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
	o.Emoji = emoji
	return o
}
