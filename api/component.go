package api

type ComponentType int

const (
	ComponentTypeButtons = iota + 1
	ComponentTypeButton
	ComponentTypeSelect
)

type Style int

const (
	StyleBlurple = iota + 1
	StyleGrey
	StyleGreen
	StyleRed
	StyleHyperlink
)

type Component interface {
	Type() ComponentType
}

type ComponentImpl struct {
	ComponentType ComponentType `json:"type"`
}

func (t ComponentImpl) Type() ComponentType {
	return t.ComponentType
}

func NewEmote(emoteID Snowflake) *Emoji {
	return &Emoji{ID: emoteID, Animated: false}
}

func NewAnimatedEmote(emoteID Snowflake) *Emoji {
	return &Emoji{ID: emoteID, Animated: true}
}

func NewEmoji(name string) *Emoji {
	return &Emoji{Name: name}
}

type Emoji struct {
	Name     string    `json:"name,omitempty"`
	ID       Snowflake `json:"id,omitempty"`
	Animated bool      `json:"animated"`
}
