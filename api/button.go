package api

func NewBlurpleButton(label string, customID string, emoji *Emoji, disabled bool) *Button {
	return &Button{
		ComponentImpl: ComponentImpl{
			ComponentType: ComponentTypeButton,
		},
		Style:    StyleBlurple,
		CustomID: customID,
		Label:    label,
		Emoji:    emoji,
		Disabled: disabled,
	}
}

func NewGreyButton(label string, customID string, emoji *Emoji, disabled bool) *Button {
	return &Button{
		ComponentImpl: ComponentImpl{
			ComponentType: ComponentTypeButton,
		},
		Style:    StyleGrey,
		CustomID: customID,
		Label:    label,
		Emoji:    emoji,
		Disabled: disabled,
	}
}

func NewGreenButton(label string, customID string, emoji *Emoji, disabled bool) *Button {
	return &Button{
		ComponentImpl: ComponentImpl{
			ComponentType: ComponentTypeButton,
		},
		Style:    StyleGreen,
		CustomID: customID,
		Label:    label,
		Emoji:    emoji,
		Disabled: disabled,
	}
}

func NewRedButton(label string, customID string, emoji *Emoji, disabled bool) *Button {
	return &Button{
		ComponentImpl: ComponentImpl{
			ComponentType: ComponentTypeButton,
		},
		Style:    StyleRed,
		CustomID: customID,
		Label:    label,
		Emoji:    emoji,
		Disabled: disabled,
	}
}

func NewLinkButton(label string, URL string, emoji *Emoji, disabled bool) *Button {
	return &Button{
		ComponentImpl: ComponentImpl{
			ComponentType: ComponentTypeButton,
		},
		Style:    StyleHyperlink,
		Label:    label,
		URL:      URL,
		Emoji:    emoji,
		Disabled: disabled,
	}
}

type Button struct {
	ComponentImpl
	Style    Style  `json:"style,omitempty"`
	CustomID string `json:"custom_id,omitempty"`
	Label    string `json:"label,omitempty"`
	URL      string `json:"url,omitempty"`
	Emoji    *Emoji `json:"emoji,omitempty"`
	Disabled bool   `json:"disabled,omitempty"`
}