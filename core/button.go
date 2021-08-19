package core

import "github.com/DisgoOrg/disgo/discord"

// NewButton creates a new Button with the provided parameters. Link Button(s) need a URL and other Button(s) need a customID
//goland:noinspection GoUnusedExportedFunction
func NewButton(style discord.ButtonStyle, label string, customID string, url string, emoji *discord.Emoji, disabled bool) Button {
	return Button{
		Component: discord.Component{
			Type:     discord.ComponentTypeButton,
			Style:    style,
			CustomID: customID,
			URL:      url,
			Label:    label,
			Emoji:    emoji,
			Disabled: disabled,
		},
	}
}

// NewPrimaryButton creates a new Button with ButtonStylePrimary & the provided parameters
//goland:noinspection GoUnusedExportedFunction
func NewPrimaryButton(label string, customID string, emoji *discord.Emoji) Button {
	return NewButton(discord.ButtonStylePrimary, label, customID, "", emoji, false)
}

// NewSecondaryButton creates a new Button with ButtonStyleSecondary & the provided parameters
//goland:noinspection GoUnusedExportedFunction
func NewSecondaryButton(label string, customID string, emoji *discord.Emoji) Button {
	return NewButton(discord.ButtonStyleSecondary, label, customID, "", emoji, false)
}

// NewSuccessButton creates a new Button with ButtonStyleSuccess & the provided parameters
//goland:noinspection GoUnusedExportedFunction
func NewSuccessButton(label string, customID string, emoji *discord.Emoji) Button {
	return NewButton(discord.ButtonStyleSuccess, label, customID, "", emoji, false)
}

// NewDangerButton creates a new Button with ButtonStyleDanger & the provided parameters
//goland:noinspection GoUnusedExportedFunction
func NewDangerButton(label string, customID string, emoji *discord.Emoji) Button {
	return NewButton(discord.ButtonStyleDanger, label, customID, "", emoji, false)
}

// NewLinkButton creates a new link Button with ButtonStyleLink & the provided parameters
//goland:noinspection GoUnusedExportedFunction
func NewLinkButton(label string, url string, emoji *discord.Emoji) Button {
	return NewButton(discord.ButtonStyleLink, label, "", url, emoji, false)
}

// Button can be attacked to all messages & be clicked by a User. If clicked it fires an events.ButtonClickEvent with the declared customID
type Button struct {
	discord.Component
}

// Type returns the ComponentType of this Component
func (b Button) Type() discord.ComponentType {
	return b.Component.Type
}

// AsDisabled returns a new Button but disabled
func (b Button) AsDisabled() Button {
	b.Disabled = true
	return b
}

// AsEnabled returns a new Button but enabled
func (b Button) AsEnabled() Button {
	b.Disabled = true
	return b
}

// WithDisabled returns a new Button but disabled/enabled
func (b Button) WithDisabled(disabled bool) Button {
	b.Disabled = disabled
	return b
}

// WithEmoji returns a new Button with the provided Emoji
func (b Button) WithEmoji(emoji *discord.Emoji) Button {
	b.Emoji = emoji
	return b
}

// WithCustomID returns a new Button with the provided custom id
func (b Button) WithCustomID(customID string) Button {
	b.CustomID = customID
	return b
}

// WithStyle returns a new Button with the provided style
func (b Button) WithStyle(style discord.ButtonStyle) Button {
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
