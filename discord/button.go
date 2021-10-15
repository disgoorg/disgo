package discord

import "github.com/DisgoOrg/disgo/json"

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

// NewButton creates a new Button with the provided parameters. Link Button(s) need a URL and other Button(s) need a customID
//goland:noinspection GoUnusedExportedFunction
func NewButton(style ButtonStyle, label string, customID string, url string, emoji *ComponentEmoji, disabled bool) Button {
	return Button{
		Style:    style,
		CustomID: customID,
		URL:      url,
		Label:    label,
		Emoji:    emoji,
		Disabled: disabled,
	}
}

// NewPrimaryButton creates a new Button with ButtonStylePrimary & the provided parameters
//goland:noinspection GoUnusedExportedFunction
func NewPrimaryButton(label string, customID string) Button {
	return NewButton(ButtonStylePrimary, label, customID, "", nil, false)
}

// NewSecondaryButton creates a new Button with ButtonStyleSecondary & the provided parameters
//goland:noinspection GoUnusedExportedFunction
func NewSecondaryButton(label string, customID string) Button {
	return NewButton(ButtonStyleSecondary, label, customID, "", nil, false)
}

// NewSuccessButton creates a new Button with ButtonStyleSuccess & the provided parameters
//goland:noinspection GoUnusedExportedFunction
func NewSuccessButton(label string, customID string) Button {
	return NewButton(ButtonStyleSuccess, label, customID, "", nil, false)
}

// NewDangerButton creates a new Button with ButtonStyleDanger & the provided parameters
//goland:noinspection GoUnusedExportedFunction
func NewDangerButton(label string, customID string) Button {
	return NewButton(ButtonStyleDanger, label, customID, "", nil, false)
}

// NewLinkButton creates a new link Button with ButtonStyleLink & the provided parameters
//goland:noinspection GoUnusedExportedFunction
func NewLinkButton(label string, url string) Button {
	return NewButton(ButtonStyleLink, label, "", url, nil, false)
}

var _ Component = (*Button)(nil)

type Button struct {
	CustomID string          `json:"custom_id"`
	Style    ButtonStyle     `json:"style,omitempty"`
	Label    string          `json:"label,omitempty"`
	Emoji    *ComponentEmoji `json:"emoji,omitempty"`
	URL      string          `json:"url,omitempty"`
	Disabled bool            `json:"disabled,omitempty"`
}

func (b Button) MarshalJSON() ([]byte, error) {
	type button Button
	v := struct {
		Type ComponentType `json:"type"`
		button
	}{
		Type:   b.Type(),
		button: button(b),
	}
	return json.Marshal(v)
}

func (_ Button) Type() ComponentType {
	return ComponentTypeButton
}

// AsEnabled returns a new Button but enabled
func (b Button) AsEnabled() Button {
	b.Disabled = false
	return b
}

// AsDisabled returns a new Button but disabled
func (b Button) AsDisabled() Button {
	b.Disabled = true
	return b
}

// WithDisabled returns a new Button but disabled/enabled
func (b Button) WithDisabled(disabled bool) Button {
	b.Disabled = disabled
	return b
}

// WithEmoji returns a new Button with the provided Emoji
func (b Button) WithEmoji(emoji ComponentEmoji) Button {
	b.Emoji = &emoji
	return b
}

// WithCustomID returns a new Button with the provided custom id
func (b Button) WithCustomID(customID string) Button {
	b.CustomID = customID
	return b
}

// WithStyle returns a new Button with the provided style
func (b Button) WithStyle(style ButtonStyle) Button {
	b.Style = style
	return b
}

// WithLabel returns a new Button with the provided label
func (b Button) WithLabel(label string) Button {
	b.Label = label
	return b
}

// WithURL returns a new Button with the provided URL
func (b Button) WithURL(url string) Button {
	b.URL = url
	return b
}
