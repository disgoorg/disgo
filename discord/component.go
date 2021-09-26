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

// Component is used for easier unmarshalling of different Component(s)
type Component struct {
	Type ComponentType `json:"type"`

	// Button && SelectMenu
	CustomID string `json:"custom_id,omitempty"`

	// Button
	Style    ButtonStyle `json:"style,omitempty"`
	Label    string      `json:"label,omitempty"`
	Emoji    *Emoji      `json:"emoji,omitempty"`
	URL      string      `json:"url,omitempty"`
	Disabled bool        `json:"disabled,omitempty"`

	// ActionRow
	Components []Component `json:"components,omitempty"`

	// SelectMenu
	Placeholder string         `json:"placeholder,omitempty"`
	MinValues   int            `json:"min_values,omitempty"`
	MaxValues   int            `json:"max_values,omitempty"`
	Options     []SelectOption `json:"options,omitempty"`
}

// SelectOption represents an option in a SelectMenu
type SelectOption struct {
	Label       string `json:"label"`
	Value       string `json:"value"`
	Description string `json:"description,omitempty"`
	Emoji       *Emoji `json:"emoji,omitempty"`
	Default     bool   `json:"default,omitempty"`
}
