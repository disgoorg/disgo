package discord

import "github.com/DisgoOrg/disgo/json"

// ButtonStyle defines how the ButtonComponent looks like (https://discord.com/assets/7bb017ce52cfd6575e21c058feb3883b.png)
type ButtonStyle int

// Supported ButtonStyle(s)
const (
	ButtonStylePrimary = iota + 1
	ButtonStyleSecondary
	ButtonStyleSuccess
	ButtonStyleDanger
	ButtonStyleLink
)

// NewButton creates a new ButtonComponent with the provided parameters. Link ButtonComponent(s) need a URL and other ButtonComponent(s) need a customID
//goland:noinspection GoUnusedExportedFunction
func NewButton(style ButtonStyle, label string, customID string, url string, emoji *ComponentEmoji, disabled bool) ButtonComponent {
	return ButtonComponent{
		Style:    style,
		CustomID: customID,
		URL:      url,
		Label:    label,
		Emoji:    emoji,
		Disabled: disabled,
	}
}

// NewPrimaryButton creates a new ButtonComponent with ButtonStylePrimary & the provided parameters
//goland:noinspection GoUnusedExportedFunction
func NewPrimaryButton(label string, customID string) ButtonComponent {
	return NewButton(ButtonStylePrimary, label, customID, "", nil, false)
}

// NewSecondaryButton creates a new ButtonComponent with ButtonStyleSecondary & the provided parameters
//goland:noinspection GoUnusedExportedFunction
func NewSecondaryButton(label string, customID string) ButtonComponent {
	return NewButton(ButtonStyleSecondary, label, customID, "", nil, false)
}

// NewSuccessButton creates a new ButtonComponent with ButtonStyleSuccess & the provided parameters
//goland:noinspection GoUnusedExportedFunction
func NewSuccessButton(label string, customID string) ButtonComponent {
	return NewButton(ButtonStyleSuccess, label, customID, "", nil, false)
}

// NewDangerButton creates a new ButtonComponent with ButtonStyleDanger & the provided parameters
//goland:noinspection GoUnusedExportedFunction
func NewDangerButton(label string, customID string) ButtonComponent {
	return NewButton(ButtonStyleDanger, label, customID, "", nil, false)
}

// NewLinkButton creates a new link ButtonComponent with ButtonStyleLink & the provided parameters
//goland:noinspection GoUnusedExportedFunction
func NewLinkButton(label string, url string) ButtonComponent {
	return NewButton(ButtonStyleLink, label, "", url, nil, false)
}

var _ Component = (*ButtonComponent)(nil)

type ButtonComponent struct {
	CustomID string          `json:"custom_id"`
	Style    ButtonStyle     `json:"style,omitempty"`
	Label    string          `json:"label,omitempty"`
	Emoji    *ComponentEmoji `json:"emoji,omitempty"`
	URL      string          `json:"url,omitempty"`
	Disabled bool            `json:"disabled,omitempty"`
}

func (b ButtonComponent) MarshalJSON() ([]byte, error) {
	type button ButtonComponent
	v := struct {
		Type ComponentType `json:"type"`
		button
	}{
		Type:   b.Type(),
		button: button(b),
	}
	return json.Marshal(v)
}

func (_ ButtonComponent) Type() ComponentType {
	return ComponentTypeButton
}

// AsEnabled returns a new ButtonComponent but enabled
func (b ButtonComponent) AsEnabled() ButtonComponent {
	b.Disabled = false
	return b
}

// AsDisabled returns a new ButtonComponent but disabled
func (b ButtonComponent) AsDisabled() ButtonComponent {
	b.Disabled = true
	return b
}

// WithDisabled returns a new ButtonComponent but disabled/enabled
func (b ButtonComponent) WithDisabled(disabled bool) ButtonComponent {
	b.Disabled = disabled
	return b
}

// WithEmoji returns a new ButtonComponent with the provided Emoji
func (b ButtonComponent) WithEmoji(emoji ComponentEmoji) ButtonComponent {
	b.Emoji = &emoji
	return b
}

// WithCustomID returns a new ButtonComponent with the provided custom id
func (b ButtonComponent) WithCustomID(customID string) ButtonComponent {
	b.CustomID = customID
	return b
}

// WithStyle returns a new ButtonComponent with the provided style
func (b ButtonComponent) WithStyle(style ButtonStyle) ButtonComponent {
	b.Style = style
	return b
}

// WithLabel returns a new ButtonComponent with the provided label
func (b ButtonComponent) WithLabel(label string) ButtonComponent {
	b.Label = label
	return b
}

// WithURL returns a new ButtonComponent with the provided URL
func (b ButtonComponent) WithURL(url string) ButtonComponent {
	b.URL = url
	return b
}
