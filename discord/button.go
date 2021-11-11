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
func NewButton(style ButtonStyle, label string, customID CustomID, url string, emoji *ComponentEmoji, disabled bool) ButtonComponent {
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
func NewPrimaryButton(label string, customID CustomID) ButtonComponent {
	return NewButton(ButtonStylePrimary, label, customID, "", nil, false)
}

// NewSecondaryButton creates a new ButtonComponent with ButtonStyleSecondary & the provided parameters
//goland:noinspection GoUnusedExportedFunction
func NewSecondaryButton(label string, customID CustomID) ButtonComponent {
	return NewButton(ButtonStyleSecondary, label, customID, "", nil, false)
}

// NewSuccessButton creates a new ButtonComponent with ButtonStyleSuccess & the provided parameters
//goland:noinspection GoUnusedExportedFunction
func NewSuccessButton(label string, customID CustomID) ButtonComponent {
	return NewButton(ButtonStyleSuccess, label, customID, "", nil, false)
}

// NewDangerButton creates a new ButtonComponent with ButtonStyleDanger & the provided parameters
//goland:noinspection GoUnusedExportedFunction
func NewDangerButton(label string, customID CustomID) ButtonComponent {
	return NewButton(ButtonStyleDanger, label, customID, "", nil, false)
}

// NewLinkButton creates a new link ButtonComponent with ButtonStyleLink & the provided parameters
//goland:noinspection GoUnusedExportedFunction
func NewLinkButton(label string, url string) ButtonComponent {
	return NewButton(ButtonStyleLink, label, "", url, nil, false)
}

var (
	_ Component            = (*ButtonComponent)(nil)
	_ InteractiveComponent = (*ButtonComponent)(nil)
)

type ButtonComponent struct {
	CustomID CustomID        `json:"custom_id"`
	Style    ButtonStyle     `json:"style,omitempty"`
	Label    string          `json:"label,omitempty"`
	Emoji    *ComponentEmoji `json:"emoji,omitempty"`
	URL      string          `json:"url,omitempty"`
	Disabled bool            `json:"disabled,omitempty"`
}

func (c ButtonComponent) MarshalJSON() ([]byte, error) {
	type buttonComponent ButtonComponent
	return json.Marshal(struct {
		Type ComponentType `json:"type"`
		buttonComponent
	}{
		Type:            c.Type(),
		buttonComponent: buttonComponent(c),
	})
}

func (c ButtonComponent) Type() ComponentType {
	return ComponentTypeButton
}

func (c ButtonComponent) ID() CustomID {
	return c.CustomID
}

func (c ButtonComponent) component()            {}
func (c ButtonComponent) interactiveComponent() {}

// AsEnabled returns a new ButtonComponent but enabled
func (c ButtonComponent) AsEnabled() ButtonComponent {
	c.Disabled = false
	return c
}

// AsDisabled returns a new ButtonComponent but disabled
func (c ButtonComponent) AsDisabled() ButtonComponent {
	c.Disabled = true
	return c
}

// WithDisabled returns a new ButtonComponent but disabled/enabled
func (c ButtonComponent) WithDisabled(disabled bool) ButtonComponent {
	c.Disabled = disabled
	return c
}

// WithEmoji returns a new ButtonComponent with the provided Emoji
func (c ButtonComponent) WithEmoji(emoji ComponentEmoji) ButtonComponent {
	c.Emoji = &emoji
	return c
}

// WithCustomID returns a new ButtonComponent with the provided custom id
func (c ButtonComponent) WithCustomID(customID CustomID) ButtonComponent {
	c.CustomID = customID
	return c
}

// WithStyle returns a new ButtonComponent with the provided style
func (c ButtonComponent) WithStyle(style ButtonStyle) ButtonComponent {
	c.Style = style
	return c
}

// WithLabel returns a new ButtonComponent with the provided label
func (c ButtonComponent) WithLabel(label string) ButtonComponent {
	c.Label = label
	return c
}

// WithURL returns a new ButtonComponent with the provided URL
func (c ButtonComponent) WithURL(url string) ButtonComponent {
	c.URL = url
	return c
}
