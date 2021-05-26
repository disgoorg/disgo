package api

type ButtonStyle int

const (
	ButtonStylePrimary = iota + 1
	ButtonStyleSecondary
	ButtonStyleSuccess
	ButtonStyleDanger
	ButtonStyleLink
)

func NewButton(style ButtonStyle, label *string, customID *string, url *string, emoji *Emote, disabled *bool) *Button {
	return &Button{
		ComponentImpl: newComponentImpl(ComponentTypeButton),
		Style:         style,
		CustomID:      customID,
		URL:           url,
		Label:         label,
		Emoji:         emoji,
		Disabled:      disabled,
	}
}

func NewPrimaryButton(label string, customID string, emoji *Emote, disabled bool) *Button {
	return NewButton(ButtonStylePrimary, &label, &customID, nil, emoji, &disabled)
}

func NewSecondaryButton(label string, customID string, emoji *Emote, disabled bool) *Button {
	return NewButton(ButtonStyleSecondary, &label, &customID, nil, emoji, &disabled)
}

func NewSuccessButton(label string, customID string, emoji *Emote, disabled bool) *Button {
	return NewButton(ButtonStyleSuccess, &label, &customID, nil, emoji, &disabled)
}

func NewDangerButton(label string, customID string, emoji *Emote, disabled bool) *Button {
	return NewButton(ButtonStyleDanger, &label, &customID, nil, emoji, &disabled)
}

func NewLinkButton(label string, url string, emoji *Emote, disabled bool) *Button {
	return NewButton(ButtonStyleLink, &label, nil, &url, emoji, &disabled)
}

type Button struct {
	ComponentImpl
	Style    ButtonStyle `json:"style,omitempty"`
	Label    *string     `json:"label,omitempty"`
	Emoji    *Emote      `json:"emoji,omitempty"`
	CustomID *string     `json:"custom_id,omitempty"`
	URL      *string     `json:"url,omitempty"`
	Disabled *bool       `json:"disabled,omitempty"`
}
